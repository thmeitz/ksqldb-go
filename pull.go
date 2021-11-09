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

Parts of this apiclient are borrowed from Zalando Skipper
https://github.com/zalando/skipper/blob/master/net/httpclient.go

Zalando licence: MIT
https://github.com/zalando/skipper/blob/master/LICENSE
*/

package ksqldb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
// 		var COL1 string
// 		var COL2 float64
// 		for _, row := range r {
// 			COL1 = row[0].(string)
// 			COL2 = row[1].(float64)
// 			// Do other stuff with the data here
// 			}
// 		}
func Pull(api *Client, ctx context.Context, q string, s bool) (h Header, r Payload, err error) {

	// first sanitize the query
	query := api.SanitizeQuery(q)
	// we're kick in our ksqlparser to check the query string
	ksqlerr := ParseKSQL(query)
	if ksqlerr != nil {
		return h, r, ksqlerr
	}

	// Create the request
	payload := strings.NewReader(`{"properties":{"ksql.query.pull.table.scan.enabled": ` + strconv.FormatBool(s) + `},"sql":"` + query + `"}`)

	req, err := api.NewQueryStreamRequest(ctx, payload)
	if err != nil {
		return h, r, fmt.Errorf("can't create new request with context: %w", err)
	}
	req.Header.Add("Accept", "application/json; charset=utf-8")

	res, err := (&api.client).Do(req)
	if err != nil {
		return h, r, fmt.Errorf("can't do request: %+w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return h, r, fmt.Errorf("can't read response body:\n%w", err)
	}

	if res.StatusCode != http.StatusOK {
		return h, r, api.handleRequestError(res.StatusCode, body)
	}

	var x []interface{}
	// Parse the output
	if err := json.Unmarshal(body, &x); err != nil {
		return h, r, fmt.Errorf("could not parse the response as json:\n%w", err)

	}

	switch len(x) {
	case 0:
		return h, r, fmt.Errorf("%w (not even a header row) returned from lookup. Maybe we got an error:%v", ErrNotFound, err)
	case 1:
		// len 1 means we just got a header, no rows
		// Should we define our own error types here so we can return more clearly
		// an indicator that no rows were found?
		// ANSWER: yes
		return h, r, ErrNotFound
	default:
		for _, z := range x {
			switch zz := z.(type) {
			case map[string]interface{}:
				// It's a header row, so extract the data
				// {"queryId":null,"columnNames":["WINDOW_START","WINDOW_END","DOG_SIZE","DOGS_CT"],"columnTypes":["STRING","STRING","STRING","BIGINT"]}
				if _, ok := zz["queryId"].(string); ok {
					h.queryId = zz["queryId"].(string)
				} else {
					// it is a hard fact, so we should throw an error?
					// log interface needs a format and a interface{}
					api.logger.Info("(query id not found - this is expected for a pull query)")
				}

				names, okn := zz["columnNames"].([]interface{})
				types, okt := zz["columnTypes"].([]interface{})
				if okn && okt {
					for col := range names {
						if n, ok := names[col].(string); n != "" && ok {
							if t, ok := types[col].(string); t != "" && ok {
								a := Column{Name: n, Type: t}
								h.columns = append(h.columns, a)

							} else {
								api.logger.Infof("nil type found for column %v", col)
							}
						} else {
							api.logger.Infof("nil name found for column %v", col)
						}
					}
				} else {
					api.logger.Infof("column names/types not found in header:\n%v", zz)
				}

			case []interface{}:
				// It's a row of data
				r = append(r, zz)
			}
		}

		return h, r, nil
	}
}
