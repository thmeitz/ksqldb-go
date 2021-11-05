package ksqldb

import (
	"context"
	"fmt"
	"strings"
)

const (
	QBErr             = "qbErr"
	QBUnsupportedType = "unsupported param type"
	EMPTY_STATEMENT   = "empty ksql statement"
)

type QueryBuilder struct {
	stmt      string
	origStmnt string
	options   *QueryBuilderOptions
	ctx       context.Context
}

// QueryBuilderOptions type
type QueryBuilderOptions struct {
	Parse   bool // set to true, to parse the sql by KSQLParser
	Context context.Context
}

// DefaultQueryBuilder returns a @QueryBuilder with default values
func DefaultQueryBuilder(stmnt string) (*QueryBuilder, error) {
	context := context.Background()
	options := QueryBuilderOptions{Parse: false}

	if err := checkEmptyStatement(stmnt); err != nil {
		return nil, err
	}

	return &QueryBuilder{ctx: context, stmt: stmnt, origStmnt: stmnt, options: &options}, nil
}

// QueryBuilderWithOptions returns a @QueryBuilder, which can be configured by @QueryBuilderOptions
func QueryBuilderWithOptions(stmnt string, options *QueryBuilderOptions) (*QueryBuilder, error) {
	ctx := context.Background()

	if err := checkEmptyStatement(stmnt); err != nil {
		return nil, err
	}

	if options != nil && options.Context != nil {
		ctx = options.Context
	}

	return &QueryBuilder{ctx: ctx, stmt: stmnt, origStmnt: stmnt, options: options}, nil
}

// GetInputStatement gets the original statement given do @DefaultQueryBuilder or @QueryBuilderWithOptions
func (qb *QueryBuilder) GetInputStatement() string {
	return qb.stmt
}

// Bind parameters to QueryBuilder
func (qb *QueryBuilder) Bind(params ...interface{}) (*string, error) {
	paramCount := len(params)
	count := strings.Count(qb.origStmnt, "?")
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
		qb.stmt = strings.Replace(qb.stmt, "?", *replace, 1)
	}
	return &qb.stmt, nil
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
		n := fmt.Sprintf("'%v'", param)
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
