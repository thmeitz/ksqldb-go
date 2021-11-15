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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/thmeitz/ksqldb-go/net"
)

type RequestParams map[string]interface{}
type Response map[string]interface{}

func newKsqlRequest(api net.HTTPClient, payload io.Reader) (*http.Request, error) {
	return http.NewRequest("POST", api.GetUrl(KSQL_ENDPOINT), payload)
}

func newQueryStreamRequest(api net.HTTPClient, ctx context.Context, payload io.Reader) (*http.Request, error) {
	req, err := newPostRequest(api, ctx, QUERY_STREAM_ENDPOINT, payload)
	return req, err
}

func newCloseQueryRequest(api net.HTTPClient, ctx context.Context, payload io.Reader) (*http.Request, error) {
	return newPostRequest(api, ctx, CLOSE_QUERY_ENDPOINT, payload)
}

func handleRequestError(code int, buf []byte) error {
	ksqlError := ResponseError{}
	fmt.Printf("handleRequestError called: %+v\n", string(buf))
	if err := json.Unmarshal(buf, &ksqlError); err != nil {
		return fmt.Errorf("ksqldb error: %w", err)
	}
	fmt.Printf("ksql: %+v\n", ksqlError)

	return ksqlError
}

func handleGetRequest(httpClient *net.HTTPClient, url string) (*[]byte, error) {
	res, err := (*httpClient).Get(url)
	if err != nil {
		return nil, fmt.Errorf("ksqldb get request failed: %v", err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res.StatusCode, body)
	}
	return &body, nil
}

func newPostRequest(api net.HTTPClient, ctx context.Context, endpoint string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", api.GetUrl(endpoint), payload)
	if err != nil {
		return req, fmt.Errorf("can't create new request with context: %w", err)
	}

	return req, nil
}
