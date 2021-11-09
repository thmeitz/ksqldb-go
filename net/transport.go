/*
Copyright © 2021 Thomas Meitz <thme219@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Parts of this apiclient are borrowed from Zalando Skipper
https://github.com/zalando/skipper/blob/master/net/httpclient.go

Zalando licence: MIT
https://github.com/zalando/skipper/blob/master/LICENSE

Next: change opentracing-go to opentelemetry if stable version is released
*/

package net

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/net/http2"
)

// Options configures the Client
type Options struct {
	// BaseUrl is the URL to your ksqldb server
	BaseUrl string
	// Credentials for BaseAuth remote user
	Credentials Credentials
	// AllowHTTP
	AllowHTTP bool
	// DisableKeepAlives see https://golang.org/pkg/net/http/#Transport.DisableKeepAlives
	DisableKeepAlives bool
	// DisableCompression see https://golang.org/pkg/net/http/#Transport.DisableCompression
	DisableCompression bool
	// ForceAttemptHTTP2 see https://golang.org/pkg/net/http/#Transport.ForceAttemptHTTP2
	ForceAttemptHTTP2 bool
	// MaxIdleConns see https://golang.org/pkg/net/http/#Transport.MaxIdleConns
	MaxIdleConns int
	// MaxIdleConnsPerHost see https://golang.org/pkg/net/http/#Transport.MaxIdleConnsPerHost
	MaxIdleConnsPerHost int
	// MaxConnsPerHost see https://golang.org/pkg/net/http/#Transport.MaxConnsPerHost
	MaxConnsPerHost int
	// WriteBufferSize see https://golang.org/pkg/net/http/#Transport.WriteBufferSize
	WriteBufferSize int
	// ReadBufferSize see https://golang.org/pkg/net/http/#Transport.ReadBufferSize
	ReadBufferSize int
	// MaxResponseHeaderBytes see
	// https://golang.org/pkg/net/http/#Transport.MaxResponseHeaderBytes
	MaxResponseHeaderBytes int64
	// Timeout sets all Timeouts, that are set to 0 to the given
	// value. Basically it's the default timeout value.
	Timeout time.Duration
	// TLSHandshakeTimeout see
	// https://golang.org/pkg/net/http/#Transport.TLSHandshakeTimeout,
	// if not set or set to 0, its using Options.Timeout.
	TLSHandshakeTimeout time.Duration
	// IdleConnTimeout see
	// https://golang.org/pkg/net/http/#Transport.IdleConnTimeout,
	// if not set or set to 0, its using Options.Timeout.
	IdleConnTimeout time.Duration
	// ResponseHeaderTimeout see
	// https://golang.org/pkg/net/http/#Transport.ResponseHeaderTimeout,
	// if not set or set to 0, its using Options.Timeout.
	ResponseHeaderTimeout time.Duration
	// ExpectContinueTimeout see
	// https://golang.org/pkg/net/http/#Transport.ExpectContinueTimeout,
	// if not set or set to 0, its using Options.Timeout.
	ExpectContinueTimeout time.Duration
	// Tracer instance, can be nil to not enable tracing
	Tracer opentracing.Tracer
	// OpentracingComponentTag sets component tag for all requests
	OpentracingComponentTag string
	// OpentracingSpanName sets span name for all requests
	OpentracingSpanName string
}

// Transport wraps an http.Transport and adds support for tracing and
// bearerToken injection.
type Transport struct {
	quit          chan struct{}
	closed        bool
	tr            *http.Transport
	tr2           *http2.Transport
	tracer        opentracing.Tracer
	spanName      string
	componentName string
	// bearerToken        string

}

