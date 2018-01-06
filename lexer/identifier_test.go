package lexer

import (
	"strings"
	"testing"
)

func TestIdentifierNameLexer(t *testing.T) {
	var l identifierNameLexer
	for _, v := range keywords {
		s := newBufioScanner(strings.NewReader(v))
		if !l.Accept(s) {
			t.Error("expected to accept")
		}
		tk, err := l.Lex(s, &context{})
		if err != nil {
			t.Fatal(err)
		}
		if tk.Text != v {
			t.Errorf("expected %s got %s", v, tk.Text)
		}
	}

	escapes := []string{
		"\u006C\u006F\u006C\u0077\u0061\u0074",
		`\u{000000000061}`,
	}
	for _, v := range escapes {
		s := newBufioScanner(strings.NewReader(v))
		if !l.Accept(s) {
			t.Error("expected to accept", v)
		}
		tk, err := l.Lex(s, &context{})
		if err != nil {
			t.Fatal(err)
		}
		if tk.Text != v {
			t.Errorf("expected %s got %s", v, tk.Text)
		}
	}
}
