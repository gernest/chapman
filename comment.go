package goes

import (
	"bytes"
	"fmt"
	"io"
)

type singleLineCommentLexer struct{}

func (singleLineCommentLexer) name() string {
	return "singleLineComment"
}

func (singleLineCommentLexer) accept(s scanner) bool {
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
			return true
		}
	}
	return false
}

func (c singleLineCommentLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	n, w, err := s.next()
	if err != nil {
		return nil, err
	}
	end.Column += w
	if n == '/' {
		nx, w, err := s.next()
		if err != nil {
			return nil, err
		}
		end.Column += w
		if nx == '/' {
			var b bytes.Buffer
			tk := &token{Kind: SingleLineComment}
			for {
				x, w, err := s.next()
				if err != nil {
					if err == io.EOF {
						tk.Text = b.String()
						tk.Start = start
						tk.End = end
						return tk, nil
					}
					return nil, err
				}
				if isLineTerminator(x) {
					tk.Text = b.String()
					tk.Start = start
					tk.End = end
					return tk, nil
				}
				end.Column += w
				b.WriteRune(x)
			}
		}
	}
	return nil, fmt.Errorf(unexpectedTkn, c.name(), end)
}
