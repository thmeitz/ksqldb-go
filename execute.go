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

Parts of this apiclient are borrowed from Zalando Skipper
https://github.com/zalando/skipper/blob/master/net/httpclient.go

Zalando licence: MIT
https://github.com/zalando/skipper/blob/master/LICENSE
*/

package ksqldb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thmeitz/ksqldb-go/internal"
	"github.com/thmeitz/ksqldb-go/parser"
)

type SessionVariablesMap map[string]interface{}

type ExecOptions struct {
	KSql                  string              `json:"ksql"`
	StreamsProperties     PropertyMap         `json:"streamsProperties,omitempty"`
	SessionVariables      SessionVariablesMap `json:"sessionVariables,omitempty"`
	CommandSequenceNumber int64               `json:"commandSequenceNumber,omitempty"`
}

func (o *ExecOptions) SanitizeQuery() {
	o.KSql = internal.SanitizeQuery(o.KSql)
}

func (o *ExecOptions) EmptyQuery() bool {
	return len(o.KSql) < 1
}

// Execute will execute a ksqlDB statement.
// All statements, except those starting with SELECT,
// can be run on this endpoint.
// To run SELECT statements use use Push or Pull functions.
//
// To use this function pass in the @ExecOptions.
//
// Ref: https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/ksql-endpoint/
func (api *KsqldbClient) Execute(options ExecOptions) (response *KsqlResponseSlice, err error) {
	response = new(KsqlResponseSlice)

	if options.EmptyQuery() {
		return nil, fmt.Errorf("empty ksql query")
	}
	// remove \t \n from query
	options.SanitizeQuery()

	if api.ParseSQLEnabled() {
		ksqlerr := parser.ParseSql(options.KSql)
		if ksqlerr != nil {
			return nil, ksqlerr
		}
	}

	jsonData, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("can't marshal input data")
	}

	// make the request
	req, err := newKsqlRequest(api.http, bytes.NewReader(jsonData))
	// api.logger.Debugf("sending ksqlDB request:%v", q)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %w", err)
	}

	res, err := api.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer func() {
		berr := res.Body.Close()
		if err == nil {
			err = berr
		}
	}()

	body, err := api.readBody(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	// this is only one side of the coin
	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res.StatusCode, body)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("could not parse the response: %w\n%v", err, string(body))
	}

	return
}
