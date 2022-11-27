package expr

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const eof = -1

// stateFunc is defined to handle implementation of the different tokentypes.
// See token lexing below
type stateFunc func(*lexer) stateFunc

// lexer functions similarly to Rob Pike's discussion
// about lexer design in this [talk](https://www.youtube.com/watch?v=HxaD_trXwRE).
type lexer struct {
	input      string
	start      int
	pos        int
	width      int
	line       int
	includeWS  bool
	doubleOps  map[string]bool
	tokens     chan Token
	tokenStack tokenStack
}

// newLexer creates and returns a lexer ready to parse the given input
func newLexer(includeWS bool, in string) *lexer {
	l := lexer{
		input:      in,
		start:      0,
		pos:        0,
		width:      0,
		line:       0,
		includeWS:  includeWS,
		doubleOps:  registerDoubleOps(),
		tokens:     make(chan Token),
		tokenStack: newTokenStack(),
	}
	go l.run()
	return &l
}

func registerDoubleOps() map[string]bool {
	ops := make(map[string]bool)
	ops["=="] = true
	ops["!="] = true
	ops["<="] = true
	ops[">="] = true
	ops["||"] = true
	ops["&&"] = true
	return ops
}

// next pulls the next rune from the lexer and returns it, moving the position
// forward in the source.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func (l lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// Current returns the value being analyzed at this moment.
func (l lexer) Current() string {
	return l.input[l.start:l.pos]
}

// Emit will receive a token and push a new token with the current analyzed
// value into the tokens channel.
func (l *lexer) emit(tok Token, send bool) {

	if send {
		l.tokens <- tok
	}
	l.start = l.pos
}

// NextToken returns the next token from the lexer and a value to denote whether
// or not the lexer is finished.
func (l *lexer) NextToken() (token Token, finished bool) {
	if st := l.tokenStack.pop(); st != nil {
		return *st, false
	}
	if tok, ok := <-l.tokens; ok {
		return tok, false
	}
	return Token{}, true

}

// PushBack returns the next token from the lexer and a value to denote whether
// or not the lexer is finished.
func (l *lexer) PushBack(t Token) {
	l.tokenStack.push(t)
}

// Private methods
func (l *lexer) run() {
	for state := lexInput; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// eat receives a string containing all acceptable strings and will contine
// over each consecutive character in the source until a token not in the given
// string is encountered. This should be used to quickly pull token parts.
func (l *lexer) eat(chars string) string {
	result := ""
	r := l.next()
	for strings.ContainsRune(chars, r) {
		result += string(r)
		r = l.next()
	}
	l.backup() // last next wasn't a match
	return result
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

//====================================================================================
//Token lexing below
//====================================================================================

const chWHITE = " \t\r\n\u00A0"
const chDIGIT = "0123456789"
const chBinary = "01"
const chOCTAL = "01234567"
const chHEXDIGIT = chDIGIT + "abcdefABCDEF"
const chALPHA = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const chIDENTSTART = chALPHA + "_"
const chIDENTBODY = chIDENTSTART + chDIGIT

// lexInput starts the lexing of the text input string
func lexInput(l *lexer) stateFunc {
	l.width = 0
	r := l.next()
	if strings.ContainsRune(chWHITE, r) {
		return lexSpace(l)
	} else if strings.ContainsRune(chIDENTSTART, r) {
		return lexIdent(l)
	} else if strings.ContainsRune(chDIGIT, r) {
		l.backup()
		return lexNumber(l)
	} else if r == '"' || r == '\'' {
		l.backup()
		return lexQuoted(l)
	} else if r == eof {
		l.emit(Token{Type: EoFTok}, true)
		return nil
	}
	l.backup()
	return lexOperator(l)
}

// lexSpace scans for space characters.
func lexSpace(l *lexer) stateFunc {
	l.eat(chWHITE)
	tok := Token{
		Type:    WhitespaceTok,
		Literal: l.Current(),
		Line:    l.line,
		Start:   l.start,
	}
	l.emit(tok, l.includeWS)
	return lexInput(l)
}

// lexIdent scans all ident body characters.
func lexIdent(l *lexer) stateFunc {
	l.eat(chIDENTBODY)
	tok := Token{
		Type:    IdentTok,
		Literal: l.Current(),
		Line:    l.line,
		Start:   l.start,
	}
	l.emit(tok, true)
	return lexInput(l)
}

// lexIdent scans all ident body characters.
func lexNumber(l *lexer) stateFunc {
	fl, err := l.scanNumber()
	if err != nil {
		return l.errorf("bad number syntax: %q. error: %q", l.Current(), err.Error())
	}
	tok := Token{
		Type:    NumberTok,
		Literal: l.Current(),
		Line:    l.line,
		Start:   l.start,
		Value:   fl,
	}
	l.emit(tok, true)
	return lexInput(l)
}
func (l *lexer) scanNumber() (interface{}, error) {
	// Optional leading sign.
	//l.accept("+-")
	digits := chDIGIT
	// Is it hex,octal,binary?
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = chHEXDIGIT
			l.eat(digits)
			return strconv.ParseInt(l.Current(), 16, 64)

		} else if l.accept("oO") {
			digits = chOCTAL
			l.eat(digits)
			return strconv.ParseInt(l.Current(), 8, 64)
		} else if l.accept("bB") {
			digits = chBinary
			l.eat(digits)
			return strconv.ParseInt(l.Current(), 2, 64)
		}
	}
	if l.accept(".") {
		l.eat(digits)
		return strconv.ParseFloat(l.Current(), 64)
	}
	return strconv.ParseInt(l.Current(), 10, 64)
}

