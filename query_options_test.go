package ksqldb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
	mocknet "github.com/thmeitz/ksqldb-go/mocks/net"
)

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

func TestAutoOffsetReset(t *testing.T) {
	o := ksqldb.QueryOptions{Sql: `select * 
	from bla`}
	o.AutoOffsetReset(ksqldb.EARLIEST)
	require.Equal(t, "earliest", o.Properties[ksqldb.KSQL_STREAMS_AUTO_OFFSET_RESET])
	o.AutoOffsetReset(ksqldb.LATEST)
	require.Equal(t, "latest", o.Properties[ksqldb.KSQL_STREAMS_AUTO_OFFSET_RESET])
}

func TestSetIdleConnectionTimeout(t *testing.T) {
	timeout := int64(10000)
	o := ksqldb.QueryOptions{Sql: `select * 
	from bla`}
	o.SetIdleConnectionTimeout(timeout)
	require.Equal(t, "10000", o.Properties[ksqldb.KSQL_IDLE_CONNECTION_TIMEOUT_SECONDS])
}

func TestNewDefaultPushQueryOptions(t *testing.T) {
	o := ksqldb.NewDefaultPushQueryOptions("select * from bla")
	require.Equal(t, "600", o.Properties[ksqldb.KSQL_IDLE_CONNECTION_TIMEOUT_SECONDS])
	require.Equal(t, "latest", o.Properties[ksqldb.KSQL_STREAMS_AUTO_OFFSET_RESET])
}

func TestNewDefaultPullQueryOptions(t *testing.T) {
	o := ksqldb.NewDefaultPullQueryOptions("select * from bla")
	require.Equal(t, "true", o.Properties[ksqldb.KSQL_QUERY_PULL_TABLE_SCAN_ENABLED])
	o.EnablePullQueryTableScan(false)
	require.Equal(t, "false", o.Properties[ksqldb.KSQL_QUERY_PULL_TABLE_SCAN_ENABLED])
}
