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
)

// KsqlServerInfo
type KsqlServerInfo struct {
	Version        string `json:"version"`
	KafkaClusterID string `json:"kafkaClusterId"`
	KsqlServiceID  string `json:"ksqlServiceId"`
	ServerStatus   string `json:"serverStatus,omitempty"`
}

// KsqlServerInfoResponse
type KsqlServerInfoResponse struct {
	KsqlServerInfo KsqlServerInfo `json:"KsqlServerInfo"`
}

// ServerInfo gets the info for your server
// api net.KsqlHTTPClient
func (api *KsqldbClient) GetServerInfo(ctx context.Context) (info *KsqlServerInfo, err error) {
	response := KsqlServerInfoResponse{}
	res, err := api.http.Get(ctx, api.http.GetUrl(INFO_ENDPOINT))
	if err != nil {
		return nil, fmt.Errorf("can't get server info: %v", err)
	}
	defer func() {
		berr := res.Body.Close()
		if err == nil {
			err = berr
		}
	}()

	body, readErr := api.readBody(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if err := api.unMarshalResp(body, &response); err != nil {
		return nil, fmt.Errorf("could not parse the response: %w", err)
	}

	info = &response.KsqlServerInfo

	return
}
