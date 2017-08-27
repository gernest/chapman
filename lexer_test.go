package goes

import (
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func TestLexComment(t *testing.T) {
	sample := []struct {
		src     string
		context string
	}{
		{"// single line comment", "single line commment"},
	}

	for _, v := range sample {
		tkns, err := lex(strings.NewReader(v.src))
		if err != nil {
			t.Fatal(err)
		}
		pretty.Println(tkns)
	}
}
