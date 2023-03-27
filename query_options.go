package ksqldb

import (
	"strconv"

	"github.com/thmeitz/ksqldb-go/internal"
)

const (
	KSQL_QUERY_PULL_TABLE_SCAN_ENABLED   = "ksql.query.pull.table.scan.enabled"
	KSQL_STREAMS_AUTO_OFFSET_RESET       = "ksql.streams.auto.offset.reset"
	KSQL_IDLE_CONNECTION_TIMEOUT_SECONDS = "ksql.idle.connection.timeout.seconds"
	DEFAULT_IDLE_CONNECTION_TIMEOUT      = int64(600) // 10 minutes
)

type PropertyMap map[string]string

type QueryOptions struct {
	Sql        string      `json:"sql"`
	Properties PropertyMap `json:"properties"`
}

type StreamOffset string

const (
	EARLIEST StreamOffset = "earliest"
	LATEST   StreamOffset = "latest"
)

/*
EnablePullQueryTableScan to control whether table scans are permitted when executing pull queries.

Without this enabled, only key lookups are used.

Enabling table scans removes various restrictions on what types of queries are allowed.

In particular, these pull query types are now permitted:

- No WHERE clause

- Range queries on keys

- Equality and range queries on non-key columns

- Multi-column key queries without specifying all key columns

There may be significant performance implications to using these types of queries,
depending on the size of the data and other workloads running, so use this config carefully.
*/
func (q *QueryOptions) EnablePullQueryTableScan(enable bool) *QueryOptions {
	// check for empty map
	if len(q.Properties) == 0 {
		q.Properties = make(PropertyMap)
	}
	q.Properties[KSQL_QUERY_PULL_TABLE_SCAN_ENABLED] = strconv.FormatBool(enable)
	return q
}

// AutoOffsetReset sets the offset to latest | earliest
//
// Determines what to do when there is no initial offset in Apache Kafka® or
// if the current offset doesn't exist on the server.
// The default value in ksqlDB is `latest`,
// which means all Kafka topics are read from the latest available offset.
func (q *QueryOptions) AutoOffsetReset(offset StreamOffset) *QueryOptions {
	if len(q.Properties) == 0 {
		q.Properties = make(PropertyMap)
	}
	q.Properties[KSQL_STREAMS_AUTO_OFFSET_RESET] = string(offset)
	return q
}

// SetIdleConnectionTimeout sets the timeout for idle connections
//
// A connection is idle if there is no data in either direction on
// that connection for the duration of the timeout.
//
// This configuration can be helpful if you are issuing push queries that only
// receive data infrequently from the server, as otherwise those connections will
// be severed when the timeout (default 10 minutes) is hit.
//
// Decreasing this timeout enables closing connections more aggressively to save
// server resources.
//
// Increasing this timeout makes the server more tolerant of low-data volume use cases.
func (q *QueryOptions) SetIdleConnectionTimeout(seconds int64) *QueryOptions {
	if len(q.Properties) == 0 {
		q.Properties = make(PropertyMap)
	}
	q.Properties[KSQL_IDLE_CONNECTION_TIMEOUT_SECONDS] = strconv.FormatInt(seconds, 10)
	return q
}

// SanitizeQuery removes `\t` and `\n` from the query
func (q *QueryOptions) SanitizeQuery() {
	q.Sql = internal.SanitizeQuery(q.Sql)
}

// EmptyQuery returns true if the query is empty
func (o *QueryOptions) EmptyQuery() bool {
	return len(o.Sql) < 1
}

// NewDefaultPushQueryOptions returns default QueryOptions for push queries
//
// - IdleConnectionTimeout: 600 seconds
// - AutoOffsetReset: "latest"
func NewDefaultPushQueryOptions(sql string) (options QueryOptions) {
	options = QueryOptions{Sql: sql}
	options.
		SetIdleConnectionTimeout(DEFAULT_IDLE_CONNECTION_TIMEOUT).
		AutoOffsetReset(LATEST)
	return
}

// NewDefaultPullQueryOptions returns default QueryOptions for pull queries
//
// - EnablePullQueryTableScan: true
func NewDefaultPullQueryOptions(sql string) (options QueryOptions) {
	options = QueryOptions{Sql: sql}
	options.EnablePullQueryTableScan(true)
	return
}
