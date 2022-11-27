package expr

import (
	"errors"
	"fmt"
	"sync"
)

// Environment is is used for registering expressions and holding the parser.
// Environment is set up with a basic set of expressions by calling RegisterBuiltins()
// Environment are further extendable by RegisterSymbol(), RegisterFunction() and RegisterScopedFunction()
type Environment struct {
	exStack             exStack
	parser              Parser
	globalFuncs         map[string]eValue
	AutoregisterGlobals bool
}

// eValue is used for keeping global registered functions
type eValue struct {
	expr     Expression
	readOnly bool
}

// NewEnvironment returns a new Environment ready set up with a minimum of bultin expressions
// Environment is is used for registering expressions and holding the parser.
// Environment is set up with a basic set of expressions by calling RegisterBuiltins()
// Environment are further extendable by RegisterSymbol(), RegisterFunction() and RegisterScopedFunction()
// The expressions currently registered are:
// - 'null':  expression that evaluates to the symbol 'null'
// - 'true': expression that evaluates to the boolean 'true'
// - 'false': expression that evaluates to the boolean 'false'
// - 'empty': expression that evaluates to an empty string
// - 'print' expression in form print(args ... Expression) prints all arguments to console
func NewEnvironment() *Environment {
	e := new(Environment)
	e.parser = newParser()
	e.globalFuncs = make(map[string]eValue)
	e.registerBuiltIns()
	return e
}

// RegisterBuiltIns will set up the environment with a minimum set of expressions
// The expressions currently registered are:
// - 'null':  expression that evaluates to the symbol 'null'
// - 'true': expression that evaluates to the boolean 'true'
// - 'false': expression that evaluates to the boolean 'false'
// - 'empty': expression that evaluates to an empty string
// - 'print' expression in form print(args ... Expression) prints all arguments to console
func (e *Environment) registerBuiltIns() error {

	e.Set("null", NewScalarExprV(nil))
	e.Set("true", NewScalarExprV(true))
	e.Set("false", NewScalarExprV(false))
	e.Set("empty", NewScalarExpr("empty", ""))
	e.Lock("null", true)
	e.Lock("false", true)
	e.Lock("true", true)

	//Native functions
	e.RegisterFunction("print", Print)

	return nil
}

// True returns a True expression from the environment
func (e *Environment) True() Expression {
	return e.Get("true")
}

// False returns a False expression from the environment
func (e *Environment) False() Expression {
	return e.Get("false")
}

// Null returns a Null expression from the environment
func (e *Environment) Null() Expression {
	return e.Get("null")
}

// Empty returns a Null expression from the environment
func (e *Environment) Empty() Expression {
	return e.Get("empty")
}

// RegisterFunction registers a native function in the Environment
// Use this to Extend the Environment
func (e *Environment) RegisterFunction(name string, callback NativeCallBack) {

	e.Set(name, NewNativeFunctionExpr(name, callback))
}

// RegisterScopedFunction registers a native function in the Environment
// Use this to Extend the Environment
// Firs argument in function is the scope
func (e *Environment) RegisterScopedFunction(name string, scope Expression, callback NativeCallBack) Expression {
	scopeFunc := NewScopedNativeFunctionExpr(name, scope, callback)
	e.Set(name, scopeFunc)
	return scopeFunc
}

// RegisterSymbol is used for registering symbols in the Environment
func (e *Environment) RegisterSymbol(symbol SymbolExpr, expr Expression, immutable bool) error {
	name := symbol.Literal()
	if _, ok := e.globalFuncs[name]; ok {
		return fmt.Errorf("symbol already defined: " + name)
	}
	val := eValue{readOnly: immutable}
	val.expr = expr
	e.globalFuncs[name] = val
	return nil
}

// Set registers a new expression in the environment
func (e *Environment) Set(name string, expr Expression) error {
	var val eValue
	var ok bool
	if val, ok = e.globalFuncs[name]; !ok {
		val = eValue{readOnly: false}
	}
	if val.readOnly {
		return fmt.Errorf("symbol %s cannot be modified", name)
	}
	val.expr = expr
	e.globalFuncs[name] = val
	return nil
}

// Get a registered expression in the environment
func (e *Environment) Get(name string) Expression {
	var val eValue
	var ok bool
	if val, ok = e.globalFuncs[name]; !ok {
		if !e.AutoregisterGlobals {
			return e.Null()
		}
		val = eValue{expr: e.Null(), readOnly: false}
		e.globalFuncs[name] = val
	}
	return val.expr
}

// Lock a registered expression in the environment. If the expression does not exist, a null expression will be registered
func (e *Environment) Lock(name string, locked bool) {
	var val eValue
	var ok bool
	if val, ok = e.globalFuncs[name]; !ok {
		val = eValue{expr: e.Null()}
		e.globalFuncs[name] = val
	}
	val.readOnly = locked
}

// Remove removes a registered expression in the environment
func (e *Environment) Remove(name string) {
	delete(e.globalFuncs, name)
}

// Pushes function information to a stack. Used when functions are evaluated
func (e *Environment) pushStack(frm exFrame) {
	e.exStack.push(frm)
}

// Pops the previous function information from the stack. Used after current function is evaluated
func (e *Environment) popStack() {
	e.exStack.pop()
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (e *Environment) peek() (exFrame, error) {
	return e.exStack.peek()
}

// GetParser gets return the parser used by the Envionment
func (e *Environment) GetParser() Parser {
	return e.parser
}

// StackFrame holds the function information for the current context
type exFrame struct {
	function Function
	args     []Expression
}

type exStack struct {
	lock sync.Mutex
	frm  []exFrame
}

func (s *exStack) push(ef exFrame) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.frm = append(s.frm, ef)
}

func (s *exStack) pop() (exFrame, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.frm)
	if l == 0 {
		return exFrame{}, errors.New("Empty Stack")
	}

	res := s.frm[l-1]
	s.frm = s.frm[:l-1]
	return res, nil
}

func (s *exStack) peek() (exFrame, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.frm)
	if l == 0 {
		return exFrame{}, errors.New("Empty Stack")
	}

	res := s.frm[l-1]
	return res, nil
}
