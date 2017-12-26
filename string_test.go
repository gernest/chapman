package chapman

import (
	"strings"
	"testing"
)

func TestStringLexer(t *testing.T) {
	badStrings := []string{
		`"\u{g0g}"'`,
		`"\u{g}"`,
		`"\u{g0}"`,
		`"\u{g0}"`,
		`\u{0g}`,
		`\u{0g0}\r\n`,
		`"\u{110000}"`,
		`"\u{11ffff}"`,
		`"\x0g"`,
		`"\xg0\r\n"`,
		`"\xgg"`,
		`"\u1"`,
	}

	var l stringLexer
	for _, v := range badStrings {
		s := newBufioScanner(strings.NewReader(v))
		_, err := l.lex(s, &context{})
		if err == nil {
			t.Errorf("expected an error for %v", v)
		}
	}

	goodStrings := []string{
		`"abc"`,
		`"\\б"`,
		`"\\u0435"`,
		`"Hello\\nworld"`,
		`"\\u0435"`,
		`"\\u0432"`,
		`"\\u180E"`,
		`"\\7"`,
		`"Hello\\012World"`,
		`"Hello\\412World"`,
		`"Hello\\712World"`,
		`"Hello\\1World"`,
		`"\\xff"`,
		`"\\u{11000}"'`,
		`"\\Щ"`,
		`"\\З"`,
		`"\\ю"`,
		`"\\б"`,
		`"a\\r\\nb"`,
		`"\\u0451"`,
		`"\\u0006A"`,
	}
	for _, v := range goodStrings {
		s := newBufioScanner(strings.NewReader(v))
		_, err := l.lex(s, &context{})
		if err != nil {
			t.Error(err)
		}
	}
}
