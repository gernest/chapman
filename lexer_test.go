package goes

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type caseFixture struct {
	dir         string
	description string
	actual      string
	expected    string
}

func cases(dir string) ([]caseFixture, error) {
	lookup := make(map[string]caseFixture)

	ferr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		key := filepath.Dir(path)
		base := filepath.Base(path)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if v, ok := lookup[key]; ok {
			switch base {
			case "actual":
				v.actual = string(b)
			case "expect":
				v.expected = string(b)
			case "desc":
				v.description = string(b)
			}
			lookup[key] = v
			return nil
		}
		v := caseFixture{dir: key}
		switch base {
		case "actual":
			v.actual = string(b)
		case "expect":
			v.expected = string(b)
		case "desc":
			v.description = string(b)
		}
		lookup[key] = v
		return nil
	})
	if ferr != nil {
		return nil, ferr
	}
	var c []caseFixture
	for _, v := range lookup {
		c = append(c, v)
	}
	return c, nil
}

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
