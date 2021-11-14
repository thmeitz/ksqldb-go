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
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
)

func TestTerminateClusterTopics_Add(t *testing.T) {
	topicList := ksqldb.TerminateClusterTopics{}
	topicList.Add("FOO", "bar.*")
	require.Equal(t, 2, topicList.Size())
	require.Equal(t, "FOO", topicList.DeleteTopicList[0])
	require.Equal(t, "bar.*", topicList.DeleteTopicList[1])
}
