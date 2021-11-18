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
	"strconv"

	"github.com/thmeitz/ksqldb-go/internal"
	"github.com/thmeitz/ksqldb-go/parser"
)

const (
	KSQL_QUERY_PULL_TABLE_SCAN_ENABLED = "ksql.query.pull.table.scan.enabled"
)

type QueryOptions struct {
	Sql        string      `json:"sql"`
	Properties PropertyMap `json:"properties"`
}

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

func (q *QueryOptions) SanitizeQuery() {
	q.Sql = internal.SanitizeQuery(q.Sql)
}

func (o *QueryOptions) EmptyQuery() bool {
	return len(o.Sql) < 1
}

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
// 		var col1 string
// 		var col2 float64
// 		for _, row := range r {
// 			col1 = row[0].(string)
// 			col2 = row[1].(float64)
// 			// Do other stuff with the data here
// 			}
// 		}
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
	defer res.Body.Close()

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

	switch len(result) {
	case 0:
		return header, payload, fmt.Errorf("%w (not even a header row) returned from lookup. Maybe we got an error:%v", ErrNotFound, err)
	case 1:
		// len 1 means we just got a header, no rows
		// Should we define our own error types here so we can return more clearly
		// an indicator that no rows were found?
		// ANSWER: no - maybe we have no data - its not an error
		return header, payload, ErrNotFound
	default:
		for _, z := range result {
			switch zz := z.(type) {
			case map[string]interface{}:
				// It's a header row, so extract the data
				// {"queryId":null,"columnNames":["WINDOW_START","WINDOW_END","DOG_SIZE","DOGS_CT"],"columnTypes":["STRING","STRING","STRING","BIGINT"]}
				if _, ok := zz["queryId"].(string); ok {
					header.queryId = zz["queryId"].(string)
				} //else {
				// api.logger.Info("(query id not found - this is expected for a pull query)")
				// why should we log this???? - check facts in java source code
				//}

				names, okn := zz["columnNames"].([]interface{})
				types, okt := zz["columnTypes"].([]interface{})
				if okn && okt {
					for col := range names {
						if n, ok := names[col].(string); n != "" && ok {
							if t, ok := types[col].(string); t != "" && ok {
								a := Column{Name: n, Type: t}
								header.columns = append(header.columns, a)

							} /*else {
								// api.logger.Infof("nil type found for column %v", col)
							}*/
						} /*else {
							// api.logger.Infof("nil name found for column %v", col)
						}*/
					}
				} /*else {
					// api.logger.Infof("column names/types not found in header:\n%v", zz)
				}*/

			case []interface{}:
				// It's a row of data
				payload = append(payload, zz)
			}
		}

		return header, payload, nil
	}
}
