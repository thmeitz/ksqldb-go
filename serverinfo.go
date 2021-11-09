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
	"net/http"

	"github.com/thmeitz/ksqldb-go/net"
)

// ServerInfo gets the info for your server
// url = api.options.BaseUrl +
func GetServerInfo(api net.KSqlDBClient) (*ServerInfo, error) {
	info := ServerInfoResponse{}
	res, err := http.Get(api.GetUrl(INFO_ENDPOINT))

	if err != nil {
		// TODO: we have to close the transport api.Close()
		return nil, fmt.Errorf("can't get server info: %v", err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("could not parse the response as JSON: %w", err)
	}

	return &info.KSQLServerInfo, nil
}
