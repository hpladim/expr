package main

import (
	"fmt"

	expr "github.com/hanslad/expr"
)

func exprtest() {
	ex, err := expr.Parse("true == true")
	if err != nil {
		fmt.Println(err.Error())
	}
	ex = ex.Evaluate(expr.NewEnvironment())

	fmt.Println()

}

func main() {

	exprtest()

}
