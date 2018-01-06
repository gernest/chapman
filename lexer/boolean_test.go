package lexer

import (
	"strings"
	"testing"
)

func TestBooleanLexer(t *testing.T) {
	var l boolLexer
	sample := []string{"true", "false"}
	for _, v := range sample {
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
