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
	"github.com/thmeitz/ksqldb-go/net"
)

func TestValidateProperty_EmptyProperty(t *testing.T) {
	ctx := context.Background()
	options := net.Options{
		BaseUrl:   "http://localhost:8088",
		AllowHTTP: true,
	}
	kcl, err := ksqldb.NewClientWithOptions(options)
	require.Nil(t, err)
	require.NotNil(t, kcl)
	val, err := kcl.ValidateProperty(ctx, "")
	require.Equal(t, "property must not empty", err.Error())
	require.Nil(t, val)
}

func TestValidateProperty_RequestError(t *testing.T) {
	ctx := context.Background()
	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/is_valid_property/test")
	m.Mock.On("Get", ctx, mock.Anything).Return(nil, errors.New("error"))
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.ValidateProperty(ctx, "test")
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "ksqldb get request failed: error", err.Error())
}

func TestValidateProperty_UnmarshalError(t *testing.T) {
	ctx := context.Background()
	json := `{"name":"Test Name"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}

	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/is_valid_property/test")
	m.Mock.On("Get", ctx, mock.Anything).Return(&res, nil)
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.ValidateProperty(ctx, "test")
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "could not parse the response:json: cannot unmarshal object into Go value of type bool", err.Error())
}

func TestValidateProperty_Successfull(t *testing.T) {
	ctx := context.Background()
	json := `true`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}

	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/is_valid_property/test")
	m.Mock.On("Get", ctx, mock.Anything).Return(&res, nil)
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.ValidateProperty(ctx, "test")
	require.NotNil(t, val)
	require.Nil(t, err)
	require.True(t, *val)
}
