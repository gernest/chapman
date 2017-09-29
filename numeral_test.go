package chapman

import (
	"strings"
	"testing"
)

func TestNumeralLexer(t *testing.T) {
	var l numeralLexer
	sample := []string{
		"3.14", "3", "123e5", "123e-5", " 0.5e2", " 0xFF",
		"0b11111111", "0o377",
	}
	for _, v := range sample {
		s := newBufioScanner(strings.NewReader(v))
		if !l.accept(s) {
			t.Error("expected to accept", v)
		}
		tk, err := l.lex(s, &context{})
		if err != nil {
			t.Errorf("case %s %v", v, err)
		}
		if tk != nil {
			if tk.Text != v {
				t.Errorf("expected %s got %s", v, tk.Text)
			}
		}

	}
}
