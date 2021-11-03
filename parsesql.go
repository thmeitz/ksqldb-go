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
