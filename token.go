package gorules

import (
	"errors"
	"go/token"
	"reflect"
)

func operate(x, y interface{}, tk token.Token) (interface{}, error) {
	xv := reflect.ValueOf(x)
	yv := reflect.ValueOf(y)
	switch tk {
	case token.ADD, token.SUB, token.MUL, token.QUO:
		return mathOp(xv, yv, tk)
	case token.LSS, token.GTR, token.LEQ, token.GEQ, token.EQL, token.NEQ:
		return compare(xv, yv, tk)
	case token.LAND, token.LOR:
		return boolOp(xv, yv, tk)
	default:
		return nil, ErrUnsupportToken
	}
}

// 数字计算操作
func mathOp(x, y reflect.Value, tk token.Token) (float64, error) {
	numx, err := number(x)
	if err != nil {
		return 0, err
	}
	numy, err := number(y)
	if err != nil {
		return 0, err
	}
	switch tk {
	case token.ADD:
		return numx + numy, nil
	case token.SUB:
		return numx - numy, nil
	case token.MUL:
		return numx * numy, nil
	case token.QUO:
		if numy == 0 {
			return 0, errors.New("x/0 error")
		}
		return numx / numy, nil
	default:
		return 0, ErrUnsupportToken
	}
}

// 数值比较，暂时支持6种 >, <, >=,<=， ==， !=
func compare(x, y reflect.Value, tk token.Token) (bool, error) {
	if x.Kind() == reflect.String && y.Kind() == reflect.String {
		return compareString(x.String(), y.String(), tk)
	}
	numx, err := number(x)
	if err != nil {
		return false, err
	}
	numy, err := number(y)
	if err != nil {
		return false, err
	}
	switch tk {
	case token.LSS:
		return numx < numy, nil
	case token.GTR:
		return numx > numy, nil
	case token.LEQ:
		return numx <= numy, nil
	case token.GEQ:
		return numx >= numy, nil
	case token.EQL:
		return numx == numy, nil
	case token.NEQ:
		return numx != numy, nil
	default:
		return false, ErrUnsupportToken
	}
}

// 布尔操作 && || 两种
func boolOp(x, y reflect.Value, tk token.Token) (bool, error) {
	if x.Kind() != reflect.Bool || y.Kind() != reflect.Bool {
		return false, ErrNotBool
	}

	switch tk {
	case token.LAND:
		return x.Bool() && y.Bool(), nil
	case token.LOR:
		return x.Bool() || y.Bool(), nil
	default:
		return false, ErrUnsupportToken
	}
}

func number(x reflect.Value) (float64, error) {
	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(x.Int()), nil
	case reflect.Float32, reflect.Float64:
		return x.Float(), nil
	default:
		return 0, ErrNotNumber
	}
}

func compareString(x, y string, tk token.Token) (bool, error) {
	switch tk {
	case token.EQL:
		return x == y, nil
	case token.NEQ:
		return x != y, nil
	default:
		return false, ErrUnsupportToken
	}
}

func stringPair(x, y reflect.Value) (string, string, error) {
	return "", "", nil
}
