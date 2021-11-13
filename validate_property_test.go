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
	"github.com/thmeitz/ksqldb-go/net"
)

func TestValidateProperty_EmptyProperty(t *testing.T) {
	var options = net.Options{
		BaseUrl:   "http://localhost:8088",
		AllowHTTP: true,
	}
	kcl, err := ksqldb.NewClientWithOptions(options)
	require.Nil(t, err)
	require.NotNil(t, kcl)
	val, err := kcl.ValidateProperty("")
	require.Equal(t, "property must not empty", err.Error())
	require.Nil(t, val)
}
