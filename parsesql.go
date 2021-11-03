/*
Copyright © 2021 Robin Moffat & Contributors
Copyright © 2021 Thomas Meitz <thme219@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Parts of this apiclient are borrowed from Zalando Skipper
https://github.com/zalando/skipper/blob/master/net/httpclient.go

Zalando licence: MIT
https://github.com/zalando/skipper/blob/master/LICENSE
*/

package ksqldb

import (
	"github.com/Masterminds/log-go"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/thmeitz/ksqldb-go/parser"
)

func (cl *Client) ParseKSQL(sql string) *parser.KSqlSyntaxErrorList {
	errors := parser.KSqlSyntaxErrorList{}

	input := antlr.NewInputStream(sql)
	lexerErrorListener := &parser.KSqlErrorListener{}
	lexer := parser.NewKSqlLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrorListener)

	stream := antlr.NewCommonTokenStream(lexer, 0)
	parserErrorListener := &parser.KSqlErrorListener{}
	p := parser.NewKSqlParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(parserErrorListener)

	antlr.ParseTreeWalkerDefault.Walk(&parser.BaseKSqlListener{}, p.Statements())

	if lexerErrorListener.HasErrors() {
		cl.logger.Errorw("lexer error(s)", log.Fields{"count": lexerErrorListener.ErrorCount(), "errors": lexerErrorListener.Errors})
		errors = append(errors, lexerErrorListener.Errors...)
	}
	if parserErrorListener.HasErrors() {
		cl.logger.Errorw("parser error(s)", log.Fields{"count": parserErrorListener.ErrorCount(), "errors": parserErrorListener.Errors})
		errors = append(errors, parserErrorListener.Errors...)
	}

	if lexerErrorListener.HasErrors() || parserErrorListener.HasErrors() {
		return &errors
	}
	return nil
}
