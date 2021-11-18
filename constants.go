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

const (
	QUERY_STREAM_ENDPOINT      = "/query-stream"
	QUERY_ENDPOINT             = "/query"
	INSERTS_ENDPOINT           = "/inserts-stream"
	CLOSE_QUERY_ENDPOINT       = "/close-query"
	KSQL_ENDPOINT              = "/ksql"
	INFO_ENDPOINT              = "/info"
	STATUS_ENDPOINT            = "/status"
	HEALTHCHECK_ENDPOINT       = "/healthcheck"
	CLUSTER_STATUS_ENDPOINT    = "/clusterStatus"
	PROP_VALIDITY_ENPOINT      = "/is_valid_property"
	TERMINATE_CLUSTER_ENDPOINT = "/ksql/terminate"
)
