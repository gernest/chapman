package lexer

import (
	"bytes"
	"encoding/json"
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
	limit := "fixture/core/uncategorised/101/actual.js"
	for _, f := range files {
		if f == limit {
			t.Run(filepath.Dir(f), func(ts *testing.T) {
				b, err := ioutil.ReadFile(f)
				if err != nil {
					ts.Fatal(err)
				}
				op, err := opts(filepath.Dir(f))
				if err != nil {
					ts.Fatal(err)
				}
				_, err = lex(bytes.NewReader(b), defaultLexMe()...)
				if op.Throws {
					if err == nil {
						abs, _ := filepath.Abs(f)
						t.Error(fmt.Errorf("%s : expected an error", abs))
					}
				} else {
					if err != nil {
						abs, _ := filepath.Abs(f)
						t.Error(fmt.Errorf("%s : %v", abs, err))
					}
				}
			})
		}
	}
}

type options struct {
	Throws bool `json:"throws"`
}

func opts(dir string) (*options, error) {
	b, err := ioutil.ReadFile(filepath.Join(dir, "options.json"))
	if err != nil {
		return nil, err
	}
	o := &options{}
	err = json.Unmarshal(b, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
