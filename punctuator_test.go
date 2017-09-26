package goes

import (
	"strings"
	"testing"
)

func TestPunctuatorLexer(t *testing.T) {
	var l punctuatorLexer
	for v := range puncs {
		s := newBufioScanner(strings.NewReader(v))
		if !l.accept(s) {
			t.Error("expected to accept", v)
		}
		tk, err := l.lex(s, &context{})
		if err != nil {
			t.Fatal(err)
		}
		if tk.Text != v {
			t.Errorf("expected %s got %s", v, tk.Text)
		}
	}
}
