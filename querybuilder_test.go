package ksqldb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go"
)

const (
	SelectBla = "select * from bla where id=?"
)

func TestDefaultQueryBuilder_NotNil(t *testing.T) {
	builder, err := ksqldb.DefaultQueryBuilder(SelectBla)
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
	builder, _ := ksqldb.DefaultQueryBuilder(SelectBla)
	assert.NotNil(t, builder)
	stmnt := builder.GetInputStatement()
	assert.Equal(t, SelectBla, stmnt)
}

func TestQueryBuilderWithOptions_EmptyStatement_NilOptions(t *testing.T) {
	builder, err := ksqldb.QueryBuilderWithOptions("", nil)
	assert.Nil(t, builder)
	assert.NotNil(t, err)
}

func TestQueryBuilderWithOptions_EmptyOptions(t *testing.T) {
	options := ksqldb.QueryBuilderOptions{}
	builder, err := ksqldb.QueryBuilderWithOptions(SelectBla, &options)
	assert.NotNil(t, builder)
	assert.Nil(t, err)
}

func TestQueryBuilderWithOptions_WithContext(t *testing.T) {
	options := ksqldb.QueryBuilderOptions{Context: context.Background()}
	builder, err := ksqldb.QueryBuilderWithOptions(SelectBla, &options)
	assert.NotNil(t, builder)
	assert.Nil(t, err)
}

func TestQueryBuilder_Bind_ToManyParamsError(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(SelectBla)
	_, err := builder.Bind(1, "bla", 31235)
	assert.NotNil(t, err)
	assert.Equal(t, "qbErr: to many params", err.Error())
}

func TestQueryBuilder_Bind_ToFewParamsError(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(SelectBla)
	_, err := builder.Bind()
	assert.NotNil(t, err)
	assert.Equal(t, "qbErr: to few params", err.Error())
}

func TestQueryBuilder_Bind_CorrectParams(t *testing.T) {
	builder, _ := ksqldb.DefaultQueryBuilder(SelectBla)
	stmnt, err := builder.Bind(1)
	assert.Nil(t, err)
	assert.NotNil(t, stmnt)
	assert.Equal(t, "select * from bla where id=1", *stmnt)
}
