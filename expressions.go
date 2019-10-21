package expr

import (
	"fmt"
	"strings"
)

//Expression is the generic expression for all expression types in the expression framework
type Expression interface {
	//Evaluate the expression
	Evaluate(env *Environment) Expression
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
	Invoke(env Environment, args []Expression) Expression
}

//Scalar is a basic scalar expression
type Scalar struct {
	literal string
	value   interface{}
}

//NewScalar registers a new scalar with the defined literal and value
func NewScalar(literal string, value interface{}) *Scalar {
	sc := Scalar{}
	sc.value = value
	sc.literal = literal
	return &sc
}

//NewScalarV registers a new scalar and creates and literal based on the value
func NewScalarV(value interface{}) *Scalar {
	sc := Scalar{}
	sc.value = value
	if sc.value != nil {
		switch v := sc.value.(type) {
		case string:
			sc.literal = escape(sc.value.(string))
		default:
			sc.literal = fmt.Sprintf("%T", v)
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
func (e *Scalar) Evaluate(env *Environment) Expression {
	return e
}

//Literal will provide a uniqe literal for the expression
func (e *Scalar) Literal() string {
	return e.literal
}

//Value will provide value after evaluation
func (e *Scalar) Value() interface{} {
	return e.value
}

//String will provide the string representation of value
func (e *Scalar) String() string {
	if e.value != nil {
		return fmt.Sprintf("%v", e.value)
	}
	return "NULL"
}

//condExpr is a conditional expression on the form cond? left: right
type condExpr struct {
	condition Expression
	left      Expression
	right     Expression
}

//NewCondExpr registers a new conditional expression on the form cond?left:right
func newCondExpr(cond Expression, left Expression, right Expression) *condExpr {
	e := condExpr{}
	e.condition = cond
	e.left = left
	e.right = right
	return &e
}

//Evaluate the expression
func (e *condExpr) Evaluate(env *Environment) Expression {
	c := e.condition.Evaluate(env)
	if c == nil || c == env.False() || c.Value() == nil || c.Value() == env.False() {
		return e.left.Evaluate(env)
	}
	return e.right.Evaluate(env)
}

//Literal will provide a uniqe literal for the expression
func (e *condExpr) Literal() string {
	return fmt.Sprintf("(%s?%s:%s)", e.condition.Literal(), e.left.Literal(), e.right.Literal())
}

//Value will provide value after evaluation
func (e *condExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *condExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//orExpr is a basic binary expression(||)
type orExpr struct {
	left  Expression
	right Expression
}

//NewOrExpr registers a new conditional expression on the form left||right
func newOrExpr(left Expression, right Expression) *orExpr {
	e := orExpr{}
	e.left = left
	e.right = right
	return &e
}

//Evaluate the expression
func (e *orExpr) Evaluate(env *Environment) Expression {
	c := e.left.Evaluate(env)
	if c != nil && c != env.False() && c.Value() != nil && c.Value() != env.False() {
		return env.True()
	}
	c = e.right.Evaluate(env)
	if c != nil && c != env.False() && c.Value() != nil && c.Value() != env.False() {
		return env.True()
	}
	return env.False()
}

//Literal will provide a uniqe literal for the expression
func (e *orExpr) Literal() string {
	return fmt.Sprintf("(%s || %s)", e.left.Literal(), e.right.Literal())
}

//Value will provide value after evaluation
func (e *orExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *orExpr) String() string {
	return fmt.Sprintf("%T", e)
}

//andExpr is a basic binary expression(&&)
type andExpr struct {
	operand string
	left    Expression
	right   Expression
}

//NewOrExpr registers a new conditional expression on the form left||right
func newAndExpr(left Expression, right Expression) *andExpr {
	e := andExpr{}
	e.operand = "&&"
	e.left = left
	e.right = right
	return &e
}

//Evaluate the expression
func (e *andExpr) Evaluate(env *Environment) Expression {
	c := e.left.Evaluate(env)
	if c == nil || c == env.False() || c.Value() == nil || c.Value() == env.False() {
		return env.False()
	}
	c = e.right.Evaluate(env)
	if c == nil || c == env.False() || c.Value() == nil || c.Value() == env.False() {
		return env.False()
	}
	return env.True()
}

//Literal will provide a uniqe literal for the expression
func (e *andExpr) Literal() string {
	return fmt.Sprintf("(%s %s %s)", e.left.Literal(), e.operand, e.right.Literal())
}

//Value will provide value after evaluation
func (e *andExpr) Value() interface{} {
	return fmt.Sprintf("[:%T:]", e)
}

//String will provide the string representation of value
func (e *andExpr) String() string {
	return fmt.Sprintf("%T", e)
}
