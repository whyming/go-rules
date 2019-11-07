package gorules

import (
	"errors"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

// 错误定义
var (
	ErrRuleEmpty      = errors.New("rule is empty")
	ErrTypeNotStruct  = errors.New("value must struct or struct pointer")
	ErrNotFoundTag    = errors.New("not found tag")
	ErrUnsupportToken = errors.New("unsupport token")
	ErrUnsupportExpr  = errors.New("unsupport expr")
	ErrNotNumber      = errors.New("not a number")
	ErrNotBool        = errors.New("not boolean")
)

// Bool 规则rule结果的布尔值，rule的参数基于base的json tag
func Bool(base interface{}, rule string) (bool, error) {
	r, err := NewRule(rule)
	if err != nil {
		return false, err
	}
	return r.Bool(base)
}

// Int 规则rule的结果如果是数值型，转换为int64，否则报错
func Int(base interface{}, rule string) (int64, error) {
	r, err := NewRule(rule)
	if err != nil {
		return 0, err
	}
	return r.Int(base)
}

// Float 返回规则rule结果，如果数值型返回float64，否则报错
func Float(base interface{}, rule string) (float64, error) {
	r, err := NewRule(rule)
	if err != nil {
		return 0, err
	}
	return r.Float(base)
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
		switch t.Kind {
		case token.STRING:
			return strings.Trim(t.Value, "\""), nil
		case token.INT:
			return strconv.ParseInt(t.Value, 10, 64)
		case token.FLOAT:
			return strconv.ParseFloat(t.Value, 64)
		default:
			return nullValue, errors.New("unsupport param")
		}
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
			if i, ok := idx.(int64); ok {
				f = float64(i)
			} else {
				return nullValue, errors.New("index must be int or float")
			}

		}
		v, err := getValue(base, t.X)
		if err != nil {
			return nullValue, err
		}
		return getSliceValue(reflect.ValueOf(v), int(f))
	case *ast.CallExpr:
		if fexp, ok := t.Fun.(*ast.Ident); ok {
			if strings.ToUpper(fexp.Name) == "IN" {
				if len(t.Args) == 2 {
					return isIn(base, t.Args[0], t.Args[1])
				}
				return nullValue, errors.New("function IN only support tow params")
			}
			return nullValue, errors.New("unsupport function: " + fexp.Name)
		}
		return nullValue, errors.New("unknow function")
	default:
		return nullValue, ErrUnsupportExpr
	}
}

func isIn(base reflect.Value, slice ast.Expr, key ast.Expr) (bool, error) {
	sv, err := getValue(base, slice)
	if err != nil {
		return false, err
	}

	kv, err := getValue(base, key)
	if err != nil {
		return false, err
	}
	svv := reflect.ValueOf(sv)
	if svv.Kind() != reflect.Slice && svv.Kind() != reflect.Array {
		return false, errors.New("function IN first param must be slice or array")
	}
	if svv.Len() == 0 {
		return false, nil
	}
	kvv := reflect.ValueOf(kv)

	switch svv.Index(0).Kind() {
	case reflect.String:
		for i := 0; i < svv.Len(); i++ {
			if svv.Index(i).CanInterface() {
				if s, ok := svv.Index(i).Interface().(string); ok {
					if s == kvv.String() {
						return true, nil
					}
				}
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		for i := 0; i < svv.Len(); i++ {
			if svv.Index(i).CanInterface() {
				if s, ok := svv.Index(i).Interface().(float64); ok {
					if s == kvv.Float() {
						return true, nil
					}
				} else if s, ok := svv.Index(i).Interface().(int64); ok {
					if float64(s) == kvv.Float() {
						return true, nil
					}
				}
			}
		}
	default:
		return false, errors.New("function IN only support: string int float")
	}
	return false, nil
}
