package lexer

import (
	"strings"
	"testing"
)

func TestPunctuatorLexer(t *testing.T) {
	var l punctuatorLexer
	for v := range puncs {
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
