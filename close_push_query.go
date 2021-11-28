/*
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
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Close Query terminates push query explicitly
func (api *KsqldbClient) ClosePushQuery(ctx context.Context, queryID string) error {
	payload := strings.NewReader(`{"queryId":"` + queryID + `"}`)
	req, err := newPostRequest(api.http, ctx, CLOSE_QUERY_ENDPOINT, payload)

	if err != nil {
		return fmt.Errorf("failed to create post request to terminate query: %w", err)
	}

	res, err := api.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute post request to terminate query: %w", err)
	}
	defer res.Body.Close()

	body, err := api.readBody(res.Body)
	if err != nil {
		return fmt.Errorf("can't read response body: %w", err)
	}
	// handleError
	if res.StatusCode != http.StatusOK {
		return handleRequestError(res.StatusCode, body)
	}
	return nil
}
