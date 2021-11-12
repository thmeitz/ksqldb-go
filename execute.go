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

	"github.com/thmeitz/ksqldb-go/internal"
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
func (api *KsqldbClient) Execute(sql string) (err error) {

	// first sanitize the query
	query := internal.SanitizeQuery(sql)
	// we're kick in our ksqlparser to check the query string
	ksqlerr := ParseSql(query)
	if ksqlerr != nil {
		return ksqlerr
	}
	//  make the request
	payload := strings.NewReader(`{"ksql":"` + query + `"}`)

	req, err := newKsqlRequest(*api.http, payload)
	// api.logger.Debugf("sending ksqlDB request:%v", q)
	if err != nil {
		return fmt.Errorf("can't create new request: %w", err)
	}

	res, err := (*api.http).Do(req)
	if err != nil {
		return fmt.Errorf("can't do request: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("can't read response body: %w", err)
	}

	// this is only one side of the coin
	if res.StatusCode != http.StatusOK {
		return handleRequestError(res.StatusCode, body)
	}

	return nil
}
