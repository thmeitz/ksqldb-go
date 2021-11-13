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
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thmeitz/ksqldb-go"
)

const (
	select1Param  = "select * from bla where column=?"
	select5Params = "insert into bla values(?,?,?,?,?,?)"
)

var qbtests = []struct {
	name    string
	stmnt   string
	value   interface{}
	message string
}{
	{"string", select1Param, "Lara", "select * from bla where column='Lara'"},
	{"nil", select1Param, nil, "select * from bla where column=NULL"},
	{"int", select1Param, 15235, "select * from bla where column=15235"},
	{"hex int", select1Param, 0xff, "select * from bla where column=255"},
	{"int8", select1Param, int8(123), "select * from bla where column=123"},
	{"int16", select1Param, int16(123), "select * from bla where column=123"},
	{"int32", select1Param, int32(123), "select * from bla where column=123"},
	{"int64", select1Param, int64(123), "select * from bla where column=123"},
	{"uint", select1Param, uint(123), "select * from bla where column=123"},
	{"uint8", select1Param, uint8(123), "select * from bla where column=123"},
	{"uint16", select1Param, uint16(123), "select * from bla where column=123"},
	{"uint32", select1Param, uint32(123), "select * from bla where column=123"},
	{"uint64", select1Param, uint64(123), "select * from bla where column=123"},
	{"float32", select1Param, float32(123.99998999998888), "select * from bla where column=123.99999"},
	{"float64", select1Param, float64(123.99998999998888), "select * from bla where column=123.99998999998888"},
	{"bool", select1Param, true, "select * from bla where column=true"},
}

func TestQueryBuilderTypes(t *testing.T) {
	for _, tt := range qbtests {
		t.Run(tt.name, func(t *testing.T) {
			builder, _ := ksqldb.DefaultQueryBuilder(tt.stmnt)
			stmnt, err := builder.Bind(tt.value)
			if err != nil {
				require.Equal(t, tt.message, err.Error())
				return
			}
			require.Equal(t, tt.message, *stmnt)
		})
	}
}

func TestDefaultQueryBuilder_Types5Params(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select5Params)
	stmnt, err := builder.Bind(nil, 1, "2", 3.5, 4, 5)

	require.Nil(t, err)
	require.NotNil(t, stmnt)
	require.Equal(t, "insert into bla values(NULL,1,'2',3.5,4,5)", *stmnt)
}

func TestDefaultQueryBuilder_NotNil(t *testing.T) {
	builder, err := ksqldb.DefaultQueryBuilder(select1Param)
	require.Nil(t, err)
	require.NotNil(t, builder)
}

func TestDefaultQueryBuilder_EmptyStatement(t *testing.T) {
	builder, err := ksqldb.DefaultQueryBuilder("")
	require.Nil(t, builder)
	require.NotNil(t, err)
	require.Equal(t, "qbErr: empty ksql statement", err.Error())
}

func TestDefaultQueryBuilder_GetStatement(t *testing.T) {
	builder, err := ksqldb.DefaultQueryBuilder(select1Param)
	require.NotNil(t, builder)
	require.Nil(t, err)
	stmnt := builder.GetInputStatement()
	require.Equal(t, select1Param, stmnt)
}

func TestQueryBuilderWithOptions_EmptyStatement_NilOptions(t *testing.T) {
	builder, err := ksqldb.QueryBuilderWithOptions("", nil)
	require.Nil(t, builder)
	require.NotNil(t, err)
}

func TestQueryBuilderWithOptions_EmptyOptions(t *testing.T) {
	options := ksqldb.QueryBuilderOptions{}
	builder, err := ksqldb.QueryBuilderWithOptions(select1Param, &options)
	require.NotNil(t, builder)
	require.Nil(t, err)
}

func TestQueryBuilderWithOptions_WithContext(t *testing.T) {
	options := ksqldb.QueryBuilderOptions{Context: context.Background()}
	builder, err := ksqldb.QueryBuilderWithOptions(select1Param, &options)
	require.NotNil(t, builder)
	require.Nil(t, err)
}

func TestQueryBuilder_Bind_ToManyParamsError(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select1Param)
	_, err := builder.Bind(1, "bla", 31235)
	require.NotNil(t, err)
	require.Equal(t, "qbErr: to many params", err.Error())
}

func TestQueryBuilder_Bind_ToFewParamsError(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select1Param)
	_, err := builder.Bind()
	require.NotNil(t, err)
	require.Equal(t, "qbErr: to few params", err.Error())
}

func TestQueryBuilder_Bind_CorrectParams(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select1Param)
	stmnt, err := builder.Bind(1)
	require.Nil(t, err)
	require.NotNil(t, stmnt)
	require.Equal(t, "select * from bla where column=1", *stmnt)
}

func TestQueryBuilder_Bind_MultiParams(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select5Params)
	stmnt, err := builder.Bind(nil, 1, "rainer", 1.98, true, nil)
	require.Nil(t, err)
	require.NotNil(t, stmnt)
	require.Equal(t, "insert into bla values(NULL,1,'rainer',1.98,true,NULL)", *stmnt)
}
