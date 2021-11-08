/*
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
*/

// https://stackoverflow.com/questions/66067549/how-to-write-a-custom-error-reporter-in-go-target-of-antlr

package main

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/thmeitz/ksqldb-go/parser"
)

func main() {
	//----------| this is the error
	k := `SELECT1 TIMESTAMPTOSTRING(WINDOWSTART,'yyyy-MM-dd HH:mm:ss','Europe/London') AS WINDOW_START, 
				  TIMESTAMPTOSTRING(WINDOWEND,'HH:mm:ss','Europe/London') AS WINDOW_END, 
					DOG_SIZE, DOGS_CT 
	FROM DOGS_BY_SIZE 
	WHERE DOG_SIZE='large';`

	input := antlr.NewInputStream(k)
	lexerErrors := &parser.KSqlErrorListener{}
	lexer := parser.NewKSqlLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrors)

	stream := antlr.NewCommonTokenStream(lexer, 0)
	parserErrors := &parser.KSqlErrorListener{}
	p := parser.NewKSqlParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(parserErrors)

	antlr.ParseTreeWalkerDefault.Walk(&parser.BaseKSqlListener{}, p.Statements())
	fmt.Println(fmt.Sprintf("lexer has errors: %v", lexerErrors.HasErrors()))
	fmt.Println(fmt.Sprintf("parser error count: %v", lexerErrors.ErrorCount()))
	fmt.Println(fmt.Sprintf("parser has errors: %v", parserErrors.HasErrors()))
	fmt.Println(fmt.Sprintf("parser error count: %v", parserErrors.ErrorCount()))
	fmt.Println(parserErrors.Errors)
}
