package goes

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type kind uint
type token struct {
	pos    int
	offset int
	text   string
	line   int
	kind   kind
}

const (
	comment kind = iota
	eof
)

type lexer interface {
	next() (rune, error)
	peek() (rune, error)
	rewind() error
}

func lex(src io.Reader) ([]*token, error) {
	r := bufio.NewReader(src)
	var currentLine int
	var currentPosition int
	var tkns []*token
	for {
		ch, n, err := r.ReadRune()
		if err != nil {
			if err == io.EOF && len(tkns) != 0 {
				return tkns, nil
			}
			return nil, err
		}
		if unicode.IsSpace(ch) {
			if ch == '\n' || ch == '\r' {
				currentLine++
				currentPosition = 0
			} else {
				currentPosition += n
			}
		}
		var b bytes.Buffer
		switch ch {
		case '/':
			nxt, _, err := r.ReadRune()
			if err != nil {
				return nil, err
			}
			if nxt == '/' { //comment
				tk := &token{kind: comment, line: currentLine}
				b.Reset()
				for {
					x, _, err := r.ReadRune()
					if err != nil {
						if err == io.EOF {
							tk.text = b.String()
							tkns = append(tkns, tk)
							break
						}
						return nil, err
					}
					switch x {
					case '\r', '\n':
						tk.text = b.String()
						tkns = append(tkns, tk)
						break
					default:
						b.WriteRune(x)
					}
				}
			}
		}
	}
}
