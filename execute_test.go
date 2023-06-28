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

func TestExecOptions_SanitizeQuery(t *testing.T) {
	o := ksqldb.ExecOptions{
		KSql: `
		CREATE STREAM IF NOT EXISTS DOGS (ID STRING KEY, 
			NAME STRING, 
			DOGSIZE STRING, 
			AGE STRING) 
		WITH (KAFKA_TOPIC='dogs', 
		VALUE_FORMAT='JSON', PARTITIONS=1);
		`,
	}
	o.SanitizeQuery()
	require.Equal(t, "CREATE STREAM IF NOT EXISTS DOGS (ID STRING KEY, NAME STRING, DOGSIZE STRING, AGE STRING) WITH (KAFKA_TOPIC='dogs', VALUE_FORMAT='JSON', PARTITIONS=1);", o.KSql)
}

func TestExecOption_EmptyQuery(t *testing.T) {
	ctx := context.Background()
	o := ksqldb.ExecOptions{}
	require.True(t, o.EmptyQuery())
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.Execute(ctx, o)
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "empty ksql query", err.Error())
}

func TestExecute_ParseSQLError(t *testing.T) {
	ctx := context.Background()
	m := mocknet.HTTPClient{}
	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(true)
	val, err := kcl.Execute(ctx, ksqldb.ExecOptions{KSql: "create table bla"})
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "1 sql syntax error(s) found", err.Error())
}

func TestExecute_NewKsqlRequest_Error(t *testing.T) {
	ctx := context.Background()
	m := mocknet.HTTPClient{}
	m.Mock.On("BasicAuth", mock.Anything).Return("")
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql")

	m.Mock.
		On("Do", mock.Anything).
		Return(nil, errors.New("error"))

	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(false)
	// error will not found by parser; ";" t the end of statement missing
	val, err := kcl.Execute(ctx, ksqldb.ExecOptions{KSql: "create table bla"})
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "can't do request: error", err.Error())
}

func TestExecute_NewKsqlRequest_ResponseStatusError(t *testing.T) {
	ctx := context.Background()
	r := ioutil.NopCloser(bytes.NewReader([]byte("")))
	res := http.Response{StatusCode: 400, Body: r}
	m := mocknet.HTTPClient{}
	m.Mock.On("BasicAuth", mock.Anything).Return("")
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql")
	m.Mock.
		On("Do", mock.Anything).
		Return(&res, nil)

	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(false)
	val, err := kcl.Execute(ctx, ksqldb.ExecOptions{KSql: "create table bla;"})
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "ksqldb error: unexpected end of JSON input", err.Error())
}

func TestExecute_NewKsqlRequest_ResponseStatusOk_JsonError(t *testing.T) {
	ctx := context.Background()
	r := ioutil.NopCloser(bytes.NewReader([]byte("")))
	res := http.Response{StatusCode: 200, Body: r}
	m := mocknet.HTTPClient{}

	m.Mock.On("BasicAuth", mock.Anything).Return("")
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql")
	m.Mock.
		On("Do", mock.Anything).
		Return(&res, nil)

	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(false)
	val, err := kcl.Execute(ctx, ksqldb.ExecOptions{KSql: "create table bla;"})
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "could not parse the response: unexpected end of JSON input\n", err.Error())
}

func TestExecute_NewKsqlRequest_ResponseStatusOk(t *testing.T) {
	ctx := context.Background()
	r := ioutil.NopCloser(bytes.NewReader([]byte(`[{"@type":"currentStatus","statementText":"CREATE TABLE IF NOT EXISTS DOGS_BY_SIZE WITH (KAFKA_TOPIC='DOGS_BY_SIZE', PARTITIONS=1, REPLICAS=1) AS SELECT\n  DOGS.DOGSIZE DOG_SIZE,\n  COUNT(*) DOGS_CT\nFROM DOGS DOGS\nWINDOW TUMBLING ( SIZE 15 MINUTES ) \nGROUP BY DOGS.DOGSIZE\nEMIT CHANGES;","commandId":"table/DOGS_BY_SIZE/create","commandStatus":{"status":"SUCCESS","message":"Cannot add table DOGS_BY_SIZE: A table with the same name already exists.","queryId":null},"commandSequenceNumber":44,"warnings":[]}]`)))
	res := http.Response{StatusCode: 200, Body: r}
	m := mocknet.HTTPClient{}

	m.Mock.On("BasicAuth", mock.Anything).Return("")
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql")
	m.Mock.
		On("Do", mock.Anything).
		Return(&res, nil)

	kcl, _ := ksqldb.NewClient(&m)
	kcl.EnableParseSQL(false)
	val, err := kcl.Execute(ctx, ksqldb.ExecOptions{KSql: "create table bla;"})
	require.Nil(t, err)
	require.NotNil(t, val)
	require.Equal(t, int64(44), (*val)[0].CommandSequenceNumber)
}
