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
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ServerInfo struct {
	Version        string `json:"version"`
	KafkaClusterID string `json:"kafkaClusterId"`
	KSQLServiceID  string `json:"ksqlServiceId"`
}

type ServerInfoResponse struct {
	KSQLServerInfo ServerInfo `json:"KsqlServerInfo"`
}

// ServerInfo gets the info for your server
func GetServerInfo(api *Client) (*ServerInfo, error) {
	info := ServerInfoResponse{}
	res, err := api.client.Get(api.options.BaseUrl + INFO_ENDPOINT)

	if err != nil {
		api.Close()
		return nil, fmt.Errorf("can't get server info: %v", err)
	}
	defer res.Body.Close()

	api.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("could not parse the response as JSON:\n%w\n%v", err, string(body))
	}

	return &info.KSQLServerInfo, nil
}
