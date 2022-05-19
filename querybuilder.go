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

package ksqldb

import (
	"fmt"
	"strings"
)

const (
	QBErr             = "qbErr"
	QBUnsupportedType = "unsupported param type"
	EMPTY_STATEMENT   = "empty ksql statement"
)

// QueryBuilder replaces ? with the correct types in the sql statement
func QueryBuilder(stmnt string, params ...interface{}) (*string, error) {
	var result *string
	var err error

	if err = checkEmptyStatement(stmnt); err != nil {
		return nil, err
	}

	if result, err = bind(stmnt, params...); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result, nil
}

// ParseQueryBuilder parses the
// func ParseQueryBuilder(stmnt string, parse bool, params ...interface{}) (*string, error) {
// 	var result *string
// 	var err error
//
// 	result, err = QueryBuilder(stmnt, params...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if !parse {
// 		return result, err
// 	}
//
// 	if ksqlErr := parser.ParseSql(*result); ksqlErr != nil {
// 		return nil, ksqlErr
// 	}
//
// 	return result, nil
// }

// bind parameters to QueryBuilder
func bind(stmnt string, params ...interface{}) (*string, error) {
	paramCount := len(params)
	count := strings.Count(stmnt, "?")
	if paramCount < count {
		return nil, fmt.Errorf("%v: %v", QBErr, "to few params")
	} else if paramCount > count {
		return nil, fmt.Errorf("%v: %v", QBErr, "to many params")
	}

	for _, param := range params {
		replace, err := getReplacement(param)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		stmnt = strings.Replace(stmnt, "?", *replace, 1)
	}
	return &stmnt, nil
}

func getReplacement(param interface{}) (*string, error) {
	switch param := param.(type) {
	case int:
		n := fmt.Sprintf("%v", int(param))
		return &n, nil
	case int8:
		n := fmt.Sprintf("%v", int8(param))
		return &n, nil
	case int16:
		n := fmt.Sprintf("%v", int16(param))
		return &n, nil
	case int32:
		n := fmt.Sprintf("%v", int32(param))
		return &n, nil
	case int64:
		n := fmt.Sprintf("%v", int64(param))
		return &n, nil
	case uint:
		n := fmt.Sprintf("%v", uint(param))
		return &n, nil
	case uint8:
		n := fmt.Sprintf("%v", uint8(param))
		return &n, nil
	case uint16:
		n := fmt.Sprintf("%v", uint16(param))
		return &n, nil
	case uint32:
		n := fmt.Sprintf("%v", uint32(param))
		return &n, nil
	case uint64:
		n := fmt.Sprintf("%v", uint64(param))
		return &n, nil
	case float32:
		n := fmt.Sprintf("%v", float32(param))
		return &n, nil
	case float64:
		n := fmt.Sprintf("%v", float64(param))
		return &n, nil
	case nil:
		n := "NULL"
		return &n, nil
	case string:
		n := fmt.Sprintf("'%v'", strings.Replace(param, "'", "''", -1))
		return &n, nil
	case bool:
		n := fmt.Sprintf("%v", param)
		return &n, nil
	default:
		return nil, fmt.Errorf("%v: %v :%v", QBErr, QBUnsupportedType, param)
	}
}

// checkEmptyStatement to keep the Factories dry
func checkEmptyStatement(stmnt string) error {
	if stmnt == "" {
		return fmt.Errorf("qbErr: %v", EMPTY_STATEMENT)
	}
	return nil
}
