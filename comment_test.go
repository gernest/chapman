package goes

import (
	"strings"
	"testing"
)

func TestSingleLineComment(t *testing.T) {
	c, err := cases("fixture/comment/single")
	if err != nil {
		t.Fatal(err)
	}
	var l singleLineCommentLexer
	for _, v := range c {
		s := newBufioScanner(strings.NewReader(v.actual))
		if !l.accept(s) {
			t.Errorf("expected to accept %s", v.dir)
		}
		tk, err := l.lex(s, &context{})
		if err != nil {
			t.Fatal(err)
		}
		nv := string(printToken(tk))
		if nv != v.expected {
			t.Errorf("expected %s got %s", v.expected, nv)
		}
	}
}
