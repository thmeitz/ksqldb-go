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

package ksqldb_test

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
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

func TestClusterStatusResponse(t *testing.T) {
	var csr ksqldb.ClusterStatusResponse
	var err error
	var input map[string]interface{}

	err = json.Unmarshal([]byte(fullBlown), &input)
	require.Nil(t, err)

	err = mapstructure.Decode(&input, &csr)
	require.Nil(t, err)
}
