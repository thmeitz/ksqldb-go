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

// GetClusterStatus returns the status of the cluster
// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/cluster-status-endpoint/
func (api *KsqldbClient) GetClusterStatus() (*ClusterStatusResponse, error) {
	var csr ClusterStatusResponse
	var input map[string]interface{}
	var body *[]byte
	var err error

	url := api.http.GetUrl(CLUSTER_STATUS_ENDPOINT)

	if body, err = handleGetRequest(api.http, url); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := api.unMarshalResp(*body, &input); err != nil {
		return nil, fmt.Errorf("could not parse the response:%w", err)
	}

	if err := mapstructure.Decode(&input, &csr); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &csr, nil
}
