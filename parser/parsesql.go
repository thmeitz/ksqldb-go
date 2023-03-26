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
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type CanParseSQL interface {
	ParseSql(string) []error
}

func ParseSql(sql string) *SqlSyntaxErrorList {
	errors := SqlSyntaxErrorList{}

	input := antlr.NewInputStream(sql)
	upper := NewUpperCaseStream(input)
	lexerErrorListener := &KSqlErrorListener{}
	lexer := NewKSqlLexer(upper)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrorListener)

	stream := antlr.NewCommonTokenStream(lexer, 0)
	parserErrorListener := &KSqlErrorListener{}
	p := NewKSqlParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(parserErrorListener)

	antlr.ParseTreeWalkerDefault.Walk(&BaseKSqlListener{}, p.Statements())

	if lexerErrorListener.HasErrors() {
		errors = append(errors, lexerErrorListener.Errors...)
	}
	if parserErrorListener.HasErrors() {
		errors = append(errors, parserErrorListener.Errors...)
	}

	if lexerErrorListener.HasErrors() || parserErrorListener.HasErrors() {
		return &errors
	}
	return nil
}
