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
)

// TODO: check type for LastStatusUpdateMs
type ClusterNode struct {
	HostAlive          bool  `json:"hostAlive"`
	LastStatusUpdateMs int64 `json:"lastStatusUpdateMs"`
}

type ClusterNodeMap map[string]ClusterNode

type ClusterStatus struct {
	Hostname ClusterNodeMap
}
type ClusterStatusResponse struct {
	ClusterStatus ClusterStatus `json:"clusterStatus"`
}

type PartitionMap map[string]Partition

type Partition struct {
	CurrentOffsetPosition uint64 `json:"currentOffsetPosition"`
	EndOffsetPosition     uint64 `json:"endOffsetPosition"`
	OffsetLag             uint64 `json:"offsetLag"`
}

type LagByPartitionMap map[string]LagByPartition

type lagByPartition LagByPartition

func (pm *LagByPartition) UnmarshalJSON(b []byte) (err error) {
	var orig lagByPartition
	if err := json.Unmarshal(b, &orig); err != nil {
		return fmt.Errorf("could not parse the response as JSON:%w", err)
	}
	if pm.Partition == nil {
		pm.Partition = make(PartitionMap)
	}
	for k, v := range orig.Partition {
		pm.Partition[k] = v
	}
	return
}

type LagByPartition struct {
	Partition PartitionMap `json:"lagByPartition"`
	Size      uint64       `json:"size"`
}

type StateStoreLags struct {
	LagByPartition LagByPartitionMap `json:"lagByPartition"`
}

func (csr *StateStoreLags) UnmarshalJSON(b []byte) (err error) {
	var orig LagByPartitionMap

	if err := json.Unmarshal(b, &orig); err != nil {
		return fmt.Errorf("could not parse the response as JSON:%w", err)
	}

	fmt.Printf("%+v\n=========\n%v", orig, string(b))
	// fmt.Printf("%+v\n=========\n", orig)
	if csr.LagByPartition == nil {
		csr.LagByPartition = make(LagByPartitionMap)
	}

	for k, v := range orig {
		csr.LagByPartition[k] = v
		//csr.LagByPartition[k].Size = uint64(1)
	}
	return
}

type HostStoreLags struct {
	StateStoreLags StateStoreLags `json:"stateStoreLags"`
	UpdateTimeMs   uint64         `json:"updateTimeMs"`
}

type HostStatus struct {
	HostAlive             bool                   `json:"hostAlive"`
	LastStatusUpdateMs    uint64                 `json:"lastStatusUpdateMs"`
	HostStoreLags         HostStoreLags          `json:"hostStoreLags"`
	ActiveStandbyPerQuery map[string]interface{} `json:"activeStandbyPerQuery"`
}

// GetClusterStatus
// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/cluster-status-endpoint/
func (c *KsqldbClient) GetClusterStatus() (*ClusterStatusResponse, error) {
	var csr ClusterStatusResponse

	url := (*c.http).GetUrl(CLUSTER_STATUS_ENDPOINT)

	res, err := (*c.http).Get(url)
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

	if err := json.Unmarshal(body, &csr); err != nil {
		return nil, fmt.Errorf("could not parse the response as JSON:%w", err)
	}

	fmt.Println(csr)

	cs := ClusterStatusResponse{}

	return &cs, nil
}
