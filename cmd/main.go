package main

import (
	"fmt"

	expr "github.com/hanslad/expr"
)

func exprTestM(str string, expectedValue interface{}) {
	//Step 1: Create environment
	env := expr.NewEnvironment()
	//Step 2: optionally extend environment
	//
	//Step 3: Get parser
	parser := env.GetParser()

	fmt.Printf("Testing %s. Expected value: %v\n", str, expectedValue)

	//Step 4: Parse the input
	ex, err := parser.Parse(str)
	if err != nil {
		fmt.Printf("Parse failed: %v\n", ex)
		fmt.Println(err.Error())
		return

	}
	//
	exEv, errEv := ex.Evaluate(env)
	if errEv != nil {
		fmt.Printf("Eval failed: %v\n", exEv)
		fmt.Println(err.Error())
		return

	}
	if expectedValue != exEv.Value() {
		fmt.Printf("Failed: %v\n", exEv)
		fmt.Printf("Evaluated: %v\n", exEv)
		return
	}
	fmt.Println("Test passed")
}

func main() {
	//exprTestM("\"expr\" + \" \" + \"rules!\"", "expr rules!")
	exprTestM("\"lunch\" in {\"breakfast\", \"lunch\", \"dinner\",\"supper\"}", true)
	//exprTestM("\"lunch\"", "lunch")

	/*	exprTestM("true && (true && true)", true)
			exprTestM("true == true", true)
		  	exprTestM("true != true", false)
		  	exprTestM("false != true", true)
		  	exprTestM("false == true", false)
		  	exprTestM("true != false", true)
		  	exprTestM("true == false", false)
		  	exprTestM("true && false", true)
		  	exprTestM("true && false", true) */

}
