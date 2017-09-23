package goes

import (
	"bytes"
	"fmt"
	"io"
)

type commentLexer struct{}

func (commentLexer) name() string {
	return "comment"
}

func (commentLexer) accept(s scanner) bool {
	n, _, err := s.peekAt(1)
	if err != nil {
		return false
	}
	if n == '/' {
		nx, _, err := s.peekAt(2)
		if err != nil {
			return false
		}
		if nx == '/' {
			fmt.Println("is a comment")
			return true
		}
	}
	return false
}

func (c commentLexer) lex(s scanner, ctx *context) (*token, error) {
	n, _, err := s.next()
	if err != nil {
		return nil, err
	}
	fmt.Println("first", string(n))
	if n == '/' {
		nx, _, err := s.next()
		if err != nil {
			return nil, err
		}
		fmt.Println("second", string(nx))
		if nx == '/' {
			var b bytes.Buffer
			tk := &token{kind: comment}
			for {
				x, _, err := s.next()
				if err != nil {
					if err == io.EOF {
						tk.text = b.String()
						return tk, nil
					}
					return nil, err
				}
				if isLineTerminator(x) {
					s.rewind()
					tk.text = b.String()
					return tk, nil
				}
				b.WriteRune(x)
			}
		}
		return nil, fmt.Errorf(unexpectedTkn, c.name(), "/", string(nx))
	}
	return nil, fmt.Errorf(unexpectedTkn, c.name(), "/", string(n))
}
