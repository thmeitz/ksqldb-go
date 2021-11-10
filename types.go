/*
Copyright © 2021 Robin Moffat & Contributors
Copyright © 2021 Thomas Meitz <thme219@gmail.com>

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

package ksqldb

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/thmeitz/ksqldb-go/net"
)

type Ksqldb interface {
	// GetServerInfo returns informations about the ksqlDB Server
	// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/info-endpoint/
	GetServerInfo() (*KsqlServerInfo, error)

	// GetServerStatus returns server status
	// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/info-endpoint/
	GetServerStatus() (*ServerStatusResponse, error)

	// GetClusterStatus
	// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/cluster-status-endpoint/
	GetClusterStatus() (*ClusterStatusResponse, error)

	// TerminateCluster terminates a ksqldb cluster - READ THE DOCS before you call this endpoint
	// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/terminate-endpoint/
	//TerminateCluster(*TopicList) error

	// ValidateProperty validates a property
	// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/is_valid_property-endpoint/
	//ValidateProperty(http.Client, string) error
}

type KsqldbClient struct {
	http *net.HTTPClient
}

func NewClient(http net.HTTPClient) (*KsqldbClient, error) {
	var client = KsqldbClient{
		http: &http,
	}

	return &client, nil
}

func NewClientWithOptions(options net.Options) (*KsqldbClient, error) {
	http, err := net.NewHTTPClient(options, nil)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return NewClient(http)
}

type ksqldbRequest interface {
	// newQueryRequest interface
	// @API https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/query-endpoint/
	newQueryRequest(http.Client, io.Reader) (*http.Request, error)

	// newQueryStreamRequest interface
	// @API https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/streaming-endpoint/
	newQueryStreamRequest(http.Client, context.Context, io.Reader) (*http.Request, error)

	// NewIntrospectQueryRequest
	// @API https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/status-endpoint/
}

type KsqlParser interface {
	ParseSql(string) []error
}

type TopicList []string

// KsqlServerInfo
// @
type KsqlServerInfo struct {
	Version        string `json:"version"`
	KafkaClusterID string `json:"kafkaClusterId"`
	KsqlServiceID  string `json:"ksqlServiceId"`
}

// KsqlServerInfoResponse
type KsqlServerInfoResponse struct {
	KsqlServerInfo KsqlServerInfo `json:"KsqlServerInfo"`
}

// Row represents a row returned from a query
type Row []interface{}

// Payload represents multiple rows
type Payload []Row

// Header represents a header returned from a query
type Header struct {
	queryId string
	columns []Column
}

// Column represents the metadata for a column in a Row
type Column struct {
	Name string
	Type string
}
