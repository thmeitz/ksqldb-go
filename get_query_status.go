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
	"fmt"
)

type QueryStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetQueryStatus returns the current command status for a CREATE, DROP, or TERMINATE statement.
//
// CREATE, DROP, and TERMINATE statements returns an object that indicates the current state of statement execution.
// A statement can be in one of the following states:
//
//    QUEUED, PARSING, EXECUTING: The statement was accepted by the server and is being processed.
//
//    SUCCESS: The statement was successfully processed.
//
//    ERROR: There was an error processing the statement. The statement was not executed.
//
// TERMINATED: The query started by the statement was terminated. Only returned for CREATE STREAM|TABLE AS SELECT.
//
// If a CREATE, DROP, or TERMINATE statement returns a command status with state
// QUEUED, PARSING, or EXECUTING from the @Execute endpoint,
// you can use the @GetQueryStatus endpoint to poll the status of the command.
func (api *KsqldbClient) GetQueryStatus(commandId string) (*QueryStatus, error) {
	var qs QueryStatus
	var body *[]byte
	var err error

	if len(commandId) == 0 {
		return nil, fmt.Errorf("commandId is empty")
	}

	url := api.http.GetUrl(STATUS_ENDPOINT + commandId)

	if body, err = handleGetRequest(api.http, url); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := api.unMarshalResp(*body, &qs); err != nil {
		return nil, fmt.Errorf("could not parse the response:%w", err)
	}

	return &qs, nil
}
