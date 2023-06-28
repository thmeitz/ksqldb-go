/*
Copyright © 2021 Thomas Meitz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ksqldb_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
	mocknet "github.com/thmeitz/ksqldb-go/mocks/net"
)

func TestQueryStatus_EmptyCommandId(t *testing.T) {
	ctx := context.Background()
	m := mocknet.HTTPClient{}

	kcl, _ := ksqldb.NewClient(&m)
	result, err := kcl.GetQueryStatus(ctx, "")

	require.Nil(t, result)
	require.NotNil(t, err)
	require.Equal(t, "commandId is empty", err.Error())
}

func TestQueryStatus_GetError(t *testing.T) {
	ctx := context.Background()
	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/status/stream/PAGEVIEWS/create")
	m.Mock.On("Get", ctx, mock.Anything).Return(nil, errors.New("error"))
	m.Mock.On("Close").Return()

	kcl, _ := ksqldb.NewClient(&m)
	kcl.Close()
	_, err := kcl.GetQueryStatus(ctx, "/stream/PAGEVIEWS/create")

	require.NotNil(t, err)
	require.Equal(t, "ksqldb get request failed: error", err.Error())
	m.AssertCalled(t, "Close")
}

func TestQueryStatus_UnmarshalError(t *testing.T) {
	ctx := context.Background()
	json := `true`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}
	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/status/stream/PAGEVIEWS/create")
	m.Mock.On("Get", ctx, mock.Anything).Return(&res, nil)

	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.GetQueryStatus(ctx, "/stream/PAGEVIEWS/create")
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "could not parse the response:json: cannot unmarshal bool into Go value of type ksqldb.QueryStatus", err.Error())
}

func TestQueryStatus_Successful(t *testing.T) {
	ctx := context.Background()
	json := `{"status": "SUCCESS","message":"Stream created and running"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}
	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/status/stream/PAGEVIEWS/create")
	m.Mock.On("Get", ctx, mock.Anything).Return(&res, nil)

	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.GetQueryStatus(ctx, "/stream/PAGEVIEWS/create")
	require.Nil(t, err)
	require.NotNil(t, val)
}
