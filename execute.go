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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Execute will execute a ksqlDB statement, such as creating
// a new stream or table. To run queries use Push or
// Pull functions.
//
// To use this function pass in the base URL of your
// ksqlDB server, and the SQL query statement
//
// Ref: https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/ksql-endpoint/
//
// TODO Add support for commandSequenceNumber and streamsProperties
// TODO Add better support for responses to CREATE/DROP/TERMINATE (e.g. commandID, commandStatus.status, etc).
func (cl *Client) Execute(q string) (err error) {
	// Create the client
	// TODO: this should be refactored, since we can't mockup the cient
	// should this client in our Client?
	client := &http.Client{}

	//  make the request
	payload := strings.NewReader(`{"ksql":"` + cl.SanitizeQuery(q) + `"}`)

	req, err := cl.newKsqlRequest(payload)
	cl.logger.Debug("sending ksqlDB request:%v", q)
	if err != nil {
		return fmt.Errorf("can't create new request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("can't do request: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("can't read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("the http request did not return a success code: %v / %v", res.StatusCode, string(body))
	}

	return nil
}
