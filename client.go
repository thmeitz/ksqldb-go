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
*/

package ksqldb

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/thmeitz/ksqldb-go/internal"
)

const (
	DefaultIdleConnTimeout = 10 * time.Second
	DefaultBaseUrl         = "http://localhost:8088"
)

// The ksqlDB client
type Client struct {
	options Options
	uri     *url.URL
	client  http.Client
	tr      *Transport
	logger  log.Logger
}

// Credentials holds the username and password
type Credentials struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

func NewClient(options Options, logger log.Logger) (*Client, error) {
	var uri *url.URL
	var err error

	if options.BaseUrl == "" {
		options.BaseUrl = DefaultBaseUrl
	}

	if uri, err = internal.GetUrl(options.BaseUrl); err != nil {
		return nil, fmt.Errorf("%+w", err)
	}

	tr := NewTransport(options)

	return &Client{
		logger: logger,
		client: http.Client{
			Transport: tr,
		},
		options: options,
		tr:      tr,
		uri:     uri,
	}, nil
}

func (c *Client) Close() {
	c.tr.Close()
}

func (cl *Client) newQueryStreamRequest(ctx context.Context, payload io.Reader) (*http.Request, error) {
	req, err := cl.newPostRequest(ctx, QUERY_STREAM_ENDPOINT, payload)
	return req, err
}

func (cl *Client) newCloseQueryRequest(ctx context.Context, payload io.Reader) (*http.Request, error) {
	return cl.newPostRequest(ctx, CLOSE_QUERY_ENDPOINT, payload)
}

func (cl *Client) newKsqlRequest(payload io.Reader) (*http.Request, error) {
	return http.NewRequest("POST", cl.options.BaseUrl+KSQL_ENDPOINT, payload)
}

func (cl *Client) newPostRequest(ctx context.Context, endpoint string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", cl.options.BaseUrl+endpoint, payload)
	if err != nil {
		return req, fmt.Errorf("can't create new request with context:\n%w", err)
	}

	return req, nil
}

// SanitizeQuery sanitizes the given content
// eventually we can use the KSqlParser to rewrite the query, so its automatically sanitized
// whitespaces will be eaten by the KSqlParser
func (cl *Client) SanitizeQuery(content string) string {
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, "\n", "")
	return content
}
