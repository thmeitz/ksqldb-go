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

	"github.com/mitchellh/mapstructure"
)

type ClusterStatusResponse struct {
	ClusterStatus ClusterStatus
}

type ClusterStatus struct {
	Host ClusterNodeMap `mapstructure:",remain"`
}

type ClusterNodeMap map[string]ClusterNode

type ClusterNode struct {
	HostAlive             bool
	LastStatusUpdateMs    int64
	HostStoreLags         HostStoreLags
	ActiveStandbyPerQuery ActiveStandbyPerQueryMap
}

type TopicPartition struct {
	Topic     string
	Partition uint64
}

type ActiveStandbyPerQueryMap map[string]ActiveStandbyPerQuery

type ActiveStandbyPerQuery struct {
	ActiveStores      []string
	ActivePartitions  []TopicPartition
	StandByStore      []string
	StandByPartitions []string
}

type HostStoreLags struct {
	StateStoreLags StateStoreLagMap
	UpdateTimeMs   uint64
}

type StateStoreLagMap map[string]StateStoreLag

type StateStoreLag struct {
	LagByPartition LagByPartitionMap
	Size           uint64
}

type LagByPartitionMap map[string]LagByPartition

type LagByPartition struct {
	Partition Partition
}

type PartitionMap map[string]Partition

type Partition struct {
	CurrentOffsetPosition uint64
	EndOffsetPosition     uint64
	OffsetLag             uint64
}

// GetClusterStatus
// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/cluster-status-endpoint/
func (api *KsqldbClient) GetClusterStatus() (*ClusterStatusResponse, error) {
	var csr ClusterStatusResponse
	var input map[string]interface{}

	url := (*api.http).GetUrl(CLUSTER_STATUS_ENDPOINT)

	res, err := (*api.http).Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't get cluster status: %v", err)
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
		return nil, fmt.Errorf("could not parse the response as JSON:%w", err)
	}

	if err := mapstructure.Decode(&input, &csr); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &csr, nil
}
