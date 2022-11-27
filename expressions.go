package expr

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Expression is the generic expression for all expression types in the expression framework
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

// Function is the more specific
type Function interface {
	Expression
	Invoke(env *Environment, args []Expression) (Expression, error)
}

//=============================================================================

// SymbolExpr is  used for attaching functions to extend the Environment
// Symbols can even be registered in a scope!
type SymbolExpr struct {
	name  string
	scope *SymbolExpr
}

// NewSymbolExpr registers a new symbol with the provided name
func NewSymbolExpr(name string) *SymbolExpr {
	sy := SymbolExpr{}
	sy.name = name
	return &sy
}

// NewSymbolExprWithScope registers a new symbol with the provided name
func NewSymbolExprWithScope(name string, scope *SymbolExpr) *SymbolExpr {
	sy := SymbolExpr{}
	sy.name = name
	sy.scope = scope
	return &sy
}

// Evaluate the expression
func (e *SymbolExpr) Evaluate(env *Environment) (Expression, error) {
	return env.Get(e.Literal()), nil
}

// Literal will provide a uniqe literal for the expression
func (e *SymbolExpr) Literal() string {
	if e.scope != nil {
		return e.scope.Literal() + "." + e.name
	}
	return e.name
}

// Value will provide value after evaluation
func (e *SymbolExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *SymbolExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// ScalarExpr is a basic scalar expression
type ScalarExpr struct {
	literal string
	value   interface{}
}

// NewScalarExpr registers a new scalar with the defined literal and value
func NewScalarExpr(literal string, value interface{}) *ScalarExpr {
	sc := ScalarExpr{}
	sc.value = value
	sc.literal = literal
	return &sc
}

// NewScalarExprV registers a new scalar and creates and literal based on the value
func NewScalarExprV(value interface{}) *ScalarExpr {
	sc := ScalarExpr{}
	sc.value = value
	if sc.value != nil {
		switch sc.value.(type) {
		case string:
			sc.literal = fmt.Sprintf("\"%v\"", escape(sc.value.(string)))
		default:
			sc.literal = fmt.Sprintf("%v", sc.value)
		}
	} else {
		sc.literal = "null"
	}
	return &sc
}

func escape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"")
}

// Evaluate the expression
func (e *ScalarExpr) Evaluate(env *Environment) (Expression, error) {
	return e, nil
}

// Literal will provide a uniqe literal for the expression
func (e *ScalarExpr) Literal() string {
	return e.literal
}

// Value will provide value after evaluation
func (e *ScalarExpr) Value() interface{} {
	return e.value
}

// String will provide the string representation of value
func (e *ScalarExpr) String() string {
	if e.value != nil {
		return fmt.Sprintf("%v", e.value)
	}
	return "NULL"
}

//=============================================================================

// CondExpr is a conditional expression on the form cond? left: right
type CondExpr struct {
	condition Expression
	left      Expression
	right     Expression
}

// NewCondExpr registers a new conditional expression on the form cond?left:right
func NewCondExpr(cond Expression, left Expression, right Expression) *CondExpr {
	e := CondExpr{}
	e.condition = cond
	e.left = left
	e.right = right
	return &e
}

// Evaluate the expression
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

// Literal will provide a uniqe literal for the expression
func (e *CondExpr) Literal() string {
	return fmt.Sprintf("(%s?%s:%s)", e.condition.Literal(), e.left.Literal(), e.right.Literal())
}

