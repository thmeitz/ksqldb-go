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

package parser_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go/parser"
)

func TestParseKSQL_Error(t *testing.T) {
	sql := "SELECT * FROM DOGS"
	err := parser.ParseSql(sql)
	require.NotNil(t, err)
	require.Equal(t, "1 sql syntax error(s) found", err.Error())
	if err != nil {
		for _, e := range *err {
			require.Equal(t, "error on line(1):column(18): missing ';' at '<EOF>'", e.Error())
		}
	}
}

func TestParseKSQL_NoError(t *testing.T) {
	sql := "SELECT * FROM DOGS;"
	err := parser.ParseSql(sql)
	require.Nil(t, err)
}
