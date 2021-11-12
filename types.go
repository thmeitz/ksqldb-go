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
	"io"
	"net/http"

	"github.com/thmeitz/ksqldb-go/net"
)

type KsqldbFactory interface {
	// NewClient factory
	NewClient(net.HTTPClient) (*KsqldbClient, error)

	// NewClientWithOptions factory
	NewClientWithOptions(options net.Options) (*KsqldbClient, error)
}

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
	//ValidateProperty(string) error

	//
	Pull(context.Context, string, bool) (Header, Payload, error)

	//
	Push(context.Context, string, chan<- Row, chan<- Header) error

	//
	EnableParseSQL(bool)

	ParseSQLEnabled() bool

	// Close closes net.HTTPClient transport
	Close()
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
