# expr

[![Go build](https://github.com/hanslad/expr/actions/workflows/go-build.yml/badge.svg)](https://github.com/hanslad/expr/actions/workflows/go-build.yml)


Expr is a small extendable script language.  
New expressions(go functions) can be registered in the script environment and extend the basic language with new functions.

## Environment
To parse and run scrips, an Environment need to be created.  
The Environment holds record of all expressions registered and provide the parser for the script input.
Environment is set up with a basic set of expressions by calling RegisterBuiltins().
Environment are further extendable by RegisterSymbol(), RegisterFunction() and RegisterScopedFunction().

## Parser
The Parser has one simple task: Parse the script input to create an epression tree. 
The expression must be evaluated after the parse to 'run' the script.  
During the evaluation, all callback registered in the Environment are call(if used in the script).

## Example
Full example is found in /cmd/main.go

```golang 
// Step 1: Create environment
env := expr.NewEnvironment()
```
Create an function that will be called when the script is evaluated. 
```golang 
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
```
Register the function in the Environment.
```golang 
// Step 2: Extend environment with our multiply method to make "Multiply" a method in the script language.
env.RegisterFunction("Multiply", Multiply)
```
Continue with getting the parser and parse the script.
The script string below consists of a sub expression which is compared with the string "Yeah!"  
The sub expression consists of a conditional expression(condition?left:right).  
The condition is 'Multiply(9,10) == 90', which is a logical expression where the left part is the function that we registered in the Environment earlier.
```golang 
// Step 3: Get parser
parser := env.GetParser()
// Step 4: Parse the input
ex, err := parser.Parse("(Multiply(9,10) == 90 ? 'Yeah!': 'Oh No!') == 'Yeah!'")
if err != nil {
    fmt.Printf("Parse failed: %v. \n\n Error: %s\n", ex, err.Error())
    return
}
```
Finally, evaluate the script to 'run' all expressions. This step will call the registered 'Multiply' function.  
After the evaluation, the evaluated expression will provide a value.

```golang 
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
```

