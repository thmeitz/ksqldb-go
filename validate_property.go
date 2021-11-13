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
	"io/ioutil"
	"net/http"
)

// ValidateProperty resource tells you whether a property is prohibited from setting.
// If prohibited the ksqlDB server api returns a 400 error
func (api *KsqldbClient) ValidateProperty(property string) (*bool, error) {
	var input bool

	if len(property) < 1 {
		return nil, fmt.Errorf("property must not empty")
	}

	url := (*api.http).GetUrl(PROP_VALIDITY_ENPOINT + "/" + property)

	res, err := (*api.http).Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't get validity of property: %v", err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res.StatusCode, body)
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return nil, fmt.Errorf("could not parse the response:%w", err)
	}

	return &input, nil
}
