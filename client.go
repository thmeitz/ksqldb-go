package ksqldb

import (
	"fmt"

	"github.com/thmeitz/ksqldb-go/net"
)

type KsqldbClient struct {
	http     *net.HTTPClient
	parseSQL bool
}

func NewClient(http net.HTTPClient) (*KsqldbClient, error) {
	var client = KsqldbClient{
		http:     &http,
		parseSQL: true,
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

// EnableParseSQL enables / disables sql parsing
func (cl *KsqldbClient) EnableParseSQL(activate bool) {
	cl.parseSQL = activate
}

// ParseSQLEnabled returns true if sql parsing is enabled; false otherwise
func (cl *KsqldbClient) ParseSQLEnabled() bool {
	return cl.parseSQL
}

//
func (cl *KsqldbClient) Close() {
	(*cl.http).Close()
}
