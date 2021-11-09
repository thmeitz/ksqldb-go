package parser_test

import (
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/stretchr/testify/assert"
	"github.com/thmeitz/ksqldb-go/mocks"
	"github.com/thmeitz/ksqldb-go/parser"
)

func TestKSqlSyntaxError_ErrorFunc(t *testing.T) {
	err := parser.KSqlSyntaxError{
		Line:   1,
		Column: 10,
		Msg:    "dummy message",
	}
	assert.Equal(t, "error on line(1):column(10): dummy message", err.Error())
}

func TestKSqlSyntaxErrorList_ErrorfFunc(t *testing.T) {
	var errorList parser.KSqlSyntaxErrorList = []parser.KSqlSyntaxError{
		{
			Line:   1,
			Column: 10,
			Msg:    "dummy message",
		},
	}
	assert.Equal(t, "1 sql syntax error(s) found", errorList.Error())
}

func TestKSqlErrorListener_HasErrors_ErrorCount(t *testing.T) {
	listener := parser.KSqlErrorListener{}
	assert.False(t, listener.HasErrors())
	assert.Equal(t, 0, listener.ErrorCount())
}

func TestKSqlErrorListener_HSyntaxError(t *testing.T) {
	listener := parser.KSqlErrorListener{}
	assert.False(t, listener.HasErrors())
	assert.Equal(t, 0, listener.ErrorCount())

	var r antlr.Recognizer
	recognizer := mocks.Recognizer{}
	r = &recognizer
	listener.SyntaxError(r, nil, 1, 10, "some error", nil)

}
