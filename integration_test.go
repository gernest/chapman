package chapman

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

type unit struct {
	dir    string
	src    string
	expect string
	desc   string
}

var cases map[string][]unit

func init() {
	cases = make(map[string][]unit)
	base := "fixtures"
	o, err := ioutil.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range o {
		var u []unit
		ok := make(map[string]bool)
		ferr := filepath.Walk(filepath.Join(base, d.Name()), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			r := filepath.Dir(path)
			if ok[r] {
				return nil
			}
			var e unit
			for _, p := range []string{"src", "expect", "desc"} {
				b, err := ioutil.ReadFile(filepath.Join(r, p))
				if err != nil {
					return err
				}
				switch p {
				case "src":
					e.src = string(b)
				case "expect":
					e.expect = string(b)
				case "desc":
					e.desc = string(b)
				}
			}
			ok[r] = true
			e.dir = r
			u = append(u, e)
			return nil
		})
		if ferr != nil {
			log.Fatal(ferr)
		}
		cases[d.Name()] = u
	}
}

func TestIntegrationLex(t *testing.T) {
	for k, v := range cases {
		t.Run(k, runCase(v))
	}
}

func runCase(u []unit) func(*testing.T) {
	var b bytes.Buffer
	return func(t *testing.T) {
		b.Reset()
		for _, v := range u {
			b.WriteString(v.src)
			tk, err := lex(&b, allLexme()...)
			if err != nil {
				t.Fatal(err)
			}
			o, err := printTokens(tk)
			if err != nil {
				t.Fatal(err)
			}
			s := string(o)
			if s != v.expect {
				t.Errorf("%s\n%s\nexpected:\n%s\n got:\n %s", v.dir, v.desc, v.expect, s)
			}
		}
	}
}

func allLexme() []lexMe {
	return []lexMe{
		singleLineCommentLexer{}, multiLineCommentLexer{}, lineTerminatorLexer{},
		whiteSpaceLexer{},
	}
}
