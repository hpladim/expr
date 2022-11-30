package main

import (
	"fmt"

	expr "github.com/hanslad/expr"
)

// Multiply is an example of an expression that could extend functionality of the script language.
func Multiply(env *expr.Environment, args []expr.Expression) (expr.Expression, error) {
	if len(args) != 2 {
		return env.Get("null"), fmt.Errorf("Multiply must have two int arguments")
	}
	var l = args[0]
	var r = args[1]
	switch t := l.Value().(type) {
	case int64:
		vl, lok := l.Value().(int64)
		vr, rok := r.Value().(int64)
		if lok && rok {
			return expr.NewScalarExpr("", vl*vr), nil
		}
		return env.Get("null"), fmt.Errorf("Multiply must have two int arguments")
	default:
		return env.False(), fmt.Errorf("scalar datatype not supported: %T", t)
	}
}

func main() {

	// Expr is a small extendable script language.
	// New expressions(go functions) can be registered in the script environment and  extend the basic language with new functions.
	// The script string below consists of a sub expression which is compared with the string "Yeah!"
	// The sub expression consists of a conditional expression(condition?left:right).
	// The condition is "Multiply(9,10) == 90" which is a logical expression where the left part is a function that we need to register in the environment.
	var script = "(Multiply(9,10) == 90 ? 'Yeah!': 'Oh No!') == 'Yeah!'"

	// Step 1: Create environment
	env := expr.NewEnvironment()
	// Step 2: Extend environment with our multiply method to make "Multiply" a method in the script language.
	env.RegisterFunction("Multiply", Multiply)
	// Step 3: Get parser
	parser := env.GetParser()
	// Step 4: Parse the input
	ex, err := parser.Parse(script)
	if err != nil {
		fmt.Printf("Parse failed: %v. \n\n Error: %s\n", ex, err.Error())
		return
	}
	// Step 5: Evaluate the parsed text
	exEv, errEv := ex.Evaluate(env)
	if errEv != nil {
		fmt.Printf("Eval failed:\n\n %v \n\n Evaluated to:\n\n %v \n\n Error: %s\n", ex, exEv, errEv.Error())
		return
	}
	// Step 6: Check result
	val, _ := exEv.Value().(bool)
	if !val {
		fmt.Printf("Check failed:\n\n %v \n\n Evaluated to:\n\n %v \n", ex, exEv.Value())
		return
	}
}
