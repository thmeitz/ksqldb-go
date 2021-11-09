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
)

type KSqlDB interface {
	GetServerInfo(http.Client) (*ServerInfo, error)
	Healthcheck(http.Client, string) (*ServerHealthResponse, error)
	NewKsqlRequest(http.Client, io.Reader) (*http.Request, error)
	NewQueryStreamRequest(http.Client, context.Context, io.Reader) (*http.Request, error)
}

type ServerHealthResponse struct {
	IsHealthy *bool `json:"isHealthy"`
	Details   struct {
		Metastore struct {
			IsHealthy *bool `json:"isHealthy"`
		} `json:"metastore"`
		Kafka struct {
			IsHealthy *bool `json:"isHealthy"`
		} `json:"kafka"`
	} `json:"details"`
	KSQLServiceID string `json:"ksqlServiceId"`
}

type ServerInfo struct {
	Version        string `json:"version"`
	KafkaClusterID string `json:"kafkaClusterId"`
	KSQLServiceID  string `json:"ksqlServiceId"`
}

type ServerInfoResponse struct {
	KSQLServerInfo ServerInfo `json:"KsqlServerInfo"`
}

// Stuff to rework
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
