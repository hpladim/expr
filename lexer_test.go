package expr

import (
	"testing"
)

func Test(t *testing.T) {
	t.Run("QuotedString", func(t *testing.T) { RunLexerQuotedStringTest(t) })
}

func RunLexerQuotedStringTest(t *testing.T) {

	t.Run("Quote#1", func(t *testing.T) {
		l := newLexer(false, "\"test\"")
		tok, fini := l.NextToken()
		if fini {
			t.Errorf("\n\nLexer failed,NextToken() failed: %v. \n", t)
		}
		if tok.Type != StringTok {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}
		if tok.Literal != "\"test\"" && tok.Value != "test" {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token literal and value:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}

	})
	t.Run("ListQuote#1", func(t *testing.T) {
		l := newLexer(false, "{\"\"}")
		tok, fini := l.NextToken()
		if fini {
			t.Errorf("\n\nLexer failed,NextToken() failed: %v. \n", t)
		}
		if tok.Type != OperatorTok {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}
		if tok.Literal != "{" {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token literal:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}
		tok, fini = l.NextToken()
		if fini {
			t.Errorf("\n\nLexer failed,NextToken() failed: %v. \n", t)
		}
		if tok.Type != StringTok {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}
		if tok.Literal != "\"\"" && tok.Value != "" {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token literal and value:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}
		tok, fini = l.NextToken()
		if fini {
			t.Errorf("\n\nLexer failed, NextToken() failed: %v. \n", t)
		}
		if tok.Type != OperatorTok {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}
		if tok.Literal != "}" {
			t.Errorf("\n\nLexer failed, Nextoken returned unexpected token value:\n\n token: %v \n\n lexer: %v\n", tok, l)
		}

	})

}
