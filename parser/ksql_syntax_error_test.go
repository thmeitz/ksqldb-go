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

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/stretchr/testify/require"
	mocks "github.com/thmeitz/ksqldb-go/mocks/parser"
	"github.com/thmeitz/ksqldb-go/parser"
)

func TestKSqlSyntaxError_ErrorFunc(t *testing.T) {
	err := parser.SqlSyntaxError{
		Line:   1,
		Column: 10,
		Msg:    "dummy message",
	}
	require.Equal(t, "error on line(1):column(10): dummy message", err.Error())
}

func TestKSqlSyntaxErrorList_ErrorfFunc(t *testing.T) {
	var errorList parser.SqlSyntaxErrorList = []parser.SqlSyntaxError{
		{
			Line:   1,
			Column: 10,
			Msg:    "dummy message",
		},
	}
	require.Equal(t, "1 sql syntax error(s) found", errorList.Error())
}

func TestKSqlErrorListener_HasErrors_ErrorCount(t *testing.T) {
	listener := parser.KSqlErrorListener{}
	require.False(t, listener.HasErrors())
	require.Equal(t, 0, listener.ErrorCount())
}

func TestKSqlErrorListener_SyntaxError(t *testing.T) {
	listener := parser.KSqlErrorListener{}
	require.False(t, listener.HasErrors())
	require.Equal(t, 0, listener.ErrorCount())

	var r antlr.Recognizer
	recognizer := mocks.Recognizer{}
	r = &recognizer
	listener.SyntaxError(r, nil, 1, 10, "some error", nil)

	require.True(t, listener.HasErrors())
	require.Equal(t, 1, listener.ErrorCount())

	listener.SyntaxError(r, nil, 1, 10, "some more error", nil)
	require.Equal(t, 2, listener.ErrorCount())
}
