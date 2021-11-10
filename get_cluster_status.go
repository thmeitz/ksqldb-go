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
)

// GetClusterStatus
// @see https://docs.ksqldb.io/en/latest/developer-guide/ksqldb-rest-api/cluster-status-endpoint/
func (c *KsqldbClient) GetClusterStatus() (*ClusterStatusResponse, error) {
	csr := ClStatResp{}
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

	if err := json.Unmarshal(body, &csr); err != nil {
		return nil, fmt.Errorf("could not parse the response as JSON:%w", err)
	}
	cs := ClusterStatusResponse{}

	return &cs, nil
}

type ClusterNode struct {
	Name               string
	HostAlive          bool  `json:"hostAlive"`
	LastStatusUpdateMs int64 `json:"lastStatusUpdateMs"`
}

type ClusterStatus struct {
	Node []ClusterNode
}

type ClusterStatusResponse struct {
	ClusterStatus ClusterStatus `json:"clusterStatus"`
}

type ClStatResp map[string]interface{}

// the return types from the ksqldb are creapy, so we translate them
// to better readable types
type hostStatus struct {
	HostAlive             bool                   `json:"hostAlive"`
	LastStatusUpdateMs    int64                  `json:"lastStatusUpdateMs"`
	ActiveStandbyPerQuery map[string]interface{} `json:"activeStandbyPerQuery"`
	HostStoreLags         map[string]interface{} `json:"hostStoreLags"`
}

func (clr *ClStatResp) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	nodes := []ClusterNode{}

	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(err)
		return err
	}
	if x, found := v["clusterStatus"]; found {
		switch x := x.(type) {
		case map[string]interface{}:
			{

				for key, value := range x {
					fmt.Println(key, value)
					node := ClusterNode{
						Name: key,
					}

					nodes = append(nodes, node)
				}

				break
			}

		}
		fmt.Println(nodes)
		//for key, value := range x {
		//
		//}
	}

	return nil
}
