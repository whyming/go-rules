package gorules

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

// 错误定义
var (
	ErrTypeNotStruct  = errors.New("value must struct or struct pointer")
	ErrNotFoundTag    = errors.New("not found tag")
	ErrUnsupportToken = errors.New("unsupport token")
	ErrUnsupportExpr  = errors.New("unsupport expr")
	ErrNotNumber      = errors.New("not a number")
	ErrNotBool        = errors.New("not boolean")
)

// Bool 规则rule结果的布尔值，rule的参数基于base的json tag
func Bool(base interface{}, rule string) (bool, error) {
	if len(rule) == 0 {
		return true, nil
	}
	expr, err := parser.ParseExpr(rule)
	if err != nil {
		return false, err
	}
	typ := reflect.ValueOf(base)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	b, err := getValue(typ, expr)
	if err != nil {
		return false, err
	}
	if r, ok := b.(bool); ok {
		return r, nil
	}
	return false, errors.New("result not bool")
}

// Int 规则rule的结果如果是数值型，转换为int64，否则报错
func Int(base interface{}, rule string) (int64, error) {
	if len(rule) == 0 {
		return 0, nil
	}
	expr, err := parser.ParseExpr(rule)
	if err != nil {
		return 0, err
	}
	typ := reflect.ValueOf(base)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	b, err := getValue(typ, expr)
	if err != nil {
		return 0, err
	}
	if r, ok := b.(float64); ok {
		return int64(r), nil
	}
	return 0, errors.New("result not int")
}

// Float 返回规则rule结果，如果数值型返回float64，否则报错
func Float(base interface{}, rule string) (float64, error) {
	if len(rule) == 0 {
		return 0, nil
	}
	expr, err := parser.ParseExpr(rule)
	if err != nil {
		return 0, err
	}
	typ := reflect.ValueOf(base)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	b, err := getValue(typ, expr)
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

// 拆解rule，支持的计算类型+-*/， && ||，其他报错
// 支持二元操作

// 从struct解析找到json Tag, 若嵌套struct则用“.”连接
func getValueByTag(x reflect.Value, tag string) (interface{}, error) {
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		return x, ErrTypeNotStruct
	}
	t := x.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if json, ok := field.Tag.Lookup("json"); ok {
			js := strings.Split(json, ",")
			for _, j := range js {
				if j == tag {
					return x.Field(i).Interface(), nil
				}
			}
		}
	}
	return nil, ErrNotFoundTag
}

func getSliceValue(x reflect.Value, idx int) (interface{}, error) {
	if x.Kind() != reflect.Slice && x.Kind() != reflect.Array {
		return nil, errors.New("only slice or array can get value by index")
	}
	if idx > x.Len()-1 {
		return nil, errors.New("slice index out of range")
	}
	return x.Index(idx).Interface(), nil

}

func getValue(base reflect.Value, expr ast.Expr) (interface{}, error) {
	nullValue := reflect.Value{}
	switch t := expr.(type) {
	case *ast.BinaryExpr:
		x, err := getValue(base, t.X)
		if err != nil {
			return nullValue, err
		}
		y, err := getValue(base, t.Y)
		if err != nil {
			return nullValue, err
		}
		return operate(x, y, t.Op)
	case *ast.Ident:
		return getValueByTag(base, t.Name)
	case *ast.BasicLit:
		return strconv.ParseFloat(t.Value, 64)
	case *ast.ParenExpr:
		return getValue(base, t.X)
	case *ast.SelectorExpr:
		v, err := getValue(base, t.X)
		if err != nil {
			return nullValue, err
		}
		return getValueByTag(reflect.ValueOf(v), t.Sel.Name)
	case *ast.IndexExpr:
		idx, err := getValue(base, t.Index)
		if err != nil {
			return nullValue, err
		}
		f, ok := idx.(float64)
		if !ok {
			return nullValue, errors.New("index must be int or float")
		}
		v, err := getValue(base, t.X)
		if err != nil {
			return nullValue, err
		}
		return getSliceValue(reflect.ValueOf(v), int(f))
	default:
		return nullValue, ErrUnsupportExpr
	}
}

func operate(x, y interface{}, tk token.Token) (interface{}, error) {
	switch tk {
	case token.ADD, token.SUB, token.MUL, token.QUO:
		return mathOp(x, y, tk)
	case token.LSS, token.GTR, token.LEQ, token.GEQ, token.EQL, token.NEQ:
		return compare(x, y, tk)
	case token.LAND, token.LOR:
		return boolOp(x, y, tk)
	default:
		return nil, ErrUnsupportToken
	}
}