// lexQuoted scans a quoted string.
func lexQuoted(l *lexer) stateFunc {
	quote := l.next()
	tok := Token{
		Type: StringTok,
	}
	value := ""
Loop:
	for {
		c := l.next()
		switch c {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case quote:
			break Loop
		}
		value += string(c)
	}
	tok.Literal = l.Current()
	tok.Value = value
	tok.Line = l.line
	tok.Start = l.start
	l.emit(tok, true)
	return lexInput(l)
}

// lexOperator lexes both single and dounle operands
func lexOperator(l *lexer) stateFunc {
	tok := Token{
		Type:  OperatorTok,
		Line:  l.line,
		Start: l.start,
	}
	c := l.next()
	c2 := l.next()
	tok.Literal = string(c) + string(c2)
	if !l.doubleOp(tok.Literal) {
		l.backup()
		tok.Literal = string(c)
	}
	l.emit(tok, true)
	return lexInput(l)
}

func (l *lexer) doubleOp(literal string) bool {
	_, ok := l.doubleOps[literal]
	return ok
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	tok := Token{
		Type:    ErrorTok,
		Literal: "ERROR",
		Value:   fmt.Sprintf(format, args...),
		Line:    l.line,
		Start:   l.start,
	}
	l.emit(tok, true)
	return nil
}

//====================================================================================
//Token below
//====================================================================================

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Value   interface{}
	Line    int
	Start   int
}

// TokenType is the different tokentypes provided by
type TokenType int

// Tokentypes
const (
	ErrorTok TokenType = iota
	EoFTok
	WhitespaceTok
	IdentTok
	StringTok
	NumberTok
	OperatorTok
)

type tokenNode struct {
	t    Token
	next *tokenNode
}

type tokenStack struct {
	start *tokenNode
}

func newTokenStack() tokenStack {
	return tokenStack{}
}

func (s *tokenStack) push(t Token) {
	node := &tokenNode{t: t}
	if s.start == nil {
		s.start = node
	} else {
		node.next = s.start
		s.start = node
	}
}

func (s *tokenStack) pop() *Token {
	if s.start == nil {
		return nil
	}
	n := s.start
	s.start = n.next
	return &n.t
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (s *tokenStack) clear() {
	s.start = nil
}
