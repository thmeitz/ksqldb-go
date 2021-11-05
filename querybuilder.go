package ksqldb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const (
	QBErr           = "qbErr"
	EMPTY_STATEMENT = "empty ksql statement"
)

type QueryBuilder struct {
	stmt      string
	origStmnt string
	lastErr   error
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

	for idx, param := range params {
		fmt.Println(idx, ":", param)
		qb.stmt = strings.Replace(qb.stmt, "?", strconv.FormatInt(int64(param.(int)), 10), 1)
	}
	return &qb.stmt, nil
}

// checkEmptyStatement to keep the Factories dry
func checkEmptyStatement(stmnt string) error {
	if stmnt == "" {
		return fmt.Errorf("qbErr: %v", EMPTY_STATEMENT)
	}
	return nil
}
