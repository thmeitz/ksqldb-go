/*
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
*/

package ksqldb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ServerHealth struct {
	IsHealthy *bool `json:"isHealthy"`
	Details   struct {
		Metastore struct {
			IsHealthy *bool `json:"isHealthy"`
		} `json:"metastore"`
		Kafka struct {
			IsHealthy *bool `json:"isHealthy"`
		} `json:"kafka"`
	} `json:"details"`
	KSQLServiceID string `json:"ksqlServiceId"`
}

// ServerInfo provides information about your server
func Healthcheck(api *Client) (*ServerHealth, error) {
	info := ServerHealth{}
	res, err := api.client.Get(api.options.BaseUrl + HEALTHCHECK_ENDPOINT)
	if err != nil {
		api.Close()
		return nil, fmt.Errorf("can't get healthcheck informations: %v", err)
	}
	defer res.Body.Close()

	// close transport layer
	api.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("could not parse the response as JSON:\n%w\n%v", err, string(body))
	}

	return &info, nil
}
