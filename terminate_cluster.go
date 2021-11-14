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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TerminateClusterTopics struct {
	DeleteTopicList []string `json:"deleteTopicList,omitempty"`
}

func (tct *TerminateClusterTopics) Add(topics ...string) {
	tct.DeleteTopicList = append(tct.DeleteTopicList, topics...)
}

func (tct *TerminateClusterTopics) Size() int {
	return len(tct.DeleteTopicList)
}

// TerminateCluster terminates your cluster
//
// This is a `Terminate` requests and the response is a terminate response
// see: https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/ksql-endpoint/#common-fields
//
// server logs:
// INFO Received: ClusterTerminateRequest{deleteTopicList=[DOGS_BY_SIZE, dogs]} (io.confluent.ksql.rest.server.resources.KsqlResource:216)
// INFO Terminating the KSQL server. (io.confluent.ksql.rest.server.computation.CommandRunner:374)
// INFO 172.18.0.1 - - "POST /ksql/terminate HTTP/2.0" 200 242 "-" "Go-http-client/2.0" 43 (io.confluent.ksql.api.server.LoggingHandler:113)
// INFO The KSQL server was terminated. (io.confluent.ksql.rest.server.computation.CommandRunner:380)
// INFO Closing command store (io.confluent.ksql.rest.server.computation.CommandRunner:479)

func (api *KsqldbClient) TerminateCluster(topics ...string) (*KsqlResponseSlice, error) {
	result := new(KsqlResponseSlice)
	tpc := TerminateClusterTopics{}
	var b []byte
	var err error

	url := (*api.http).GetUrl(TERMINATE_CLUSTER_ENDPOINT)
	if len(topics) > 0 {
		tpc.Add(topics...)
	}

	if b, err = json.Marshal(&tpc); err != nil {
		return nil, fmt.Errorf("can't marshal data %w", err)
	}

	res, err := (*api.http).Post(url, "application/vnd.ksql.v1+json", bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if res.StatusCode != http.StatusOK {
		return nil, handleRequestError(res.StatusCode, body)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("could not parse the response:%w", err)
	}

	return result, nil
}
