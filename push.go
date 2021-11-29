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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thmeitz/ksqldb-go/parser"
)

// Push queries are continuous queries in which new events
// or changes to a table's state are pushed to the client.
// You can think of them as subscribing to a stream of changes.
//
// Since push queries never end, this function expects a channel
// to which it can write new rows of data as and when they are
// received.
//
// To use this function pass in a context, the SQL query statement,
// and two channels:
//
// * ksqldb.Row - rows of data
// * ksqldb.Header - header (including column definitions).
//
// If you don't want to block before receiving
// row data then make this channel buffered.
//
// The channel is populated with ksqldb.Row which represents
// one row of data. You will need to define variables to hold
// each column's value. You can adopt this pattern to do this:
// 		var DATA_TS float64
// 		var ID string
// 		for row := range rc {
// 			if row != nil {
//				DATA_TS = row[0].(float64)
// 				ID = row[1].(string)
func (api *KsqldbClient) Push(ctx context.Context, options QueryOptions,
	rowChannel chan<- Row, headerChannel chan<- Header) (err error) {
	if options.EmptyQuery() {
		return fmt.Errorf("empty ksql query")
	}

	// remove \t \n from query
	options.SanitizeQuery()

	if api.ParseSQLEnabled() {
		ksqlerr := parser.ParseSql(options.Sql)
		if ksqlerr != nil {
			return ksqlerr
		}
	}

	jsonData, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("can't marshal input data")
	}

	req, err := newQueryStreamRequest(api.http, ctx, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating new request with context: %w", err)
	}

	//  make the request
	res, err := api.http.Do(req)

	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)

	doThis := true
	var row interface{}
	var header Header

	for doThis {
		select {
		case <-ctx.Done():
			// close the channels and terminate the loop regardless
			defer close(rowChannel)
			defer close(headerChannel)
			defer func() { doThis = false }()
			if err := api.ClosePushQuery(ctx, header.QueryId); err != nil {
				return fmt.Errorf("%w", err)
			}
		default:
			// Read the next chunk
			body, err := reader.ReadBytes('\n')
			if err != nil {
				doThis = false
			}
			if res.StatusCode != http.StatusOK {
				return handleRequestError(res.StatusCode, body)
			}

			if len(body) > 0 {
				// Parse the output
				if err := json.Unmarshal(body, &row); err != nil {
					return fmt.Errorf("could not parse the response: %w\n%v", err, string(body))
				}

				switch resultSetTypes := row.(type) {
				case map[string]interface{}:
					headerChannel <- processHeader(resultSetTypes)

				case []interface{}:
					// It's a row of data
					rowChannel <- resultSetTypes
				}
			}
		}
	}
	return nil
}
