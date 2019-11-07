package gorules

import (
	"errors"
	"go/ast"
	"go/parser"
	"reflect"
)

// Rule ...
type Rule interface {
	Bool(interface{}) (bool, error)
	Int(interface{}) (int64, error)
	Float(interface{}) (float64, error)
}
type rule struct {
	expr ast.Expr
}

// NewRule 提前解析规则,不用每次都重新解析
func NewRule(r string) (Rule, error) {
	if len(r) == 0 {
		return nil, ErrRuleEmpty
	}
	expr, err := parser.ParseExpr(r)
	if err != nil {
		return nil, err
	}
	return &rule{expr}, nil
}

func (r *rule) Bool(x interface{}) (bool, error) {
	typ := reflect.ValueOf(x)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	b, err := getValue(typ, r.expr)
	if err != nil {
		return false, err
	}
	if r, ok := b.(bool); ok {
		return r, nil
	}
	return false, errors.New("result not bool")
}

func (r *rule) Int(x interface{}) (int64, error) {
	typ := reflect.ValueOf(x)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	b, err := getValue(typ, r.expr)
	if err != nil {
		return 0, err
	}
	if r, ok := b.(float64); ok {
		return int64(r), nil
	} else if i, ok := b.(int64); ok {
		return i, nil
	}
	return 0, errors.New("result not int")
}

func (r *rule) Float(x interface{}) (float64, error) {
	typ := reflect.ValueOf(x)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	b, err := getValue(typ, r.expr)
	if err != nil {
		return 0, err
	}
	if r, ok := b.(float64); ok {
		return r, nil
	}
	if r, ok := b.(int64); ok {
		return float64(r), nil
	}
	return 0, errors.New("result not float")
}
