package goes

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexComment(t *testing.T) {
	sample := []struct {
		src     string
		context string
	}{
		{"// single line comment", "single line commment"},
		{`// single with line terminator 
			
			`, "single with line terminator"},
	}

	for _, v := range sample {
		tkns, err := lex(
			strings.NewReader(v.src),
			terminatorLexer{}, commentLexer{},
		)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(tkns)
	}
}
