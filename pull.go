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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thmeitz/ksqldb-go/parser"
)

// Pull queries are like "traditional" RDBMS queries in which
// the query terminates once the state has been queried.
//
// To use this function pass in the the SQL query statement, and
// a boolean for whether full table scans should be enabled.
//
// The function returns a ksqldb.Header and ksqldb.Payload
// which will hold one or more rows of data. You will need to
// define variables to hold each column's value. You can adopt
// this pattern to do this:
//
//	var col1 string
//	var col2 float64
//	for _, row := range r {
//		col1 = row[0].(string)
//		col2 = row[1].(float64)
//		... Do other stuff with the data here
//		}
//	}
func (api *KsqldbClient) Pull(ctx context.Context, options QueryOptions) (header Header, payload Payload, err error) {

	if options.EmptyQuery() {
		return header, payload, fmt.Errorf("empty ksql query")
	}

	// remove \t \n from query
	options.SanitizeQuery()

	if api.ParseSQLEnabled() {
		ksqlerr := parser.ParseSql(options.Sql)
		if ksqlerr != nil {
			return header, payload, ksqlerr
		}
	}

	jsonData, err := json.Marshal(options)
	if err != nil {
		return header, payload, fmt.Errorf("can't marshal input data")
	}

	// Create the request
	req, err := newQueryStreamRequest(api.http, ctx, bytes.NewReader(jsonData))
	if err != nil {
		return header, payload, fmt.Errorf("can't create new request with context: %w", err)
	}
	req.Header.Add("Accept", "application/json; charset=utf-8")

	res, err := api.http.Do(req)
	if err != nil {
		return header, payload, fmt.Errorf("can't do request: %+w", err)
	}
	defer func() {
		berr := res.Body.Close()
		if err == nil {
			err = berr
		}
	}()

	body, err := api.readBody(res.Body)
	if err != nil {
		return header, payload, fmt.Errorf("can't read response body:\n%w", err)
	}

	if res.StatusCode != http.StatusOK {
		return header, payload, handleRequestError(res.StatusCode, body)
	}

	var result []interface{}
	// Parse the output
	if err := json.Unmarshal(body, &result); err != nil {
		return header, payload, fmt.Errorf("could not parse the response:\n%w", err)
	}

	if len(result) == 0 {
		return header, payload, fmt.Errorf("%w", ErrNotFound)
	}

	for _, resultSet := range result {
		switch resultSetTypes := resultSet.(type) {
		case map[string]interface{}:
			// It's the Header
			header = processHeader(resultSetTypes)
		case []interface{}:
			// It's a row of data
			payload = append(payload, resultSetTypes)
		}
	}

	return
}

func processHeader(data map[string]interface{}) (header Header) {
	if _, ok := data["queryId"].(string); ok {
		header.QueryId = data["queryId"].(string)
	}

	names, okNames := data["columnNames"].([]interface{})
	types, okTypes := data["columnTypes"].([]interface{})
	if okNames && okTypes {
		for col := range names {
			if colName, ok := names[col].(string); colName != "" && ok {
				if colType, ok := types[col].(string); colType != "" && ok {
					header.Columns = append(header.Columns, Column{Name: colName, Type: colType})
				}
			}
		}
	}
	return
}
