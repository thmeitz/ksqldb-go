package ksqldb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go"
)

const (
	select1Param  = "select * from bla where column=?"
	select5Params = "insert into bla values(null,?,?,?,?,?)"
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
				assert.Equal(t, tt.message, err.Error())
				return
			}
			assert.Equal(t, tt.message, *stmnt)
		})
	}
}

func TestDefaultQueryBuilder_Types5Params(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select5Params)
	stmnt, err := builder.Bind(nil, 1, "2", 3.5, 4, 5)
	if err != nil {
		fmt.Println(err)
		assert.NotNil(t, err)
		return
	}
	assert.NotNil(t, stmnt)
	assert.Equal(t, "", *stmnt)
}

func TestDefaultQueryBuilder_NotNil(t *testing.T) {
	builder, err := ksqldb.DefaultQueryBuilder(select1Param)
	assert.Nil(t, err)
	assert.NotNil(t, builder)
}

func TestDefaultQueryBuilder_EmptyStatement(t *testing.T) {
	builder, err := ksqldb.DefaultQueryBuilder("")
	assert.Nil(t, builder)
	assert.NotNil(t, err)
	assert.Equal(t, "qbErr: empty ksql statement", err.Error())
}

func TestDefaultQueryBuilder_GetStatement(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select1Param)
	assert.NotNil(t, builder)
	stmnt := builder.GetInputStatement()
	assert.Equal(t, select1Param, stmnt)
}

func TestQueryBuilderWithOptions_EmptyStatement_NilOptions(t *testing.T) {
	builder, err := ksqldb.QueryBuilderWithOptions("", nil)
	assert.Nil(t, builder)
	assert.NotNil(t, err)
}

func TestQueryBuilderWithOptions_EmptyOptions(t *testing.T) {
	options := ksqldb.QueryBuilderOptions{}
	builder, err := ksqldb.QueryBuilderWithOptions(select1Param, &options)
	assert.NotNil(t, builder)
	assert.Nil(t, err)
}

func TestQueryBuilderWithOptions_WithContext(t *testing.T) {
	options := ksqldb.QueryBuilderOptions{Context: context.Background()}
	builder, err := ksqldb.QueryBuilderWithOptions(select1Param, &options)
	assert.NotNil(t, builder)
	assert.Nil(t, err)
}

func TestQueryBuilder_Bind_ToManyParamsError(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select1Param)
	_, err := builder.Bind(1, "bla", 31235)
	assert.NotNil(t, err)
	assert.Equal(t, "qbErr: to many params", err.Error())
}

func TestQueryBuilder_Bind_ToFewParamsError(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select1Param)
	_, err := builder.Bind()
	assert.NotNil(t, err)
	assert.Equal(t, "qbErr: to few params", err.Error())
}

func TestQueryBuilder_Bind_CorrectParams(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select1Param)
	stmnt, err := builder.Bind(1)
	assert.Nil(t, err)
	assert.NotNil(t, stmnt)
	assert.Equal(t, "select * from bla where column=1", *stmnt)
}

func TestQueryBuilder_Bind_MultiParams(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(select5Params)
	stmnt, err := builder.Bind(1, "rainer", 1.98, true, nil)
	assert.Nil(t, err)
	assert.NotNil(t, stmnt)
	assert.Equal(t, "insert into bla values(null,1,'rainer',1.98,true,NULL)", *stmnt)
}
