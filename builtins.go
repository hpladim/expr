//A set of native function that is used in NativeFunctionExpr or ScopedNativeFunctionExpr
//Has followin profile: "type NativeCallBack func(env *Environment, args []Expression) (Expression, error)"

package expr

import (
	"fmt"
)

// Print is a built-in expression for simple print to console
func Print(env *Environment, args []Expression) (Expression, error) {
	for _, a := range args {
		fmt.Printf("%s", a)
	}
	return env.Get("null"), nil
}
