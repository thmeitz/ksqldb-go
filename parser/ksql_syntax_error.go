package parser

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type KSqlSyntaxError struct {
	Line, Column int
	Msg          string
}

func (kse *KSqlSyntaxError) Error() string {
	return fmt.Sprintf("error on line(%v):column(%v): %v", kse.Line, kse.Column, kse.Msg)
}

type KSqlSyntaxErrorList []KSqlSyntaxError

func (ksl *KSqlSyntaxErrorList) Error() string {
	return fmt.Sprintf("%v sql syntax errors found", len(*ksl))
}

type KSqlErrorListener struct {
	*antlr.DefaultErrorListener
	Errors KSqlSyntaxErrorList
}

func (c *KSqlErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	c.Errors = append(c.Errors, KSqlSyntaxError{
		Line:   line,
		Column: column,
		Msg:    msg,
	})
}

func (c *KSqlErrorListener) HasErrors() bool {
	return len(c.Errors) > 0
}

func (c *KSqlErrorListener) ErrorCount() int {
	return len(c.Errors)
}
