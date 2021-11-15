package ksqldb_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

func TestHandleGetRequest(t *testing.T) {

}

func TestHandleRequestError(t *testing.T) {

}

func TestNewCloseQueryRequest(t *testing.T) {

}

func TestNewKsqlRequest(t *testing.T) {

}

func TestNewPostRequest_Error(t *testing.T) {
	postFn := ksqldb.NewPostRequest

	client, _ := net.NewHTTPClient(net.Options{}, nil)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hallo")))
	req, err := postFn(client, nil, "/bla", r)
	require.Nil(t, req)
	require.NotNil(t, err)
	require.Equal(t, "can't create new request with context: net/http: nil Context", err.Error())
}

func TestNewPostRequest_Successful(t *testing.T) {
	postFn := ksqldb.NewPostRequest

	client, _ := net.NewHTTPClient(net.Options{}, nil)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hallo")))
	req, err := postFn(client, context.TODO(), "/bla", r)
	require.NotNil(t, req)
	require.Nil(t, err)
}

func TestNewQueryRequest(t *testing.T) {

}

func TestNewQueryStreamRequest(t *testing.T) {

}
