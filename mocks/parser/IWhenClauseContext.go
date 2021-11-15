// Code generated by mockery v2.9.4. DO NOT EDIT.

package parser

import (
	antlr "github.com/antlr/antlr4/runtime/Go/antlr"
	mock "github.com/stretchr/testify/mock"

	parser "github.com/thmeitz/ksqldb-go/parser"
)

// IWhenClauseContext is an autogenerated mock type for the IWhenClauseContext type
type IWhenClauseContext struct {
	mock.Mock
}

// Accept provides a mock function with given fields: Visitor
func (_m *IWhenClauseContext) Accept(Visitor antlr.ParseTreeVisitor) interface{} {
	ret := _m.Called(Visitor)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(antlr.ParseTreeVisitor) interface{}); ok {
		r0 = rf(Visitor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// AddChild provides a mock function with given fields: child
func (_m *IWhenClauseContext) AddChild(child antlr.RuleContext) antlr.RuleContext {
	ret := _m.Called(child)

	var r0 antlr.RuleContext
	if rf, ok := ret.Get(0).(func(antlr.RuleContext) antlr.RuleContext); ok {
		r0 = rf(child)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(antlr.RuleContext)
		}
	}

	return r0
}

// AddErrorNode provides a mock function with given fields: badToken
func (_m *IWhenClauseContext) AddErrorNode(badToken antlr.Token) *antlr.ErrorNodeImpl {
	ret := _m.Called(badToken)

	var r0 *antlr.ErrorNodeImpl
	if rf, ok := ret.Get(0).(func(antlr.Token) *antlr.ErrorNodeImpl); ok {
		r0 = rf(badToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*antlr.ErrorNodeImpl)
		}
	}

	return r0
}

// AddTokenNode provides a mock function with given fields: token
func (_m *IWhenClauseContext) AddTokenNode(token antlr.Token) *antlr.TerminalNodeImpl {
	ret := _m.Called(token)

	var r0 *antlr.TerminalNodeImpl
	if rf, ok := ret.Get(0).(func(antlr.Token) *antlr.TerminalNodeImpl); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*antlr.TerminalNodeImpl)
		}
	}

	return r0
}

// EnterRule provides a mock function with given fields: listener
func (_m *IWhenClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	_m.Called(listener)
}

// ExitRule provides a mock function with given fields: listener
func (_m *IWhenClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	_m.Called(listener)
}

// GetAltNumber provides a mock function with given fields:
func (_m *IWhenClauseContext) GetAltNumber() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetBaseRuleContext provides a mock function with given fields:
func (_m *IWhenClauseContext) GetBaseRuleContext() *antlr.BaseRuleContext {
	ret := _m.Called()

	var r0 *antlr.BaseRuleContext
	if rf, ok := ret.Get(0).(func() *antlr.BaseRuleContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*antlr.BaseRuleContext)
		}
	}

	return r0
}

// GetChild provides a mock function with given fields: i
func (_m *IWhenClauseContext) GetChild(i int) antlr.Tree {
	ret := _m.Called(i)

	var r0 antlr.Tree
	if rf, ok := ret.Get(0).(func(int) antlr.Tree); ok {
		r0 = rf(i)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(antlr.Tree)
		}
	}

	return r0
}

// GetChildCount provides a mock function with given fields:
func (_m *IWhenClauseContext) GetChildCount() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetChildren provides a mock function with given fields:
func (_m *IWhenClauseContext) GetChildren() []antlr.Tree {
	ret := _m.Called()

	var r0 []antlr.Tree
	if rf, ok := ret.Get(0).(func() []antlr.Tree); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]antlr.Tree)
		}
	}

	return r0
}

// GetCondition provides a mock function with given fields:
func (_m *IWhenClauseContext) GetCondition() parser.IExpressionContext {
	ret := _m.Called()

	var r0 parser.IExpressionContext
	if rf, ok := ret.Get(0).(func() parser.IExpressionContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(parser.IExpressionContext)
		}
	}

	return r0
}

// GetInvokingState provides a mock function with given fields:
func (_m *IWhenClauseContext) GetInvokingState() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetParent provides a mock function with given fields:
func (_m *IWhenClauseContext) GetParent() antlr.Tree {
	ret := _m.Called()

	var r0 antlr.Tree
	if rf, ok := ret.Get(0).(func() antlr.Tree); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(antlr.Tree)
		}
	}

	return r0
}

