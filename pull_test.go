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

// var data = `[
// 	{
// 		"queryId":null,
// 		"columnNames":[
// 			"WINDOW_START","WINDOW_END","DOG_SIZE","DOGS_CT"
// 		],
// 		"columnTypes":[
// 			"STRING","STRING","STRING","BIGINT"
// 		]
// 	},
// 	["2021-11-16 06:00:00","06:15:00","medium",23],
// 	["2021-11-16 06:15:00","06:30:00","medium",250],
// 	["2021-11-16 06:30:00","06:45:00","medium",234],
// 	["2021-11-16 06:45:00","07:00:00","medium",242],
// 	["2021-11-16 07:00:00","07:15:00","medium",228],
// 	["2021-11-16 07:15:00","07:30:00","medium",242],
// 	["2021-11-16 07:30:00","07:45:00","medium",226],
// 	["2021-11-16 07:45:00","08:00:00","medium",222],
// 	["2021-11-16 08:00:00","08:15:00","medium",216],
// 	["2021-11-16 08:15:00","08:30:00","medium",219]
// ]
// `

func TestQueryOptions_SanitizeQuery(t *testing.T) {
	o := ksqldb.QueryOptions{Sql: `select * 
	from bla`}
	o.EnablePullQueryTableScan(true)
	require.Equal(t, "true", o.Properties[ksqldb.KSQL_QUERY_PULL_TABLE_SCAN_ENABLED])
	o.SanitizeQuery()
	require.Equal(t, "select * from bla", o.Sql)
}

func TestQueryOptions_TestEmptyQuery(t *testing.T) {
	o := ksqldb.QueryOptions{Sql: ""}
	require.True(t, o.EmptyQuery())
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	_, _, err := kcl.Pull(context.TODO(), o)
	require.NotNil(t, err)
	require.Equal(t, "empty ksql query", err.Error())
}

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
