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

package ksqldb_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go"
)

var fullBlown = `{
  "clusterStatus": {
    "localhost:8088": {
      "hostAlive": true,
      "lastStatusUpdateMs": 1617609098808,
      "activeStandbyPerQuery": {
        "CTAS_MY_AGG_TABLE_3": {
          "activeStores": [
            "Aggregate-Aggregate-Materialize"
          ],
          "activePartitions": [
            {
              "topic": "my_stream",
              "partition": 1
            },
            {
              "topic": "my_stream",
              "partition": 3
            },
            {
              "topic": "_confluent-ksql-default_query_CTAS_MY_AGG_TABLE_3-Aggregate-GroupBy-repartition",
              "partition": 1
            },
            {
              "topic": "_confluent-ksql-default_query_CTAS_MY_AGG_TABLE_3-Aggregate-GroupBy-repartition",
              "partition": 3
            }
          ],
          "standByStores": [],
          "standByPartitions": []
        }
      },
      "hostStoreLags": {
        "stateStoreLags": {
          "_confluent-ksql-default_query_CTAS_MY_AGG_TABLE_3#Aggregate-Aggregate-Materialize": {
            "lagByPartition": {
              "1": {
                "currentOffsetPosition": 0,
                "endOffsetPosition": 0,
                "offsetLag": 0
              },
              "3": {
                "currentOffsetPosition": 0,
                "endOffsetPosition": 0,
                "offsetLag": 0
              }
            },
            "size": 2
          }
        },
        "updateTimeMs": 1617609168917
      }
    },
    "other.ksqldb.host:8088": {
      "hostAlive": true,
      "lastStatusUpdateMs": 1617609172614,
      "activeStandbyPerQuery": {
        "CTAS_MY_AGG_TABLE_3": {
          "activeStores": [
            "Aggregate-Aggregate-Materialize"
          ],
          "activePartitions": [
            {
              "topic": "my_stream",
              "partition": 0
            },
            {
              "topic": "my_stream",
              "partition": 2
            },
            {
              "topic": "_confluent-ksql-default_query_CTAS_MY_AGG_TABLE_3-Aggregate-GroupBy-repartition",
              "partition": 0
            },
            {
              "topic": "_confluent-ksql-default_query_CTAS_MY_AGG_TABLE_3-Aggregate-GroupBy-repartition",
              "partition": 2
            }
          ],
          "standByStores": [],
          "standByPartitions": []
        }
      },
      "hostStoreLags": {
        "stateStoreLags": {
          "_confluent-ksql-default_query_CTAS_MY_AGG_TABLE_3#Aggregate-Aggregate-Materialize": {
            "lagByPartition": {
              "0": {
                "currentOffsetPosition": 1,
                "endOffsetPosition": 1,
                "offsetLag": 0
              },
              "2": {
                "currentOffsetPosition": 0,
                "endOffsetPosition": 0,
                "offsetLag": 0
              }
            },
            "size": 2
          }
        },
        "updateTimeMs": 1617609170111
      }
    }
  }
}
`

var hostStatusJson = `
{
	"hostAlive": true,
	"lastStatusUpdateMs": 1617609098808,
	"activeStandbyPerQuery": {},
	"hostStoreLags": {
		"stateStoreLags": {
			"_lag1": {
				"lagByPartition": {
					"1": {
						"currentOffsetPosition": 0,
						"endOffsetPosition": 0,
						"offsetLag": 0
					},
					"3": {
						"currentOffsetPosition": 0,
						"endOffsetPosition": 0,
						"offsetLag": 0
					}
				},
				"size": 2
			},
			"_lag2": {
				"lagByPartition": {
					"1": {
						"currentOffsetPosition": 0,
						"endOffsetPosition": 0,
						"offsetLag": 0
					}            
				},
				"size": 1
			}
		},				
		"updateTimeMs": 1617609168917
	}
}
`

var lagByPartition = `
{
	"lagByPartition": {
		"1": {
			"currentOffsetPosition": 11,
			"endOffsetPosition": 12,
			"offsetLag":13
		},
		"3": {
			"currentOffsetPosition": 31,
			"endOffsetPosition": 32,
			"offsetLag": 33
		}
	}
}
`

var part = ` 
	{
		"currentOffsetPosition": 1,
		"endOffsetPosition": 2,
		"offsetLag": 3
	}
`

func XTestClusterStatusResponse(t *testing.T) {
	var csr ksqldb.ClusterStatusResponse
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(fullBlown), &input); err != nil {
		fmt.Printf("could not parse the response as JSON:%+v", err)
	}
	//fmt.Println(input)
	err := mapstructure.Decode(input, &csr)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v", csr)
}

func TestDecodeClusterNode_HostStatusJson(t *testing.T) {
	var host ksqldb.HostStatus
	if err := json.Unmarshal([]byte(hostStatusJson), &host); err != nil {
		fmt.Printf("could not parse the response as JSON:%+v", err)
	}
	fmt.Printf("%+v", host)
	assert.True(t, host.HostAlive)
	assert.Equal(t, uint64(1617609098808), host.LastStatusUpdateMs)
	assert.Equal(t, uint64(1617609168917), host.HostStoreLags.UpdateTimeMs)
	assert.Equal(t, uint64(2), host.HostStoreLags.StateStoreLags.LagByPartition["_lag1"].Size)
}

func TestDecodeClusterNode_LagByPartitionJSON(t *testing.T) {
	var lag ksqldb.LagByPartition
	if err := json.Unmarshal([]byte(lagByPartition), &lag); err != nil {
		fmt.Printf("could not parse the response as JSON:%+v", err)
	}
	assert.Equal(t, uint64(11), lag.Partition["1"].CurrentOffsetPosition)
	assert.Equal(t, uint64(12), lag.Partition["1"].EndOffsetPosition)
	assert.Equal(t, uint64(13), lag.Partition["1"].OffsetLag)
	assert.Equal(t, uint64(31), lag.Partition["3"].CurrentOffsetPosition)
	assert.Equal(t, uint64(32), lag.Partition["3"].EndOffsetPosition)
	assert.Equal(t, uint64(33), lag.Partition["3"].OffsetLag)
}

func TestDecodeClusterNode_PartitionJSON(t *testing.T) {
	var partition ksqldb.Partition
	if err := json.Unmarshal([]byte(part), &partition); err != nil {
		fmt.Printf("could not parse the response as JSON:%+v", err)
	}
	assert.Equal(t, uint64(1), partition.CurrentOffsetPosition)
	assert.Equal(t, uint64(2), partition.EndOffsetPosition)
	assert.Equal(t, uint64(3), partition.OffsetLag)
}