// GetParser provides a mock function with given fields:
func (_m *IWhenClauseContext) GetParser() antlr.Parser {
	ret := _m.Called()

	var r0 antlr.Parser
	if rf, ok := ret.Get(0).(func() antlr.Parser); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(antlr.Parser)
		}
	}

	return r0
}

// GetPayload provides a mock function with given fields:
func (_m *IWhenClauseContext) GetPayload() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// GetResult provides a mock function with given fields:
func (_m *IWhenClauseContext) GetResult() parser.IExpressionContext {
	ret := _m.Called()

	var r0 parser.IExpressionContext
	if rf, ok := ret.Get(0).(func() parser.IExpressionContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(parser.IExpressionContext)
		}
	}

	return r0
}

// GetRuleContext provides a mock function with given fields:
func (_m *IWhenClauseContext) GetRuleContext() antlr.RuleContext {
	ret := _m.Called()

	var r0 antlr.RuleContext
	if rf, ok := ret.Get(0).(func() antlr.RuleContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(antlr.RuleContext)
		}
	}

	return r0
}

// GetRuleIndex provides a mock function with given fields:
func (_m *IWhenClauseContext) GetRuleIndex() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetSourceInterval provides a mock function with given fields:
func (_m *IWhenClauseContext) GetSourceInterval() *antlr.Interval {
	ret := _m.Called()

	var r0 *antlr.Interval
	if rf, ok := ret.Get(0).(func() *antlr.Interval); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*antlr.Interval)
		}
	}

	return r0
}

// GetStart provides a mock function with given fields:
func (_m *IWhenClauseContext) GetStart() antlr.Token {
	ret := _m.Called()

	var r0 antlr.Token
	if rf, ok := ret.Get(0).(func() antlr.Token); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(antlr.Token)
		}
	}

	return r0
}

// GetStop provides a mock function with given fields:
func (_m *IWhenClauseContext) GetStop() antlr.Token {
	ret := _m.Called()

	var r0 antlr.Token
	if rf, ok := ret.Get(0).(func() antlr.Token); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(antlr.Token)
		}
	}

	return r0
}

// GetText provides a mock function with given fields:
func (_m *IWhenClauseContext) GetText() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsEmpty provides a mock function with given fields:
func (_m *IWhenClauseContext) IsEmpty() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsWhenClauseContext provides a mock function with given fields:
func (_m *IWhenClauseContext) IsWhenClauseContext() {
	_m.Called()
}

// RemoveLastChild provides a mock function with given fields:
func (_m *IWhenClauseContext) RemoveLastChild() {
	_m.Called()
}

// SetAltNumber provides a mock function with given fields: altNumber
func (_m *IWhenClauseContext) SetAltNumber(altNumber int) {
	_m.Called(altNumber)
}

// SetCondition provides a mock function with given fields: _a0
func (_m *IWhenClauseContext) SetCondition(_a0 parser.IExpressionContext) {
	_m.Called(_a0)
}

// SetException provides a mock function with given fields: _a0
func (_m *IWhenClauseContext) SetException(_a0 antlr.RecognitionException) {
	_m.Called(_a0)
}

// SetInvokingState provides a mock function with given fields: _a0
func (_m *IWhenClauseContext) SetInvokingState(_a0 int) {
	_m.Called(_a0)
}

// SetParent provides a mock function with given fields: _a0
func (_m *IWhenClauseContext) SetParent(_a0 antlr.Tree) {
	_m.Called(_a0)
}

// SetResult provides a mock function with given fields: _a0
func (_m *IWhenClauseContext) SetResult(_a0 parser.IExpressionContext) {
	_m.Called(_a0)
}

// SetStart provides a mock function with given fields: _a0
func (_m *IWhenClauseContext) SetStart(_a0 antlr.Token) {
	_m.Called(_a0)
}

// SetStop provides a mock function with given fields: _a0
func (_m *IWhenClauseContext) SetStop(_a0 antlr.Token) {
	_m.Called(_a0)
}

// String provides a mock function with given fields: _a0, _a1
func (_m *IWhenClauseContext) String(_a0 []string, _a1 antlr.RuleContext) string {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func([]string, antlr.RuleContext) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ToStringTree provides a mock function with given fields: _a0, _a1
func (_m *IWhenClauseContext) ToStringTree(_a0 []string, _a1 antlr.Recognizer) string {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func([]string, antlr.Recognizer) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}