/*
Copyright © 2021 Robin Moffat & Contributors
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

package ksqldb

import (
	"context"

	"github.com/thmeitz/ksqldb-go/net"
)

type NewClientFactory interface {
	// NewClient factory
	NewClient(net.HTTPClient) (*KsqldbClient, error)
}

type NewClientWithOptionsFactory interface {
	// NewClientWithOptions factory
	NewClientWithOptions(options net.Options) (*Ksqldb, error)
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
	TerminateCluster(topics ...string) (*KsqlResponseSlice, error)

	// ValidateProperty validates a property
	// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/is_valid_property-endpoint/
	ValidateProperty(property string) (*bool, error)

	//
	Pull(context.Context, string, bool) (Header, Payload, error)

	//
	Push(context.Context, string, chan<- Row, chan<- Header) error

	// EnableParseSQL enables/disables query parsing for push/pull/execute requests
	EnableParseSQL(bool)

	// ParseSQLEnabled returns true if query parsing is enabled or not
	ParseSQLEnabled() bool

	// Close closes net.HTTPClient transport
	Close()
}

//type ksqldbRequest interface {
// 	// newQueryRequest interface
// 	// @API https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/query-endpoint/
// 	newQueryRequest(http.Client, io.Reader) (*http.Request, error)
//
// 	newKsqlRequest(api net.HTTPClient, payload io.Reader) (*http.Request, error)
//
// 	// newQueryStreamRequest interface
// 	// @API https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/streaming-endpoint/
// 	newQueryStreamRequest(api net.HTTPClient, ctx context.Context, payload io.Reader) (*http.Request, error)
//
// 	newCloseQueryRequest(api net.HTTPClient, ctx context.Context, payload io.Reader) (*http.Request, error)
//
// 	handleRequestError(code int, buf []byte) error
//
//handleGetRequest(httpClient *net.HTTPClient, url string) (*[]byte, error)
//
// 	newPostRequest(api net.HTTPClient, ctx context.Context, endpoint string, payload io.Reader) (*http.Request, error)
//
// 	// NewIntrospectQueryRequest
// 	// @API https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/status-endpoint/
//}

type SessionVariablesMap map[string]interface{}

type PropertyMap map[string]string

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
