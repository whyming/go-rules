package gorules

import (
	"go/ast"
	"go/parser"
	"reflect"
	"testing"
)

func Test_getValueByTag(t *testing.T) {
	type args struct {
		x   reflect.Value
		tag string
	}
	type WithInt struct {
		IntValue int64 `json:"int_value"`
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "int",
			args: args{
				x:   reflect.ValueOf(WithInt{123}),
				tag: "int_value",
			},
			want: int64(123),
		},
		{
			name: "not found",
			args: args{
				x:   reflect.ValueOf(WithInt{123}),
				tag: "xx",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValueByTag(tt.args.x, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValueByTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValueByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValue(t *testing.T) {
	type args struct {
		base reflect.Value
		expr ast.Expr
	}
	type Abc struct {
		A int64 `json:"a"`
		B int64 `json:"b"`
	}
	type Xy struct {
		X   float64  `json:"x,omitempty"`
		Abc Abc      `json:"abc,omitempty"`
		Y   []int64  `json:"y,omitempty"`
		Z   []string `json:"z,omitempty"`
	}
	type More struct {
		Xy Xy `json:"xy,omitempty"`
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "easy-add",
			args: args{
				base: reflect.ValueOf(Abc{A: 10, B: 8}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`a+b`)
					return expr
				}(),
			},
			want: float64(18),
		},
		{
			name: "multi-sub-add-mul-div",
			args: args{
				base: reflect.ValueOf(Abc{A: 10, B: 8}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`a+b*b-a*b+(a+b)/(a-b)`)
					return expr
				}(),
			},
			want: float64(3),
		},
		{
			name: "less than",
			args: args{
				base: reflect.ValueOf(Abc{A: 10, B: 8}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`a+b<a*b`)
					return expr
				}(),
			},
			want: true,
		},
		{
			name: "bool",
			args: args{
				base: reflect.ValueOf(Abc{A: 10, B: 8}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`a>b && b<5 || a>8 &&b<9`)
					return expr
				}(),
			},
			want: true,
		}, {
			name: "a.b",
			args: args{
				base: reflect.ValueOf(Xy{X: 10, Abc: Abc{A: 8, B: 20}}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`abc.a-abc.b+x`)
					return expr
				}(),
			},
			want: float64(-2),
		}, {
			name: "a.b.c",
			args: args{
				base: reflect.ValueOf(More{Xy: Xy{X: 10, Abc: Abc{A: 8, B: 20}}}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`xy.abc.a-xy.abc.b+xy.x`)
					return expr
				}(),
			},
			want: float64(-2),
		}, {
			name: "a.b[c+d]",
			args: args{
				base: reflect.ValueOf(More{Xy: Xy{X: 10, Abc: Abc{A: 8, B: 20}, Y: []int64{3, 6, 9}}}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`xy.y[1]-xy.y[xy.abc.b/xy.x]`)
					return expr
				}(),
			},
			want: float64(-3),
		}, {
			name: "function IN int64 yes",
			args: args{
				base: reflect.ValueOf(More{Xy: Xy{X: 10, Abc: Abc{A: 8, B: 20}, Y: []int64{3, 6, 9}}}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`in(xy.y,6.0)`)
					return expr
				}(),
			},
			want: true,
		}, {
			name: "function IN int64 no",
			args: args{
				base: reflect.ValueOf(More{Xy: Xy{X: 10, Abc: Abc{A: 8, B: 20}, Y: []int64{3, 6, 9}}}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`in(xy.y,5.0)`)
					return expr
				}(),
			},
			want: false,
		}, {
			name: "function IN string yes",
			args: args{
				base: reflect.ValueOf(More{Xy: Xy{X: 10, Abc: Abc{A: 8, B: 20}, Y: []int64{3, 6, 9}, Z: []string{"abc", "bcd"}}}),
				expr: func() ast.Expr {
					expr, _ := parser.ParseExpr(`in(xy.z,"abc")`)
					return expr
				}(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValue(tt.args.base, tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	type args struct {
		base interface{}
		rule string
	}
	type TypeInt struct {
		A int64 `json:"a"`
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "int",
			args: args{
				base: TypeInt{A: 100},
				rule: "a-5*8",
			},
			want: 60,
		}, {
			name: "null",
			args: args{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int(tt.args.base, tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("Int() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	type args struct {
		base interface{}
		rule string
	}
	type Ftype struct {
		A int64 `json:"a"`
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "float",
			args: args{
				base: Ftype{A: 100},
				rule: "3*a-20",
			},
			want: 280,
		}, {
			name: "int",
			args: args{
				base: &Ftype{A: 10},
				rule: "a",
			},
			want: 10,
		}, {
			name: "null",
			args: args{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Float(tt.args.base, tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("Float() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Float() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	type args struct {
		base interface{}
		rule string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "a+b<c*2",
			args: args{
				base: struct {
					A int64   `json:"a,omitempty"`
					B int32   `json:"b,omitempty"`
					C float64 `json:"c,omitempty"`
				}{
					A: 8,
					B: 12,
					C: 16,
				},
				rule: "a+b<c*2",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bool(tt.args.base, tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}
