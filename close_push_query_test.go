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

func TestClosePushQuery(t *testing.T) {
	var data = `[]`

	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)

	r := ioutil.NopCloser(bytes.NewReader([]byte(data)))
	res := http.Response{StatusCode: 200, Body: r}

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/close-query")
	m.On("Do", mock.Anything).Return(&res, nil)

	err := kcl.ClosePushQuery(context.TODO(), "12345")
	require.Nil(t, err)
}

func TestClosePushQuery_FailedDoRequest(t *testing.T) {
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/close-query")
	m.On("Do", mock.Anything).Return(nil, errors.New("error"))

	err := kcl.ClosePushQuery(context.TODO(), "12345")
	require.NotNil(t, err)
	require.Equal(t, "failed to execute post request to terminate query: error", err.Error())
}

func TestClosePushQuery_RequestStatusCode(t *testing.T) {
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)

	json := `{}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 400, Body: r}

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/close-query")
	m.On("Do", mock.Anything).Return(&res, nil)

	err := kcl.ClosePushQuery(context.TODO(), "12345")
	require.NotNil(t, err)
}
