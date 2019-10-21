package expr

import (
	"errors"
	"fmt"
	"sync"
)

/* //EvaluateEnv facilitates the needed methods for the evaluation of expressions
type EvaluateEnv interface {

	//Registers a new expression in the environment
	set(name string, expr Expression)
	//Get a registered epression in the environment
	get(name string) Expression
	//True returns a True expression from the environment
	true() Expression
	//False returns a False expression from the environment
	false() Expression
	//Null returns a Null expression from the environment
	null() Expression
	//Pushes function information to a stack. Used when functions are evaluated
	pushStack(frm exFrame)
	//Pops the previous function information from the stack. Used after current function is evaluated
	popStack()
} */

type eValue struct {
	expr     Expression
	readOnly bool
}

//Environment is is used for registering expressions and holding the parser.
//Environment are normally set up with a basic set of expressions by calling RegisterBuiltins()
type Environment struct {
	exStack             exStack
	parser              Parser
	asignLock           sync.Mutex
	globalFuncs         map[string]eValue
	AutoregisterGlobals bool
}

//NewEnvironment returns a new Environment ready set up.
//Remember to call RegisterBuiltins() to set up a minimum of expressions
func NewEnvironment() *Environment {
	e := new(Environment)
	e.parser = newParser()
	return e
}

//RegisterBuiltIns will set up the environment with a minimum set of expressions
//The expressions currently registered are:
//- 'true': expression that evaluates to the boolean 'true'
//- 'false': expression that evaluates to the boolean 'false'
//- 'null':  expression that evaluates to the symbol 'null'
//- 'empty': expression that evaluates to an empty string
func (e *Environment) RegisterBuiltIns() error {

	e.Set("null", NewScalarV(nil))
	e.Set("true", NewScalarV(true))
	e.Set("false", NewScalarV(true))
	e.Set("empty", NewScalar("empty", ""))
	e.Lock("null", true)
	e.Lock("false", true)
	e.Lock("true", true)
	return nil
}

//Set registers a new expression in the environment
func (e *Environment) Set(name string, expr Expression) error {
	e.asignLock.Lock()
	defer e.asignLock.Unlock()
	var val eValue
	if val, ok := e.globalFuncs[name]; !ok {
		val = eValue{readOnly: false}
		e.globalFuncs[name] = val
	}
	if val.readOnly {
		return fmt.Errorf("Symbol %s cannot be modified", name)
	}
	val.expr = expr
	return nil
}

//Get a registered expression in the environment
func (e *Environment) Get(name string) Expression {
	e.asignLock.Lock()
	defer e.asignLock.Unlock()
	var val eValue
	var ok bool
	if val, ok = e.globalFuncs[name]; !ok {
		if !e.AutoregisterGlobals {
			return e.Null()
		}
		val := eValue{expr: e.Null(), readOnly: false}
		e.globalFuncs[name] = val
	}
	return val.expr
}

//Lock a registered expression in the environment. If the expression does not exist, a null expression will be registered
func (e *Environment) Lock(name string, locked bool) {
	e.asignLock.Lock()
	defer e.asignLock.Unlock()
	var val eValue
	var ok bool
	if val, ok = e.globalFuncs[name]; !ok {
		val := eValue{expr: e.Null()}
		e.globalFuncs[name] = val
	}
	val.readOnly = locked
}

//True returns a True expression from the environment
func (e *Environment) True() Expression {
	return e.Get("true")
}

//False returns a False expression from the environment
func (e *Environment) False() Expression {
	return e.Get("false")
}

//Null returns a Null expression from the environment
func (e *Environment) Null() Expression {
	return e.Get("null")
}

//Empty returns a Null expression from the environment
func (e *Environment) Empty() Expression {
	return e.Get("empty")
}

//Pushes function information to a stack. Used when functions are evaluated
func (e *Environment) pushStack(frm exFrame) {
	e.exStack.push(frm)
}

//Pops the previous function information from the stack. Used after current function is evaluated
func (e *Environment) popStack() {
	e.exStack.pop()
}

func (e *Environment) peek() (exFrame, error) {
	return e.exStack.peek()
}

//GetParser gets return the parser used by the Envionment
func (e *Environment) GetParser() Parser {
	return e.parser
}

//StackFrame holds the function information for the current context
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
