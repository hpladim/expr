package expr

import (
	"errors"
	"fmt"
	"strings"
)

/* Expression BNF
* ==============
* expr			::=		condExpr
* condExpr		::=		orExpr['?' expr ':' expr]
* orExpr		::= 	andExpr['||' orExpr]
* andExpr		::=		cmpExpr['&&' andExpr]
* cmpExpr		::=		concatExpr['==', concatExpr]
* concatExpr	::=		atom ('+' concatExpr)*
* atom			::=		(text | symbol | subexpr)
* symbol		::=		ident[funcall]['.' symbol]
* funcall		::=		'(' [arglist] ')'
* arglist		::=		expr(',' expr) *
* subexpr		::=		'(' expr ')'
* listexpr		::=		'{' arglist '}'
 */

//Parser parses
type Parser struct {
}

func newParser() Parser {
	return *new(Parser)
}

func compareOp(operand string) bool {

	switch operand {
	case "==", "!=", ">=", "<=":
		return true
	default:
		return false
	}
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
	l := newLexer(false, input)
	return parseExpr(l)
}

//parseExpr parses the lexer tokens
func parseExpr(lex *lexer) (Expression, error) {
	return parseCond(lex)
}

func parseCond(lex *lexer) (Expression, error) {

	expr, err := parseOr(lex)
	if err != nil {
		return expr, err
	}

	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return expr, nil
	}
	if t.Type == OperatorTok && t.Literal == "?" {

		cond := CondExpr{}
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
func parseOr(lex *lexer) (Expression, error) {
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
		r, err := parseOr(lex)
		if err != nil {
			return r, err
		}
		return NewOrExpr(expr, r), nil
	}
	lex.PushBack(t)
	return expr, nil
}

//parseand. && Operator.
func parseAnd(lex *lexer) (Expression, error) {
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
		r, err := parseAnd(lex)
		if err != nil {
			return r, err
		}
		return NewAndExpr(expr, r), nil
	}
	lex.PushBack(t)
	return expr, nil
}

//parseCmp parses the compare expression
func parseCmp(lex *lexer) (Expression, error) {
	left, err := parseConcat(lex)
	if err != nil {
		return left, err
	}
	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return left, nil
	}
	if t.Type == OperatorTok {
		if compareOp(t.Literal) {
			right, err := parseConcat(lex)
			if err != nil {
				return right, err
			}
			return NewCompareExpr(t.Literal, left, right), nil
		}
		lex.PushBack(t)

	} else if t.Type == IdentTok && t.Literal == "like" {
		right, err := parseConcat(lex)
		if err != nil {
			return left, err
		}
		return NewLikeExpr(left, right), nil

	} else if t.Type == IdentTok && t.Literal == "in" {
		right, err := parseConcat(lex)
		if err != nil {
			return left, err
		}
		return NewInExpr(left, right), nil
	} else {
		lex.PushBack(t)
	}
	return left, nil
}

func parseConcat(lex *lexer) (Expression, error) {
	left, err := parseAtom(lex)
	if err != nil {
		return left, err
	}
	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return left, nil
	}
	if t.Type == OperatorTok && t.Literal == "+" {

		right, err := parseConcat(lex)
		if err != nil {
			return right, err
		}
		return NewConcatExpr(left, right), nil

	}

	lex.PushBack(t)
	return left, nil
}

func parseAtom(lex *lexer) (Expression, error) {

	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return NewScalarExpr("", nil), nil
	}
	switch t.Type {
	case NumberTok, StringTok:
		return NewScalarExpr(t.Literal, t.Value), nil
	case IdentTok:
		return parseIdent(lex, t)
	case OperatorTok:
		switch t.Literal {
		case "(":
			sub, err := parseSubExpr(lex)
			if err != nil {
				return sub, err
			}
			_, err = expect(lex, OperatorTok, ")")
			if err != nil {
				return sub, err
			}
			return sub, nil
		case "{":
			return parseListExpr(lex)
		//case "[":
		//    expr = parseStruct(lex);
		//    break;
		default:
			return nil, fmt.Errorf("Unexpected operator when expecting expression term: %s", t.Literal)
		}
	default:
		return nil, fmt.Errorf("Unexpected token type when expecting expression term: %v '%s'", t.Type, t.Literal)
	}
}

