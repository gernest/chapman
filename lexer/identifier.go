package lexer

import (
	"bytes"
	"fmt"
	"io"
)

const reverseSolidus = 0x005C // backslash

type identifierNameLexer struct{}

func (identifierNameLexer) Name() string {
	return "identifierName"
}

func (identifierNameLexer) Accept(s scanner) bool {
	ch, _, err := s.Peek()
	if err != nil {
		return false
	}
	return isUnicodeIDStart(ch) ||
		ch == '$' || ch == '_' || escapeSequence(ch, s)
}

func escapeSequence(ch rune, s scanner) bool {
	if ch == reverseSolidus {
		n, _, err := s.PeekAt(2)
		if err != nil {
			return false
		}
		return n == 'u'
	}
	return false
}

func (i identifierNameLexer) Lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	var b bytes.Buffer
	tk := &token{Kind: IdentifierName, Start: start}
	e, err := i.lexStart(s, ctx, end, &b)
	if err != nil {
		return nil, err
	}
	e, err = i.lexPart(s, ctx, *e, &b)
	if err != nil {
		return nil, err
	}
	end = *e
	tk.End = end
	tk.Text = b.String()
	return tk, nil
}

func (i identifierNameLexer) lexStart(s scanner, ctx *context, end position, b *bytes.Buffer) (*position, error) {
	n, w, err := s.Next()
	if err != nil {
		return nil, err
	}
	end.Column += w
	if isUnicodeIDStart(n) || n == '$' || n == '_' {
		b.WriteRune(n)
	} else if n == reverseSolidus {
		b.WriteRune(n)
		nx, w, err := s.Next()
		if err != nil {
			return nil, err
		}
		end.Column += w
		if nx != 'u' {
			return nil, fmt.Errorf(unexpectedTkn, i.Name(), end)
		}
		b.WriteRune(nx)

		// wer are lexing a valid UnicodeEscapeSequence
		nx, w, err = s.Next()
		if err != nil {
			return nil, err
		}
		end.Column += w
		b.WriteRune(nx)
		if isHexDigit(nx) {
			// four hex digits, we already have one three to go
			for k := 0; k < 3; k++ {
				nx, w, err = s.Next()
				if err != nil {
					return nil, err
				}
				end.Column += w
				b.WriteRune(nx)
				if !isHexDigit(nx) {
					return nil, fmt.Errorf(unexpectedTkn, i.Name(), end)
				}
			}
			return &end, nil
		}
		if nx == '{' {
			for {
				nx, w, err = s.Next()
				if err != nil {
					return nil, err
				}
				end.Column += w
				b.WriteRune(nx)
				if !isHexDigit(nx) {
					if nx == '}' {
						return &end, nil
					}
					return nil, fmt.Errorf(unexpectedTkn, i.Name(), end)
				}
			}
		}
	}
	return &end, nil
}
func (i identifierNameLexer) lexPart(s scanner, ctx *context, end position, b *bytes.Buffer) (*position, error) {
	for {
		if i.Accept(s) {
			e, err := i.lexStart(s, ctx, end, b)
			if err != nil {
				return nil, err
			}
			end.Column = e.Column
		} else {
			nx, w, err := s.Peek()
			if err != nil {
				if err == io.EOF {
					return &end, nil
				}
				return nil, err
			}
			if isUnicodeIDContinue(nx) || nx == 0x200C || nx == 0x200D {
				s.Next()
				end.Column += w
				b.WriteRune(nx)
				continue
			}
			return &end, nil
		}
	}
}
