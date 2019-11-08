# go-rules

简单的规则引擎？

## install 
go get github.com/whyming/go-rules

## Usage

Demo：

#### 基本用法
```go
package main

import (
	"fmt"

	gorules "github.com/whyming/go-rules"
)

type Abc struct {
	A int64   `json:"a,omitempty"`
	B float64 `json:"b,omitempty"`
	C int     `json:"c,omitempty"`
}

func main() {
	abc := Abc{
		A: 128,
		B: 36.58,
		C: 19,
	}
	// 用法1
	rule, err := gorules.NewRule("a-b > c")
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := rule.Bool(abc)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("exp result is %t\n", r) // true

	// 用法2
	exp1 := "a/b+c"
	r1, err := gorules.Float(abc, exp1)
	if err != nil {
		fmt.Printf("exp1 err is %v", err)
		return
	}

	fmt.Printf("exp1 result is %f\n", r1) // 22.499180
}
```
#### 数组
可以通过下标找到数组的元素，也可以判断元素是否在数组中
```go
	type Abc struct {
		A int64    `json:"a,omitempty"`
		B string   `json:"b,omitempty"`
		C []int64  `json:"c,omitempty"`
		D []string `json:"d,omitempty"`
	}
	a := Abc{
		A: 3,
		B: "xxx",
		C: []int64{20, 197, 5, 12, 7},
		D: []string{"abc", "def", "xxx", "jqk"},
	}
	rule, err := gorules.NewRule("c[2]>a && c[a]>10 && in(d,b)")
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := rule.Bool(a)
	fmt.Println(r, err)
```