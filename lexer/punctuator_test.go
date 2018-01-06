package lexer

import (
	"strings"
	"testing"
)

func TestPunctuationLexer(t *testing.T) {
	var l punctuationLexer
	for v := range punctuation {
		s := newBufioScanner(strings.NewReader(v))
		if !l.Accept(s) {
			t.Error("expected to accept", v)
		}
		tk, err := l.Lex(s, &context{})
		if err != nil {
			t.Fatal(err)
		}
		if tk.Text != v {
			t.Errorf("expected  %s got %s", v, tk.Text)
		}
		k := punctuationKind[v]
		if k != tk.Kind {
			t.Errorf("%s :expected %s got %s", v, k, tk.Kind)
		}
	}
}
