# go-rules

## install 
go get github.com/whyming/go-rules

## Usage

Demo：
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
	exp := "a-b > c"
	r, err := gorules.Bool(abc, exp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("exp result is %t\n", r)

	exp1 := "a/b+c"
	r1, err := gorules.Float(abc, exp1)
	if err != nil {
		fmt.Printf("exp1 err is %v", err)
		return
	}

	fmt.Printf("exp1 result is %f\n", r1)
}

```