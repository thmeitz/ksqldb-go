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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/log-go"
)

const (
	HEARTBEAT_TRESHOLD = 9 // After 9 minutes the connection will be closed
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
//                  If you don't want to block before receiving
//				    row data then make this channel buffered.
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
func Push(cl *Client, ctx context.Context, q string, rc chan<- Row, hc chan<- Header) (err error) {

	// first sanitize the query
	query := cl.SanitizeQuery(q)
	// we're kick in our ksqlparser to check the query string
	ksqlerr := cl.ParseKSQL(query)
	if ksqlerr != nil {
		return ksqlerr
	}

	// https://docs.confluent.io/5.0.4/ksql/docs/installation/server-config/config-reference.html#ksql-streams-auto-offset-reset
	payload := strings.NewReader(`{"properties":{"ksql.streams.auto.offset.reset": "latest"},"sql":"` + query + `"}`)

	req, err := cl.newQueryStreamRequest(ctx, payload)
	if err != nil {
		return fmt.Errorf("error creating new request with context: %w", err)
	}

	// TODO If we've got creds to pass, let's pass them
	// if cl.username != "" {
	// 	req.SetBasicAuth(cl.username, cl.password)
	// }

	// don't know if we are needing this stuff in the new client
	go cl.heartbeat(&cl.client, &ctx)

	//  make the request
	cl.logger.Debugf("sending ksqlDB request:\n\t%v", q)
	res, err := (&cl.client).Do(req)

	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
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
			defer close(rc)
			defer close(hc)
			defer func() { doThis = false }()
			// Try to close the query
			payload := strings.NewReader(`{"queryId":"` + header.queryId + `"}`)
			// cl.log("payload: %v", *payload)
			req, err := cl.newCloseQueryRequest(ctx, payload)

			cl.logger.Debugf("closing ksqlDB query\t%v", header.queryId)
			if err != nil {
				return fmt.Errorf("failed to construct http request to cancel query\n%w", err)
			}

			res, err := cl.client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to execute http request to cancel query\n%w", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("close query failed:\n%v", res)
			}
			cl.logger.Info("query closed.")
		default:

			// Read the next chunk
			body, err := reader.ReadBytes('\n')
			if err != nil {
				doThis = false
			}
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("the http request did not return a success code:\n%v / %v", res.StatusCode, string(body))
			}

			if len(body) > 0 {
				// Parse the output
				if err := json.Unmarshal(body, &row); err != nil {
					return fmt.Errorf("could not parse the response as JSON: %w\n%v", err, string(body))
				}

				switch zz := row.(type) {
				case map[string]interface{}:
					// It's a header row, so extract the data
					// {"queryId":null,"columnNames":["WINDOW_START","WINDOW_END","DOG_SIZE","DOGS_CT"],"columnTypes":["STRING","STRING","STRING","BIGINT"]}
					if _, ok := zz["queryId"].(string); ok {
						header.queryId = zz["queryId"].(string)
					} else {
						cl.logger.Debug("query id not found - this is expected for a pull query")
					}

					names, okn := zz["columnNames"].([]interface{})
					types, okt := zz["columnTypes"].([]interface{})
					if okn && okt {
						for col := range names {
							if n, ok := names[col].(string); n != "" && ok {
								if t, ok := types[col].(string); t != "" && ok {
									a := Column{Name: n, Type: t}
									header.columns = append(header.columns, a)

								} else {
									cl.logger.Infof("nil type found for column %v", col)
								}
							} else {
								cl.logger.Infof("Nil name found for column %v", col)
							}
						}
					} else {
						cl.logger.Infof("Column names/types not found in header:\n%v", zz)
					}
					cl.logger.Debugf("Header: %v", header)
					hc <- header

				case []interface{}:
					// It's a row of data
					cl.logger.Debugf("Row: %v", zz)
					rc <- zz
				}
			}
		}
	}
	return nil
}

// heartbeat sends a heartbeat to the server
//
// The default for KSQL server is a 10 minute timeout, which is a problem on low volume connections.
// `heartbeat` must be used on a go routine like this `go cl.heartbeat(*client, ctx)`
//
// This fixes issuue #17 by adding a gorountine which lists the streams every minute to keep the connection alive.
// If we miss 9 heartbeats (9 minutes), then close the connection since KSQL Server only keeps it alive for 10 minutes by default.
//
// TODO: we have introduced a healthcheck endpoint, we should refactor this if the http.Client lives in the ksqldb.Client{}
// and reusing http connections
func (cl *Client) heartbeat(client *http.Client, ctx *context.Context) {
	missedHeartbeat := 0
	heartbeatThreshold := HEARTBEAT_TRESHOLD // Default for KSQL Server is close connection after 10 minutes of no activity
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		cl.logger.Info("sending heartbeat...")

		pingPayload := strings.NewReader(`{"ksql":"SHOW STREAMS;"}`)
		pingReq, err := cl.newKsqlRequest(pingPayload)
		cl.logger.Debugf("sending ksqlDB request:\n\t%v", pingPayload)
		if err != nil {
			missedHeartbeat += 1
			cl.logger.Errorf("Couldn't create new HTTP request, %s", err)
		} else {

			res, err := client.Do(pingReq)
			if err != nil {
				missedHeartbeat += 1
				cl.logger.Errorf("failed to send heartbeat: %v", res.StatusCode)
			} else {

				bodyBytes, err := ioutil.ReadAll(res.Body)
				if err != nil {
					missedHeartbeat += 1
					cl.logger.Errorw("failed to read heartbeat body", log.Fields{"status": res.StatusCode})
				} else {
					res.Body.Close()

					body := string(bodyBytes)

					if res.StatusCode != 200 {
						missedHeartbeat += 1
						cl.logger.Debugw("the heartbeat did not return a success code", log.Fields{"status": res.StatusCode, "body": string(body)})
					} else {
						missedHeartbeat = 0
						cl.logger.Info("got heartbeat")
					}
				}
			}
		}

		if missedHeartbeat == heartbeatThreshold {
			(*ctx).Done()

			cl.logger.Infof("missed %s heartbeats, close connection", heartbeatThreshold)
			ticker.Stop()
		}
	}
}
