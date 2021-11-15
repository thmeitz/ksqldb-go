package ksqldb_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
	mocknet "github.com/thmeitz/ksqldb-go/mocks/net"
	"github.com/thmeitz/ksqldb-go/net"
)

func TestHandleGetRequest_StatusCodeError(t *testing.T) {
	fn := ksqldb.HandleGetRequest
	r := ioutil.NopCloser(bytes.NewReader([]byte("")))
	res := http.Response{StatusCode: 400, Body: r}

	m := &mocknet.HTTPClient{}
	m.On("Get", mock.Anything).Return(&res, nil)

	_, err := fn(m, "/bla")

	require.NotNil(t, err)
}

func TestHandleRequestError(t *testing.T) {

}

func TestNewCloseQueryRequest(t *testing.T) {
	postFn := ksqldb.NewCloseQueryRequest
	client, _ := net.NewHTTPClient(net.Options{}, nil)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hallo")))
	req, err := postFn(&client, context.TODO(), r)
	require.NotNil(t, req)
	require.Nil(t, err)
	require.Equal(t, "/close-query", req.URL.Path)
}

func TestNewKsqlRequest(t *testing.T) {
	postFn := ksqldb.NewKsqlRequest
	client, _ := net.NewHTTPClient(net.Options{}, nil)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hallo")))
	req, err := postFn(&client, r)
	require.NotNil(t, req)
	require.Nil(t, err)
	require.Equal(t, "/ksql", req.URL.Path)
}

func TestNewPostRequest_Error(t *testing.T) {
	postFn := ksqldb.NewPostRequest

	client, _ := net.NewHTTPClient(net.Options{}, nil)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hallo")))
	req, err := postFn(&client, nil, "/bla", r)
	require.Nil(t, req)
	require.NotNil(t, err)
	require.Equal(t, "can't create new request with context: net/http: nil Context", err.Error())
}

func TestNewPostRequest_Successful(t *testing.T) {
	postFn := ksqldb.NewPostRequest

	client, _ := net.NewHTTPClient(net.Options{}, nil)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hallo")))
	req, err := postFn(&client, context.TODO(), "/bla", r)
	require.NotNil(t, req)
	require.Nil(t, err)
}

func TestNewQueryRequest(t *testing.T) {

}

func TestNewQueryStreamRequest(t *testing.T) {
	postFn := ksqldb.NewQueryStreamRequest
	client, _ := net.NewHTTPClient(net.Options{}, nil)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hallo")))
	req, err := postFn(&client, context.TODO(), r)
	require.NotNil(t, req)
	require.Nil(t, err)
	require.Equal(t, "localhost:8088", req.Host)
	require.Equal(t, "POST", req.Method)
	require.Equal(t, "/query-stream", req.URL.Path)
}
