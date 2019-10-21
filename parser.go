package expr

import (
	"fmt"
	"strings"
)

/* Expression BNF
* ==============
* expr	::=	condExpr
* condExpr::=	orExpr['?' expr ':' expr]
* orExpr	::= 	andExpr['||' orExpr]
* andExpr ::=	cmpExpr['&&' andExpr]
* cmpExpr ::=     catExpr['==', catExpr]
* catExpr ::=     atom ('+' CatExpr)*
* atom    ::=     (text | symbol | subexpr)
* symbol::=	ident[funcall]['.' symbol]
* funcall ::=	'(' [arglist] ')'
* arglist::= expr(',' expr) *
* subexpr::= '(' expr ')'
* listexpr::= '{' arglist '}'
 */

//Parser parses
type Parser struct {
}

func newParser() Parser {
	return *new(Parser)
}

//Parse creates an parser witch parser the input string
//Remember to evaluate the returned expression in the prefered environment!
//Enjoy!
func Parse(input string) (Expression, error) {
	p := newParser()
	return p.Parse(input)
}

//Parse parses the expressing from a string format
func (p Parser) Parse(input string) (Expression, error) {
	l := lex(false, input)
	return parseExpr(*l)
}

//parseExpr parses the lexer tokens
func parseExpr(lex lexer) (Expression, error) {
	return parseCond(lex)
}

func parseCond(lex lexer) (Expression, error) {

	expr, err := parseOr(lex)
	if err != nil {
		return expr, err
	}

	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return expr, nil
	}
	if t.Type == OperatorTok && t.Literal == "?" {

		cond := condExpr{}
		cond.condition = expr
		cond.left, err = parseExpr(lex)
		if err != nil {
			return expr, err
		}
		if _, err := expect(lex, OperatorTok, ":"); err == nil {
			return &cond, err
		}
		cond.right, err = parseExpr(lex)
		expr = &cond

	} else {

		lex.PushBack(t)
	}
	return expr, nil
}

//parseOr. || Operator.
func parseOr(lex lexer) (Expression, error) {
	expr, err := parseAnd(lex)
	if err != nil {
		return expr, err
	}
	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return expr, nil
	}
	if t.Type == OperatorTok && t.Literal == "||" {
		//left has precedens
		l, err := parseOr(lex)
		if err != nil {
			return l, err
		}
		return newOrExpr(expr, l), nil
	}
	lex.PushBack(t)
	return expr, nil
}

//parseand. && Operator.
func parseAnd(lex lexer) (Expression, error) {
	expr, err := parseCmp(lex)
	if err != nil {
		return expr, err
	}
	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return expr, nil
	}
	if t.Type == OperatorTok && t.Literal == "&&" {
		//left has precedens
		l, err := parseAnd(lex)
		if err != nil {
			return l, err
		}
		return newOrExpr(expr, l), nil
	}
	lex.PushBack(t)
	return expr, nil
}

//To be continued
func parseCmp(lex lexer) (Expression, error) {
	return &condExpr{}, nil
}

func expect(lex lexer, tt TokenType, literal string) (Token, error) {

	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return t, nil
	}
	if t.Type != tt {
		//TODO check this!! %d not correct? String representation?
		return t, fmt.Errorf("Unexpected TokenType: %d", t.Type)
	}
	if strings.EqualFold(t.Literal, literal) {
		return t, fmt.Errorf("Unexpected literal value: %s (expected %s )", t.Literal, literal)
	}
	return t, nil
}
