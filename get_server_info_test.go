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

func TestGetServerInfo_ResponseError(t *testing.T) {
	ctx := context.Background()

	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/info")
	m.Mock.
		On("Get", ctx, mock.Anything).
		Return(nil, errors.New("error"))

	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.GetServerInfo(ctx)
	require.Nil(t, val)
	require.NotNil(t, err)
	require.Equal(t, "can't get server info: error", err.Error())
}

//

func TestGetServerInfo_SuccessfullResponse(t *testing.T) {
	ctx := context.Background()
	json := `{"KsqlServerInfo":{"version":"0.21.0","kafkaClusterId":"kgqdUfEoTBSutJd1JWHIyQ","ksqlServiceId":"confluent_rmoff_01","serverStatus":"RUNNING"}}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	res := http.Response{StatusCode: 200, Body: r}
	m := mocknet.HTTPClient{}
	m.Mock.On("GetUrl", mock.Anything).Return("http://localhost/info")
	m.Mock.
		On("Get", ctx, mock.Anything).
		Return(&res, nil)

	kcl, _ := ksqldb.NewClient(&m)
	val, err := kcl.GetServerInfo(ctx)
	require.Nil(t, err)
	require.NotNil(t, val)
	require.Equal(t, "0.21.0", val.Version)
	require.Equal(t, "kgqdUfEoTBSutJd1JWHIyQ", val.KafkaClusterID)
	require.Equal(t, "confluent_rmoff_01", val.KsqlServiceID)
	require.Equal(t, "RUNNING", val.ServerStatus)
}
