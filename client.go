package ksqldb

import (
	"fmt"

	"github.com/thmeitz/ksqldb-go/net"
)

type KsqldbClient struct {
	http     *net.HTTPClient
	parseSQL bool
}

// NewClient returns a new KsqldbClient with the given net.HTTPclient
func NewClient(http net.HTTPClient) (*KsqldbClient, error) {
	var client = KsqldbClient{
		http:     &http,
		parseSQL: true,
	}

	return &client, nil
}

// NewClientWithOptions returns a new @KsqldbClient with Options
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

// Close closes the underlying http transport
func (cl *KsqldbClient) Close() {
	(*cl.http).Close()
}
