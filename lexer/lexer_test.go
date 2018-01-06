package lexer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		x, _, err := b.PeekAt(i)
		if err != nil {
			t.Fatal(err)
		}
		ch := e[i-1]
		if ch != x {
			t.Errorf("expected %d: %s got %s", i, string(ch), string(x))
		}
	}
}

func TestLexer(t *testing.T) {
	var files []string
	ferr := filepath.Walk("fixture", func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if i.IsDir() {
			return nil
		}
		if filepath.Base(p) == "actual.js" {
			files = append(files, p)
		}
		return nil
	})
	if ferr != nil {
		t.Fatal(ferr)
	}
	for _, f := range files {
		t.Run(filepath.Dir(f), func(ts *testing.T) {
			b, err := ioutil.ReadFile(f)
			if err != nil {
				ts.Fatal(err)
			}
			_, err = lex(bytes.NewReader(b), defaultLexMe()...)
			if err != nil {
				abs, _ := filepath.Abs(f)
				t.Error(fmt.Errorf("%s : %v", abs, err))
			}
		})
	}
}