// NewTransport creates a new Transport with Options
func NewTransport(options Options) *Transport {
	// set default tracer
	if options.Tracer == nil {
		options.Tracer = &opentracing.NoopTracer{}
	}

	// set timeout defaults
	if options.TLSHandshakeTimeout == 0 {
		options.TLSHandshakeTimeout = options.Timeout
	}
	if options.IdleConnTimeout == 0 {
		if options.Timeout != 0 {
			options.IdleConnTimeout = options.Timeout
		} else {
			options.IdleConnTimeout = DefaultIdleConnTimeout
		}
	}
	if options.ResponseHeaderTimeout == 0 {
		options.ResponseHeaderTimeout = options.Timeout
	}
	if options.ExpectContinueTimeout == 0 {
		options.ExpectContinueTimeout = options.Timeout
	}

	htransport := &http.Transport{
		DisableKeepAlives:      options.DisableKeepAlives,
		DisableCompression:     options.DisableCompression,
		ForceAttemptHTTP2:      options.ForceAttemptHTTP2,
		MaxIdleConns:           options.MaxIdleConns,
		MaxIdleConnsPerHost:    options.MaxIdleConnsPerHost,
		MaxConnsPerHost:        options.MaxConnsPerHost,
		WriteBufferSize:        options.WriteBufferSize,
		ReadBufferSize:         options.ReadBufferSize,
		MaxResponseHeaderBytes: options.MaxResponseHeaderBytes,
		ResponseHeaderTimeout:  options.ResponseHeaderTimeout,
		TLSHandshakeTimeout:    options.TLSHandshakeTimeout,
		IdleConnTimeout:        options.IdleConnTimeout,
		ExpectContinueTimeout:  options.ExpectContinueTimeout,
	}
	var htransport2 = &http2.Transport{}
	if options.AllowHTTP {
		// ksqlDB uses HTTP2 and if the server is on HTTP then Golang will not
		// use HTTP2 unless we force it to, thus.
		// Without this you get the error `http2: unsupported scheme`
		htransport2.AllowHTTP = options.AllowHTTP
		// Pretend we are dialing a TLS endpoint.
		// Note, we ignore the passed tls.Config
		htransport2.DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		}
		t2 := &Transport{
			quit:   make(chan struct{}),
			tr2:    htransport2,
			tracer: options.Tracer,
		}
		if t2.tracer != nil {
			if options.OpentracingComponentTag != "" {
				t2 = WithComponentTag(t2, options.OpentracingComponentTag)
			}
			if options.OpentracingSpanName != "" {
				t2 = WithSpanName(t2, options.OpentracingSpanName)
			}
		}
		go func() {
			for {
				select {
				case <-time.After(options.IdleConnTimeout):
					htransport2.CloseIdleConnections()
				case <-t2.quit:
					return
				}
			}
		}()

		return t2

	} else {
		t := &Transport{
			quit:   make(chan struct{}),
			tr:     htransport,
			tracer: options.Tracer,
		}

		if t.tracer != nil {
			if options.OpentracingComponentTag != "" {
				t = WithComponentTag(t, options.OpentracingComponentTag)
			}
			if options.OpentracingSpanName != "" {
				t = WithSpanName(t, options.OpentracingSpanName)
			}
		}

		go func() {
			for {
				select {
				case <-time.After(options.IdleConnTimeout):
					htransport.CloseIdleConnections()
				case <-t.quit:
					return
				}
			}
		}()

		return t
	}

}

// Close the transport
func (t *Transport) Close() {
	if !t.closed {
		t.closed = true
		close(t.quit)
	}
}

// CloseIdleConnection closes idle connections
// func (t *Transport) CloseIdleConnections() {
// 	t.tr.CloseIdleConnections()
// }

// RoundTrip the request with tracing and add client
// tracing: DNS, TCP/IP, TLS handshake, connection pool access. Client
// traces are added as logs into the created span.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var span opentracing.Span
	var err error
	var rsp *http.Response
	if t.tr != nil {
		if t.spanName != "" {
			req, span = t.injectSpan(req)
			defer span.Finish()
			req = injectClientTrace(req, span)
			span.LogKV("http_do", "start")
		}

		rsp, err = t.tr.RoundTrip(req)
		if span != nil {
			span.LogKV("http_do", "stop")
			if rsp != nil {
				ext.HTTPStatusCode.Set(span, uint16(rsp.StatusCode))
			}
		}
	} else {
		rsp, err = t.tr2.RoundTrip(req)
	}

	return rsp, err
}

// injectSpan injects an opentracing span into the request
func (t *Transport) injectSpan(req *http.Request) (*http.Request, opentracing.Span) {
	parentSpan := opentracing.SpanFromContext(req.Context())
	var span opentracing.Span
	if parentSpan != nil {
		req = req.WithContext(opentracing.ContextWithSpan(req.Context(), parentSpan))
		span = t.tracer.StartSpan(t.spanName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		span = t.tracer.StartSpan(t.spanName)
	}

	// add Tags
	ext.Component.Set(span, t.componentName)
	ext.HTTPUrl.Set(span, req.URL.String())
	ext.HTTPMethod.Set(span, req.Method)
	ext.SpanKind.Set(span, "client")

	_ = t.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	return req, span
}

// WithComponentTag sets the component name, if you have an enabled
// tracing Transport.
func WithComponentTag(t *Transport, componentName string) *Transport {
	tt := t.shallowCopy()
	tt.componentName = componentName
	return tt
}

// WithSpanName sets the name of the span, if you have an enabled
// tracing Transport.
func WithSpanName(t *Transport, spanName string) *Transport {
	tt := t.shallowCopy()
	tt.spanName = spanName
	return tt
}

// shallowCopy copies the transport
func (t *Transport) shallowCopy() *Transport {
	tt := *t
	return &tt
}

// injectClientTrace injects traces into the context of the http.Request and returns it
func injectClientTrace(req *http.Request, span opentracing.Span) *http.Request {
	trace := &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) {
			span.LogKV("DNS", "start")
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			span.LogKV("DNS", "end")
		},
		ConnectStart: func(string, string) {
			span.LogKV("connect", "start")
		},
		ConnectDone: func(string, string, error) {
			span.LogKV("connect", "end")
		},
		TLSHandshakeStart: func() {
			span.LogKV("TLS", "start")
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			span.LogKV("TLS", "end")
		},
		GetConn: func(string) {
			span.LogKV("get_conn", "start")
		},
		GotConn: func(httptrace.GotConnInfo) {
			span.LogKV("get_conn", "end")
		},
		WroteHeaders: func() {
			span.LogKV("wrote_headers", "done")
		},
		WroteRequest: func(wri httptrace.WroteRequestInfo) {
			if wri.Err != nil {
				span.LogKV("wrote_request", wri.Err.Error())
			} else {
				span.LogKV("wrote_request", "done")
			}
		},
		GotFirstResponseByte: func() {
			span.LogKV("got_first_byte", "done")
		},
	}
	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
}
