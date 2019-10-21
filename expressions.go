package expr

import (
	"errors"
	"fmt"
	"strings"
)

//Expression is the generic expression for all expression types in the expression framework
type Expression interface {
	//Evaluate the expression
	Evaluate(env *Environment) (Expression, error)
	//Literal will provide a uniqe literal for the expression
	Literal() string
	//Value will provide value after evaluation
	Value() interface{}
	//String will provide the string representation of value
	String() string
}

//Function is the more specific
type Function interface {
	Expression
	Invoke(env *Environment, args []Expression) (Expression, error)
}

//=============================================================================

//SymbolExpr is  used for attaching functions to extend the Environment
//Symbols can even be registered in a scope!
type SymbolExpr struct {
	name  string
	scope *SymbolExpr
}

//NewSymbolExpr registers a new symbol with the provided name
func NewSymbolExpr(name string) *SymbolExpr {
	sy := SymbolExpr{}
	sy.name = name
	return &sy
}

//NewSymbolExprWithScope registers a new symbol with the provided name
func NewSymbolExprWithScope(name string, scope *SymbolExpr) *SymbolExpr {
	sy := SymbolExpr{}
	sy.name = name
	sy.scope = scope
	return &sy
}

//Evaluate the expression
func (e *SymbolExpr) Evaluate(env *Environment) (Expression, error) {
	return env.Get(e.Literal()), nil
}

//Literal will provide a uniqe literal for the expression
func (e *SymbolExpr) Literal() string {
	if e.scope != nil {
		return e.scope.Literal() + "." + e.name
	}
	return e.name
}

//Value will provide value after evaluation
func (e *SymbolExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *SymbolExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

//ScalarExpr is a basic scalar expression
type ScalarExpr struct {
	literal string
	value   interface{}
}

//NewScalarExpr registers a new scalar with the defined literal and value
func NewScalarExpr(literal string, value interface{}) *ScalarExpr {
	sc := ScalarExpr{}
	sc.value = value
	sc.literal = literal
	return &sc
}

//NewScalarExprV registers a new scalar and creates and literal based on the value
func NewScalarExprV(value interface{}) *ScalarExpr {
	sc := ScalarExpr{}
	sc.value = value
	if sc.value != nil {
		switch sc.value.(type) {
		case string:
			sc.literal = fmt.Sprintf("\"%s\"", escape(sc.value.(string)))
		default:
			sc.literal = fmt.Sprintf("%s", sc.value)
		}
	} else {
		sc.literal = "null"
	}
	return &sc
}

func escape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"")
}

//Evaluate the expression
func (e *ScalarExpr) Evaluate(env *Environment) (Expression, error) {
	return e, nil
}

//Literal will provide a uniqe literal for the expression
func (e *ScalarExpr) Literal() string {
	return e.literal
}

//Value will provide value after evaluation
func (e *ScalarExpr) Value() interface{} {
	return e.value
}

//String will provide the string representation of value
func (e *ScalarExpr) String() string {
	if e.value != nil {
		return fmt.Sprintf("%v", e.value)
	}
	return "NULL"
}

//=============================================================================

//CondExpr is a conditional expression on the form cond? left: right
type CondExpr struct {
	condition Expression
	left      Expression
	right     Expression
}

//NewCondExpr registers a new conditional expression on the form cond?left:right
func NewCondExpr(cond Expression, left Expression, right Expression) *CondExpr {
	e := CondExpr{}
	e.condition = cond
	e.left = left
	e.right = right
	return &e
}

//Evaluate the expression
func (e *CondExpr) Evaluate(env *Environment) (Expression, error) {
	c, err := e.condition.Evaluate(env)
	if err != nil {
		return c, err
	}
	if c == nil || c == env.False() || c.Value() == nil || c.Value() == env.False() {
		return e.left.Evaluate(env)
	}
	return e.right.Evaluate(env)
}

//Literal will provide a uniqe literal for the expression
func (e *CondExpr) Literal() string {
	return fmt.Sprintf("(%s?%s:%s)", e.condition.Literal(), e.left.Literal(), e.right.Literal())
}

