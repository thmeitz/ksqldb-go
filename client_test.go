package ksqldb_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
)

var options = net.Options{
	Credentials: net.Credentials{Username: "user", Password: "password"},
	BaseUrl:     "http://localhost:8088",
	AllowHTTP:   true,
}

func TestNewClient(t *testing.T) {}

func TestNewClientWithOptions(t *testing.T) {

	//factory := new(mocks.KsqldbFactory)
	//kcl := factory.On("NewClientWithOptions", options).Return(nil, nil)
	//fmt.Println(kcl.ReturnArguments)

}

func TestClient_EnableParseSQL(t *testing.T) {

	kcl, err := ksqldb.NewClientWithOptions(options)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, kcl.ParseSQLEnabled())
	kcl.EnableParseSQL(false)
	assert.False(t, kcl.ParseSQLEnabled())
}
func TestClient_ParseSQLEnabled(t *testing.T) {}
