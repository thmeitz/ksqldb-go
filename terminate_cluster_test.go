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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
	mocknet "github.com/thmeitz/ksqldb-go/mocks/net"
)

func TestTerminateClusterTopics_Add(t *testing.T) {
	topicList := ksqldb.TerminateClusterTopics{}
	topicList.Add("FOO", "bar.*")
	require.Equal(t, 2, topicList.Size())
	require.Equal(t, "FOO", topicList.DeleteTopicList[0])
	require.Equal(t, "bar.*", topicList.DeleteTopicList[1])
}

func TestTerminateCluster_WithoutTopicsPostError(t *testing.T) {
	tpc := ksqldb.TerminateClusterTopics{}
	b, _ := json.Marshal(&tpc)
	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql/terminate")
	m.Mock.
		On("Post", mock.Anything, "application/vnd.ksql.v1+json", bytes.NewBuffer(b)).
		Return(nil, errors.New("error"))
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.TerminateCluster()
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "error", err.Error())
}

func TestTerminateCluster_WithTopicsUnmarshalError(t *testing.T) {
	tpc := ksqldb.TerminateClusterTopics{}
	tpc.Add("test", "test2")
	b, _ := json.Marshal(&tpc)

	json := `{"name":"Test Name"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}

	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql/terminate")
	m.Mock.
		On("Post", mock.Anything, "application/vnd.ksql.v1+json", bytes.NewBuffer(b)).
		Return(&res, nil)
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.TerminateCluster("test", "test2")
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "could not parse the response:json: cannot unmarshal object into Go value of type ksqldb.KsqlResponseSlice", err.Error())
}

func TestTerminateCluster_HttpStatusNotOk(t *testing.T) {
	tpc := ksqldb.TerminateClusterTopics{}
	tpc.Add("test", "test2")
	b, _ := json.Marshal(&tpc)

	json := `{"@type":"some error", "error_code": 4000, "message": "some message"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 400, Body: r}

	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql/terminate")
	m.Mock.
		On("Post", mock.Anything, "application/vnd.ksql.v1+json", bytes.NewBuffer(b)).
		Return(&res, nil)
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.TerminateCluster("test", "test2")
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "some message", err.Error())
}

func TestTerminateCluster_EmptyResponseBody(t *testing.T) {
	tpc := ksqldb.TerminateClusterTopics{}
	tpc.Add("test", "test2")
	b, _ := json.Marshal(&tpc)

	r := ioutil.NopCloser(bytes.NewReader([]byte("")))
	res := http.Response{StatusCode: 400, Body: r}

	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql/terminate")
	m.Mock.
		On("Post", mock.Anything, "application/vnd.ksql.v1+json", bytes.NewBuffer(b)).
		Return(&res, nil)
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.TerminateCluster("test", "test2")
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "ksqldb error: unexpected end of JSON input", err.Error())
}

func TestTerminateCluster_TerminatedCluster(t *testing.T) {
	tpc := ksqldb.TerminateClusterTopics{}
	tpc.Add("test", "test2")
	b, _ := json.Marshal(&tpc)

	json := `[{"@type":"currentStatus","statementText":"TERMINATE CLUSTER;","commandId":"terminate/CLUSTER/execute","commandStatus":{"status":"QUEUED","message":"Statement written to command topic","queryId":null},"commandSequenceNumber":4,"warnings":[]}]`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}

	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/ksql/terminate")
	m.Mock.
		On("Post", mock.Anything, "application/vnd.ksql.v1+json", bytes.NewBuffer(b)).
		Return(&res, nil)
	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.TerminateCluster("test", "test2")
	require.NotNil(t, val)
	require.Nil(t, err)
}