//Value will provide value after evaluation
func (e *CondExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *CondExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

//OrExpr is a basic binary expression(||)
type OrExpr struct {
	left  Expression
	right Expression
}

//NewOrExpr registers a new conditional expression on the form left||right
func NewOrExpr(left Expression, right Expression) *OrExpr {
	e := OrExpr{}
	e.left = left
	e.right = right
	return &e
}

//Evaluate the expression
func (e *OrExpr) Evaluate(env *Environment) (Expression, error) {
	c, err := e.left.Evaluate(env)
	if err != nil {
		return c, err
	}
	if c != nil && c != env.False() && c.Value() != nil && c.Value() != env.False() {
		return env.True(), nil
	}
	c, err = e.right.Evaluate(env)
	if err != nil {
		return c, err
	}
	if c != nil && c != env.False() && c.Value() != nil && c.Value() != env.False() {
		return env.True(), nil
	}
	return env.False(), nil
}

//Literal will provide a uniqe literal for the expression
func (e *OrExpr) Literal() string {
	return fmt.Sprintf("(%s || %s)", e.left.Literal(), e.right.Literal())
}

//Value will provide value after evaluation
func (e *OrExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *OrExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

//AndExpr is a basic binary expression(&&)
type AndExpr struct {
	left  Expression
	right Expression
}

//NewAndExpr registers a new conditional expression on the form left||right
func NewAndExpr(left Expression, right Expression) *AndExpr {
	e := AndExpr{}
	e.left = left
	e.right = right
	return &e
}

//Evaluate the expression
func (e *AndExpr) Evaluate(env *Environment) (Expression, error) {
	c, err := e.left.Evaluate(env)
	if err != nil {
		return c, err
	}
	if c == nil || c == env.False() || c.Value() == nil || c.Value() == env.False() {
		return env.False(), nil
	}
	c, err = e.right.Evaluate(env)
	if err != nil {
		return c, err
	}
	if c == nil || c == env.False() || c.Value() == nil || c.Value() == env.False() {
		return env.False(), nil
	}
	return env.True(), nil
}

//Literal will provide a uniqe literal for the expression
func (e *AndExpr) Literal() string {
	return fmt.Sprintf("(%s && %s)", e.left.Literal(), e.right.Literal())
}

//Value will provide value after evaluation
func (e *AndExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *AndExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

//CompareExpr is a basic compare expression handling the following compare operands:
// '==', '!=', '>=','>','<=','<'
// only scalar expression with same value type is compared
type CompareExpr struct {
	operand string
	left    Expression
	right   Expression
}

//NewCompareExpr registers a new compare expression on the form left comparator right
func NewCompareExpr(operand string, left Expression, right Expression) *CompareExpr {
	e := CompareExpr{}
	e.operand = operand
	e.left = left
	e.right = right

	return &e
}

//Evaluate the expression, will do a compare supporting the following operands: '==', '!=', '>=','>','<=','<'
//Only scalar expression with same value type is compared
func (e *CompareExpr) Evaluate(env *Environment) (Expression, error) {

	l, err := e.left.Evaluate(env)
	if err != nil {
		return l, err
	}
	r, err := e.right.Evaluate(env)
	if err != nil {
		return r, err
	}
	ls, lok := l.(*ScalarExpr)
	rs, rok := r.(*ScalarExpr)
	if lok && rok {
		return compare(env, e.operand, *ls, *rs)
	}
	return nil, errors.New("Equality-operator is only supported on Scalar values")
}

//Literal will provide a uniqe literal for the expression
func (e *CompareExpr) Literal() string {
	return fmt.Sprintf("(%s %s %s)", e.left.Literal(), e.operand, e.right.Literal())
}

//Value will provide value after evaluation
func (e *CompareExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *CompareExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

//ConCatExpr is a basic concatenation expression
type ConCatExpr struct {
	left  Expression
	right Expression
}

//NewConcatExpr registers a new concatenation expression on the form left + right
func NewConcatExpr(left Expression, right Expression) *ConCatExpr {
	e := ConCatExpr{}
	e.left = left
	e.right = right
	return &e
}

//Evaluate the expression
func (e *ConCatExpr) Evaluate(env *Environment) (Expression, error) {
	l, err := e.left.Evaluate(env)
	if err != nil {
		return l, err
	}
	r, err := e.right.Evaluate(env)
	if err != nil {
		return r, err
	}
	if l == nil || l.Value() == nil {
		if r == nil || r.Value() == nil {

			return NewScalarExpr("", ""), nil
		}
		return NewScalarExpr("", r.String()), nil
	}

	if r == nil || r.Value() == nil {
		return NewScalarExpr("", r.String()), nil
	}

	return NewScalarExpr("", l.String()+r.String()), nil
}

//Literal will provide a uniqe literal for the expression
func (e *ConCatExpr) Literal() string {
	return fmt.Sprintf("(%s%s)", e.left.Literal(), e.right.Literal())
}

//Value will provide value after evaluation
func (e *ConCatExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *ConCatExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

//FuncCallExpr is a expression holding a function an its arguments
//Fantastic stuff
type FuncCallExpr struct {
	function Expression
	args     []Expression
}

//NewFuncCallExpr registers a new concatenation expression on the form left + right
func NewFuncCallExpr() *FuncCallExpr {
	e := FuncCallExpr{}
	e.args = make([]Expression, 0)
	return &e
}

//SetFunc sets the contained function in FuncExpr
func (e *FuncCallExpr) SetFunc(expr Expression) {
	e.function = expr
}

//GetFunc gets the contained function in FuncExpr
func (e *FuncCallExpr) GetFunc() Expression {
	return e.function
}

//AddArg adds argument to function
func (e *FuncCallExpr) AddArg(expr Expression) {
	e.args = append(e.args, expr)
}

//GetArgs returns all the registered args on the function
func (e *FuncCallExpr) GetArgs() []Expression {
	return e.args
}

//Evaluate the expression
func (e *FuncCallExpr) Evaluate(env *Environment) (Expression, error) {
	args := make([]Expression, len(e.GetArgs()))
	f, err := e.GetFunc().Evaluate(env)
	if err != nil {
		return f, err
	}
	fun, ok := f.(Function)
	if ok {
		for _, aex := range e.GetArgs() {
			a, err := aex.Evaluate(env)
			if err != nil {
				return a, err
			}
			args = append(args, a)
		}
		env.pushStack(exFrame{function: fun, args: e.GetArgs()})
		defer env.popStack()
		return fun.Invoke(env, args)

	}

	return f, fmt.Errorf("Not a function: %s", e.function.Literal())

}

//Literal will provide a uniqe literal for the expression
func (e *FuncCallExpr) Literal() string {
	return fmt.Sprintf("(#invoke-function:%s#)", e.function.Literal())
}

//Value will provide value after evaluation
func (e *FuncCallExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *FuncCallExpr) String() string {
	return fmt.Sprintf("%T", e)
}
