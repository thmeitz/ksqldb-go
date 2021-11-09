package ksqldb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go"
)

func TestParseKSQL_Error(t *testing.T) {
	sql := "select * from dogs"
	err := ksqldb.ParseKSQL(sql)
	assert.NotNil(t, err)
	assert.Equal(t, "1 sql syntax error(s) found", err.Error())
}

func TestParseKSQL_NoError(t *testing.T) {
	sql := "SELECT * FROM DOGS;"
	err := ksqldb.ParseKSQL(sql)
	assert.Nil(t, err)
}
