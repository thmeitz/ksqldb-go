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

package parser

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type SqlSyntaxError struct {
	Line, Column int
	Msg          string
}

func (kse *SqlSyntaxError) Error() string {
	return fmt.Sprintf("error on line(%v):column(%v): %v", kse.Line, kse.Column, kse.Msg)
}

type SqlSyntaxErrorList []SqlSyntaxError

func (ksl *SqlSyntaxErrorList) Error() string {
	return fmt.Sprintf("%v sql syntax error(s) found", len(*ksl))
}

type KSqlErrorListener struct {
	*antlr.DefaultErrorListener
	Errors SqlSyntaxErrorList
}

func (c *KSqlErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	c.Errors = append(c.Errors, SqlSyntaxError{
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