//parseSubExpr parses list expression in format (expr)
func parseSubExpr(lex *lexer) (Expression, error) {
	return parseExpr(lex)
}

//parseListExpr parses list expression in format {expr,expr,....}
func parseListExpr(lex *lexer) (Expression, error) {
	result := NewListExpr()
	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return result, nil
	}
	for t.Type != EoFTok || fini {
		if t.Type == OperatorTok && t.Literal == "}" {
			return result, nil
		}

		lex.PushBack(t)
		lele, err := parseExpr(lex)
		if err != nil {
			return lele, err
		}
		result.Append(lele)
		t, fini = lex.NextToken()
		if !fini && t.Type == OperatorTok && t.Literal == "," {
			t, fini = lex.NextToken()
		}
	}
	return result, errors.New("List unterminated")
}

func parseScopedIdent(lex *lexer, t Token, scope *SymbolExpr) (Expression, error) {

	ident := t.Literal
	sym := NewSymbolExprWithScope(t.Literal, scope)
	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return sym, nil
	}
	if t.Type == OperatorTok && t.Literal == "." {

		t, err := expect(lex, IdentTok, "")
		if err != nil {

			return sym, err
		}
		return parseScopedIdent(lex, t, sym)
	}
	if t.Type != OperatorTok || t.Literal != "(" {
		lex.PushBack(t)
		return sym, nil
	}
	f, err := NewScopedFuncCallExpr(ident, scope)
	if err != nil {
		return f, err
	}
	expectArg := true

	t, fini = lex.NextToken()
	if fini || t.Type == EoFTok {
		return f, nil
	}
	for !(t.Type == OperatorTok && t.Literal == ")") {
		lex.PushBack(t)
		if expectArg {
			a, err := parseExpr(lex)
			if err != nil {
				return a, err
			}
			f.AddArg(a)
		} else {
			expect(lex, OperatorTok, ",")

		}
		t, fini = lex.NextToken()
		if fini || t.Type == EoFTok {
			return sym, nil
		}

		expectArg = !expectArg

	}
	t, fini = lex.NextToken()
	if fini || t.Type == EoFTok {
		return sym, nil
	}
	if t.Type == OperatorTok && t.Literal == "." {

		t, err := expect(lex, IdentTok, "")
		if err != nil {

			return sym, err
		}
		return parseScopedIdent(lex, t, sym)
	}
	return f, nil
}

func parseIdent(lex *lexer, t Token) (Expression, error) {

	sym := NewSymbolExpr(t.Literal)
	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return sym, nil
	}
	if t.Type == OperatorTok && t.Literal == "." {

		t, err := expect(lex, IdentTok, "")
		if err != nil {
			return sym, err
		}
		return parseScopedIdent(lex, t, sym)
	}
	if t.Type != OperatorTok || t.Literal != "(" {
		lex.PushBack(t)
		return sym, nil
	}
	f := NewFuncCallExpr()
	expectArg := true
	f.SetFunc(sym)

	t, fini = lex.NextToken()
	if fini || t.Type == EoFTok {
		return f, nil
	}
	for !(t.Type == OperatorTok && t.Literal == ")") {
		lex.PushBack(t)
		if expectArg {
			a, err := parseExpr(lex)
			if err != nil {
				return a, err
			}
			f.AddArg(a)
		} else {
			expect(lex, OperatorTok, ",")

		}
		t, fini = lex.NextToken()
		if fini || t.Type == EoFTok {
			return sym, nil
		}
		expectArg = !expectArg

	}
	return f, nil
}

func expect(lex *lexer, tt TokenType, literal string) (Token, error) {

	t, fini := lex.NextToken()
	if fini || t.Type == EoFTok {
		return t, nil
	}
	if t.Type != tt {
		//TODO check this!! %v not correct? String representation?
		return t, fmt.Errorf("Unexpected TokenType: %v", t.Type)
	}
	if !strings.EqualFold(t.Literal, literal) {
		return t, fmt.Errorf("Unexpected literal value: %s, (expected %s)", t.Literal, literal)
	}
	return t, nil
}
