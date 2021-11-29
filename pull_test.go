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

func TestPull_ParseSQLError(t *testing.T) {
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)
	_, _, err := kcl.Pull(context.TODO(), ksqldb.QueryOptions{Sql: "select * from bla"})
	require.NotNil(t, err)
	require.Equal(t, "1 sql syntax error(s) found", err.Error())
}

func TestPull_RequestError(t *testing.T) {
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/query-stream")
	m.On("Do", mock.Anything).Return(nil, errors.New("error"))

	_, _, err := kcl.Pull(context.TODO(), ksqldb.QueryOptions{Sql: "select * from bla;"})
	require.NotNil(t, err)
	require.Equal(t, "can't do request: error", err.Error())
}

func TestPull_RequestStatusCode(t *testing.T) {
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)

	json := `{"name":"Test Name"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 400, Body: r}

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/query-stream")
	m.On("Do", mock.Anything).Return(&res, nil)

	_, _, err := kcl.Pull(context.TODO(), ksqldb.QueryOptions{Sql: "select * from bla;"})
	require.NotNil(t, err)
	require.Equal(t, "", err.Error())
}

func TestPull_UnmarshallError(t *testing.T) {
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)

	json := `{"name":"Test Name"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/query-stream")
	m.On("Do", mock.Anything).Return(&res, nil)

	_, _, err := kcl.Pull(context.TODO(), ksqldb.QueryOptions{Sql: "select * from bla;"})
	require.NotNil(t, err)
	require.Equal(t, "could not parse the response:\njson: cannot unmarshal object into Go value of type []interface {}", err.Error())
}

func TestPull_HeaderWithoutData(t *testing.T) {
	var nodata = `[
	{
		"queryId":null,
		"columnNames":[
			"WINDOW_START","WINDOW_END","DOG_SIZE","DOGS_CT"
		],
		"columnTypes":[
			"STRING","STRING","STRING","BIGINT"
		]
	}]`

	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)

	r := ioutil.NopCloser(bytes.NewReader([]byte(nodata)))
	res := http.Response{StatusCode: 200, Body: r}

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/query-stream")
	m.On("Do", mock.Anything).Return(&res, nil)

	header, _, err := kcl.Pull(context.TODO(), ksqldb.QueryOptions{Sql: "select * from bla;"})
	require.Nil(t, err)
	require.Equal(t, "", header.QueryId)
	require.Equal(t, 4, len(header.Columns))
	require.Equal(t, "WINDOW_START", header.Columns[0].Name)
	require.Equal(t, "STRING", header.Columns[0].Type)
	require.Equal(t, "WINDOW_END", header.Columns[1].Name)
	require.Equal(t, "STRING", header.Columns[1].Type)
	require.Equal(t, "DOG_SIZE", header.Columns[2].Name)
	require.Equal(t, "STRING", header.Columns[2].Type)
	require.Equal(t, "DOGS_CT", header.Columns[3].Name)
	require.Equal(t, "BIGINT", header.Columns[3].Type)
}

func TestPull_HeaderWithData(t *testing.T) {
	var data = `[
	{
		"queryId":"0815",
		"columnNames":[
			"WINDOW_START","WINDOW_END","DOG_SIZE","DOGS_CT"
		],
		"columnTypes":[
			"STRING","STRING","STRING","BIGINT"
		]
	},
	["2021-11-16 06:00:00","06:15:00","medium",23],
	["2021-11-16 06:15:00","06:30:00","medium",250]
]
`

	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)

	r := ioutil.NopCloser(bytes.NewReader([]byte(data)))
	res := http.Response{StatusCode: 200, Body: r}

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/query-stream")
	m.On("Do", mock.Anything).Return(&res, nil)

	header, payload, err := kcl.Pull(context.TODO(), ksqldb.QueryOptions{Sql: "select * from bla;"})
	require.Nil(t, err)
	require.Equal(t, "0815", header.QueryId)
	require.Equal(t, 4, len(header.Columns))
	require.Equal(t, "WINDOW_START", header.Columns[0].Name)
	require.Equal(t, "STRING", header.Columns[0].Type)
	require.Equal(t, "WINDOW_END", header.Columns[1].Name)
	require.Equal(t, "STRING", header.Columns[1].Type)
	require.Equal(t, "DOG_SIZE", header.Columns[2].Name)
	require.Equal(t, "STRING", header.Columns[2].Type)
	require.Equal(t, "DOGS_CT", header.Columns[3].Name)
	require.Equal(t, "BIGINT", header.Columns[3].Type)
	require.Equal(t, 2, len(payload))
	require.Equal(t, "2021-11-16 06:00:00", payload[0][0])
	require.Equal(t, "06:15:00", payload[0][1])
	require.Equal(t, "medium", payload[0][2])
	require.Equal(t, float64(23), payload[0][3])
}

func TestPull_NoData(t *testing.T) {
	var data = `[]`

	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)

	r := ioutil.NopCloser(bytes.NewReader([]byte(data)))
	res := http.Response{StatusCode: 200, Body: r}

	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/query-stream")
	m.On("Do", mock.Anything).Return(&res, nil)

	_, _, err := kcl.Pull(context.TODO(), ksqldb.QueryOptions{Sql: "select * from bla;"})
	require.NotNil(t, err)
	require.Equal(t, "no result found", err.Error())
}
