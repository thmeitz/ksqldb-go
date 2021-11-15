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
	"encoding/json"
	"fmt"
)

// ValidateProperty resource tells you whether a property is prohibited from setting.
// If prohibited the ksqlDB server api returns a 400 error
func (api *KsqldbClient) ValidateProperty(property string) (*bool, error) {
	var input bool
	var body *[]byte
	var err error

	if len(property) < 1 {
		return nil, fmt.Errorf("property must not empty")
	}

	url := api.http.GetUrl(PROP_VALIDITY_ENPOINT + "/" + property)

	if body, err = handleGetRequest(api.http, url); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := json.Unmarshal(*body, &input); err != nil {
		return nil, fmt.Errorf("could not parse the response:%w", err)
	}

	return &input, nil
}
