package chapman

import (
	"strings"
	"testing"
)

const expectTokenDecode = `{
	"Text": " single line comment",
	"Kind": "SINGLE_LINE_COMMENT",
	"Start": {
		"Line": 0,
		"Column": 0
	},
	"End": {
		"Line": 0,
		"Column": 22
	}
}`

func TestTokenUnmarshalJSON(t *testing.T) {
	tk, err := decodeToken([]byte(expectTokenDecode))
	if err != nil {
		t.Fatal(err)
	}
	e := string(printToken(tk))
	if e != expectTokenDecode {
		t.Errorf("expected %s got %s", expectTokenDecode, e)
	}
}

func TestBufioScanner_peekAt(t *testing.T) {
	e := []rune{'h', 'e', 'l', 'l', 'o'}
	a := "hello"
	b := newBufioScanner(strings.NewReader(a))

	for i := 1; i <= len(a); i++ {
		x, _, err := b.peekAt(i)
		if err != nil {
			t.Fatal(err)
		}
		ch := e[i-1]
		if ch != x {
			t.Errorf("expected %d: %s got %s", i, string(ch), string(x))
		}
	}
}
