package ksqldb_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

var options = net.Options{
	Credentials: net.Credentials{Username: "user", Password: "password"},
	BaseUrl:     "http://localhost:8088",
	AllowHTTP:   true,
}

func TestNewClientWithOptions_Error(t *testing.T) {
	var opt1 = net.Options{
		Credentials: net.Credentials{Username: "user", Password: "password"},
		BaseUrl:     "fa",
		AllowHTTP:   true,
	}

	_, err := ksqldb.NewClientWithOptions(opt1)

	require.NotNil(t, err)
	require.Equal(t, "invalid host name given", err.Error())

	// Ensures that the Ksqldb interface is implemented.
	// Aborts the compiler if it does not.
	// var _ ksqldb.Ksqldb = kcl
}

func TestClient_EnableParseSQL(t *testing.T) {
	kcl, err := ksqldb.NewClientWithOptions(options)
	require.Nil(t, err)
	require.True(t, kcl.ParseSQLEnabled())
	kcl.EnableParseSQL(false)
	require.False(t, kcl.ParseSQLEnabled())
}
