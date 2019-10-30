package gorules

import (
	"errors"
	"go/token"
)

// 数字计算操作
func mathOp(x, y interface{}, tk token.Token) (float64, error) {
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
func compare(x, y interface{}, tk token.Token) (bool, error) {
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
func boolOp(x, y interface{}, tk token.Token) (bool, error) {
	boolx, ok := x.(bool)
	if !ok {
		return false, ErrNotBool
	}
	booly, ok := y.(bool)
	if !ok {
		return false, ErrNotBool
	}
	switch tk {
	case token.LAND:
		return boolx && booly, nil
	case token.LOR:
		return boolx || booly, nil
	default:
		return false, ErrUnsupportToken
	}
}

func number(x interface{}) (float64, error) {
	switch num := x.(type) {
	case int:
		return float64(num), nil
	case int8:
		return float64(num), nil
	case int16:
		return float64(num), nil
	case int32:
		return float64(num), nil
	case int64:
		return float64(num), nil
	case float32:
		return float64(num), nil
	case float64:
		return num, nil
	default:
		return 0, ErrNotNumber
	}
}