// Value will provide value after evaluation
func (e *CondExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *CondExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// OrExpr is a basic binary expression(||)
type OrExpr struct {
	left  Expression
	right Expression
}

// NewOrExpr registers a new conditional expression on the form left||right
func NewOrExpr(left Expression, right Expression) *OrExpr {
	e := OrExpr{}
	e.left = left
	e.right = right
	return &e
}

// Evaluate the expression
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

// Literal will provide a uniqe literal for the expression
func (e *OrExpr) Literal() string {
	return fmt.Sprintf("(%s || %s)", e.left.Literal(), e.right.Literal())
}

// Value will provide value after evaluation
func (e *OrExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *OrExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// AndExpr is a basic binary expression(&&)
type AndExpr struct {
	left  Expression
	right Expression
}

// NewAndExpr registers a new conditional expression on the form left||right
func NewAndExpr(left Expression, right Expression) *AndExpr {
	e := AndExpr{}
	e.left = left
	e.right = right
	return &e
}

// Evaluate the expression
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

// Literal will provide a uniqe literal for the expression
func (e *AndExpr) Literal() string {
	return fmt.Sprintf("(%s && %s)", e.left.Literal(), e.right.Literal())
}

// Value will provide value after evaluation
func (e *AndExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *AndExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// CompareExpr is a basic compare expression handling the following compare operands:
// '==', '!=', '>=','>','<=','<'
// only scalar expression with same value type is compared
type CompareExpr struct {
	operand string
	left    Expression
	right   Expression
}

// NewCompareExpr registers a new compare expression on the form left comparator right
func NewCompareExpr(operand string, left Expression, right Expression) *CompareExpr {
	e := CompareExpr{}
	e.operand = operand
	e.left = left
	e.right = right

	return &e
}

// Evaluate the expression, will do a compare supporting the following operands: '==', '!=', '>=','>','<=','<'
// Only scalar expression with same value type is compared
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
		return compare(env, e.operand, ls, rs)
	}
	return nil, errors.New("equality-operator is only supported on Scalar values")
}

// Literal will provide a uniqe literal for the expression
func (e *CompareExpr) Literal() string {
	return fmt.Sprintf("(%s %s %s)", e.left.Literal(), e.operand, e.right.Literal())
}

// Value will provide value after evaluation
func (e *CompareExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *CompareExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// ConCatExpr is a basic concatenation expression
type ConCatExpr struct {
	left  Expression
	right Expression
}

// NewConcatExpr registers a new concatenation expression on the form left + right
func NewConcatExpr(left Expression, right Expression) *ConCatExpr {
	e := ConCatExpr{}
	e.left = left
	e.right = right
	return &e
}

// Evaluate the expression
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

// Literal will provide a uniqe literal for the expression
func (e *ConCatExpr) Literal() string {
	return fmt.Sprintf("(%s%s)", e.left.Literal(), e.right.Literal())
}

// Value will provide value after evaluation
func (e *ConCatExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *ConCatExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// ListExpr is a epression for list definitions
type ListExpr struct {
	exprs []Expression
}

// NewListExpr registers a new list expression in form [expr,....]
// Neat?
func NewListExpr() *ListExpr {
	e := ListExpr{}
	e.exprs = make([]Expression, 0)
	return &e
}

// Append add one expression to end of expression list(slice)
func (e *ListExpr) Append(expr Expression) {
	e.exprs = append(e.exprs, expr)
}

// Count return length of expression list(slice)
func (e *ListExpr) Count() int {
	return len(e.exprs)
}

// Get the expression at index, careful! Preferrably used together with Count()
// Does not check for out of bound referencing
func (e *ListExpr) Get(idx int) Expression {
	return e.exprs[idx]
}

// Evaluate the expression
func (e *ListExpr) Evaluate(env *Environment) (Expression, error) {
	expr := NewListExpr()
	for _, ex := range e.exprs {
		le, err := ex.Evaluate(env)
		if err != nil {
			return le, err
		}
		expr.Append(le)
	}
	return expr, nil

}

// Literal will provide a uniqe literal for the expression
func (e *ListExpr) Literal() string {
	var sb strings.Builder
	prefix := "{"
	for _, ex := range e.exprs {
		sb.WriteString(prefix)
		sb.WriteString(ex.Literal())
		prefix = ", "
	}
	return sb.String()
}

// Value will provide value after evaluation
func (e *ListExpr) Value() interface{} {
	return e
}

// String will provide the string representation of value
func (e *ListExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// InExpr is a expression for list definitions
type InExpr struct {
	left  Expression
	right Expression
}

// NewInExpr registers a new list expression in form: left in [expr...]
// Neat?
func NewInExpr(left Expression, right Expression) *InExpr {
	e := InExpr{}
	e.left = left
	e.right = right
	return &e
}

// Evaluate the expression
func (e *InExpr) Evaluate(env *Environment) (Expression, error) {
	val, err := e.left.Evaluate(env)
	if err != nil {
		return val, err
	}
	list, err := e.right.Evaluate(env)
	if err != nil {
		return list, err
	}
	listex, ok := list.(*ListExpr)
	if !ok {
		return list, fmt.Errorf("not a list for matching values: %s", e.Literal())
	}
	str := fmt.Sprintf("%v", val.Value())
	for _, ex := range listex.exprs {
		if str == fmt.Sprintf("%v", ex.Value()) {
			return env.True(), nil
		}
	}
	return env.False(), nil
}

// Literal will provide a uniqe literal for the expression
func (e *InExpr) Literal() string {
	return fmt.Sprintf("(%s in %s)", e.left.Literal(), e.right.Literal())
}

// Value will provide value after evaluation
func (e *InExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *InExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// LikeExpr is a basic binary expression(&&)
type LikeExpr struct {
	left  Expression
	right Expression
}

// NewLikeExpr registers a new conditional expression on the left like right
func NewLikeExpr(left Expression, right Expression) *LikeExpr {
	e := LikeExpr{}
	e.left = left
	e.right = right
	return &e
}

// Evaluate the expression
func (e *LikeExpr) Evaluate(env *Environment) (Expression, error) {
	l, err := e.left.Evaluate(env)
	if err != nil {
		return l, err
	}
	r, err := e.right.Evaluate(env)
	if err != nil {
		return r, err
	}
	ls := fmt.Sprintf("%v", l.Value())
	rs := fmt.Sprintf("%v", r.Value())
	if Match(ls, rs) {
		return env.True(), nil
	}
	return env.False(), nil
}

// Match checks for wildcard match in string
func Match(txt string, pattern string) bool {
	return MatchI(txt, 0, pattern, 0)
}

// MatchI checks for wildcard match in string at index
func MatchI(txt string, ti int, pattern string, pi int) bool {
	tl := len(txt)
	pl := len(pattern)
	for ti < tl && pi < pl {
		r := unicode.ToLower([]rune(pattern)[pi])
		pi++
		switch r {

		case '*':

			if pi == pl {
				return true
			}
			for ti < tl && !MatchI(txt, ti, pattern, pi) {
				ti++
			}
			return (ti < tl)

		case '?':

			ti++

		default:
			if r != unicode.ToLower([]rune(txt)[ti]) {
				return false
			}

			ti++
		}
	}
	return (ti == tl && (pi == pl || (pi == pl-1 && pattern[pi] == '*')))
}

// Literal will provide a uniqe literal for the expression
func (e *LikeExpr) Literal() string {
	return fmt.Sprintf("(%s && %s)", e.left.Literal(), e.right.Literal())
}

// Value will provide value after evaluation
func (e *LikeExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *LikeExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// FuncCallExpr is a expression holding a function and its arguments
// Fantastic stuff
type FuncCallExpr struct {
	function Expression
	args     []Expression
}

// NewFuncCallExpr registers a new concatenation expression on the form left + right
func NewFuncCallExpr() *FuncCallExpr {
	e := FuncCallExpr{}
	e.args = make([]Expression, 0)
	return &e
}

// SetFunc sets the contained function in FuncExpr
func (e *FuncCallExpr) SetFunc(expr Expression) {
	e.function = expr
}

// GetFunc gets the contained function in FuncExpr
func (e *FuncCallExpr) GetFunc() Expression {
	return e.function
}

// AddArg adds argument to function
func (e *FuncCallExpr) AddArg(expr Expression) {
	e.args = append(e.args, expr)
}

// GetArgs returns all the registered args on the function
func (e *FuncCallExpr) GetArgs() []Expression {
	return e.args
}

// Evaluate the expression
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

	return f, fmt.Errorf("not a function: %s", e.function.Literal())

}

// Literal will provide a uniqe literal for the expression
func (e *FuncCallExpr) Literal() string {
	return fmt.Sprintf("(#invoke-function:%s#)", e.function.Literal())
}

// Value will provide value after evaluation
func (e *FuncCallExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *FuncCallExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// ScopedFuncCallExpr is a expression holding a scoped function and its arguments
// Fantastic stuff
type ScopedFuncCallExpr struct {
	name  string
	scope Expression
	args  []Expression
}

// NewScopedFuncCallExpr registers a new scoped function in the form symbol.symbol(args)
func NewScopedFuncCallExpr(name string, scope Expression) (*ScopedFuncCallExpr, error) {
	//TODO: What happens if scope is nil?
	_, sok := scope.(*SymbolExpr)
	_, sfok := scope.(*ScopedFuncCallExpr)
	if !(sok || sfok) {
		return nil, errors.New("the provided scope must be a symbol or scoped function")
	}
	e := ScopedFuncCallExpr{}
	e.name = name
	e.scope = scope
	e.args = make([]Expression, 0)
	return &e, nil
}

// AddArg adds argument to function
func (e *ScopedFuncCallExpr) AddArg(expr Expression) {
	e.args = append(e.args, expr)
}

// GetArgs returns all the registered args on the function
func (e *ScopedFuncCallExpr) GetArgs() []Expression {
	return e.args
}

// Evaluate the expression
func (e *ScopedFuncCallExpr) Evaluate(env *Environment) (Expression, error) {

	args := make([]Expression, len(e.GetArgs()))

	f, err := env.Get(e.Literal()).Evaluate(env)

	if err != nil {
		return f, err
	}
	fun, ok := f.(*ScopedNativeFunctionExpr)
	if ok {
		evscope, scerr := e.scope.Evaluate(env)
		if scerr != nil {
			return f, scerr
		}
		args = append(args, evscope)
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

	return f, fmt.Errorf("not a function: %s", f.Literal())

}

// Literal will provide a uniqe literal for the expression
func (e *ScopedFuncCallExpr) Literal() string {
	_, sok := e.scope.(*SymbolExpr)
	if sok {
		return "#ScopeFunc:" + e.scope.Literal() + "." + e.name
	}
	return e.scope.Literal() + "." + e.name
}

// Value will provide value after evaluation
func (e *ScopedFuncCallExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *ScopedFuncCallExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// NativeCallBack is the hook for native functions in the environment
// Used for extending the environment with new functions
type NativeCallBack func(env *Environment, args []Expression) (Expression, error)

// NativeFunctionExpr is a expression holding a native function and its arguments
type NativeFunctionExpr struct {
	name   string
	native NativeCallBack
}

// NewNativeFunctionExpr registers a native function in the form symbol.symbol(args)
func NewNativeFunctionExpr(name string, native NativeCallBack) *NativeFunctionExpr {
	e := NativeFunctionExpr{}
	e.name = name
	e.native = native
	return &e
}

// Invoke for implementation of function interface
func (e *NativeFunctionExpr) Invoke(env *Environment, args []Expression) (Expression, error) {
	return e.native(env, args)
}

// Evaluate the expression
func (e *NativeFunctionExpr) Evaluate(env *Environment) (Expression, error) {
	return e, nil
}

// Literal will provide a uniqe literal for the expression
func (e *NativeFunctionExpr) Literal() string {
	return fmt.Sprintf("(#native-function:%s#)", e.name)
}

// Value will provide value after evaluation
func (e *NativeFunctionExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *NativeFunctionExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//=============================================================================

// ScopedNativeFunctionExpr is a expression holding a scoped native! function and its arguments
// Fantastic stuff!
type ScopedNativeFunctionExpr struct {
	name   string
	scope  Expression
	native NativeCallBack
}

// NewScopedNativeFunctionExpr registers a native function in the form symbol.symbol(args)
func NewScopedNativeFunctionExpr(name string, scope Expression, native NativeCallBack) *ScopedNativeFunctionExpr {
	e := ScopedNativeFunctionExpr{}
	e.name = name
	e.scope = scope
	e.native = native
	return &e
}

// Invoke for implementation of function interface
func (e *ScopedNativeFunctionExpr) Invoke(env *Environment, args []Expression) (Expression, error) {
	return e.native(env, args)
}

// Evaluate the expression
func (e *ScopedNativeFunctionExpr) Evaluate(env *Environment) (Expression, error) {
	return e, nil
}

// Literal will provide a uniqe literal for the expression
func (e *ScopedNativeFunctionExpr) Literal() string {
	return fmt.Sprintf("(#native-function:%s#)", e.name)
}

// Value will provide value after evaluation
func (e *ScopedNativeFunctionExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

// String will provide the string representation of value
func (e *ScopedNativeFunctionExpr) String() string {
	return fmt.Sprintf("%T", e)
}
