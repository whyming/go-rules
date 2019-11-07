package gorules

import (
	"go/token"
	"testing"
)

func Test_mathOp(t *testing.T) {
	type args struct {
		x  interface{}
		y  interface{}
		tk token.Token
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "+",
			args: args{
				x:  5,
				y:  10,
				tk: token.ADD,
			},
			want: 15,
		}, {
			name: "-",
			args: args{
				x:  5,
				y:  10,
				tk: token.SUB,
			},
			want: -5,
		}, {
			name: "*",
			args: args{
				x:  5,
				y:  10,
				tk: token.MUL,
			},
			want: 50,
		}, {
			name: "/",
			args: args{
				x:  5,
				y:  10,
				tk: token.QUO,
			},
			want: 0.5,
		}, {
			name: "x/0",
			args: args{
				x:  5,
				y:  0,
				tk: token.QUO,
			},
			wantErr: true,
		}, {
			name: "&",
			args: args{
				x:  5,
				y:  0,
				tk: token.AND,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mathOp(tt.args.x, tt.args.y, tt.args.tk)
			if (err != nil) != tt.wantErr {
				t.Errorf("mathOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mathOp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compare(t *testing.T) {
	type args struct {
		x  interface{}
		y  interface{}
		tk token.Token
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "<",
			args: args{
				x:  5,
				y:  8,
				tk: token.LSS,
			},
			want: true,
		}, {
			name: ">",
			args: args{
				x:  5,
				y:  8,
				tk: token.GTR,
			},
			want: false,
		}, {
			name: "<=",
			args: args{
				x:  5,
				y:  8,
				tk: token.LEQ,
			},
			want: true,
		}, {
			name: ">=",
			args: args{
				x:  5,
				y:  8,
				tk: token.GEQ,
			},
			want: false,
		}, {
			name: "==",
			args: args{
				x:  8,
				y:  10,
				tk: token.EQL,
			},
			want: false,
		}, {
			name: "=!",
			args: args{
				x:  8,
				y:  8,
				tk: token.NEQ,
			},
			want: false,
		}, {
			name: "=!",
			args: args{
				x:  5,
				y:  8,
				tk: token.ELLIPSIS,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compare(tt.args.x, tt.args.y, tt.args.tk)
			if (err != nil) != tt.wantErr {
				t.Errorf("compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_number(t *testing.T) {
	type args struct {
		x interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "int8",
			args: args{
				x: int8(12),
			},
			want: float64(12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := number(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("number() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("number() = %v, want %v", got, tt.want)
			}
		})
	}
}
