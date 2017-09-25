package goes

import (
	"bytes"
	"fmt"
)

const reverseSolidus = 0x005C // backslash

type identifierNameLexer struct{}

func (identifierNameLexer) name() string {
	return "identifierName"
}

func (identifierNameLexer) accept(s scanner) bool {
	ch, _, err := s.peek()
	if err != nil {
		return false
	}
	return isUnicodeIDStart(ch) ||
		ch == '$' || ch == '_' || escapeSequence(ch, s)
}

func escapeSequence(ch rune, s scanner) bool {
	if ch == reverseSolidus {
		n, _, err := s.peek()
		if err != nil {
			return false
		}
		return n == 'u'
	}
	return false
}

func (i identifierNameLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	var b bytes.Buffer
	tk := &token{Kind: IdentifierName, Start: start}
	for i.accept(s) {
		e, err := i.lexStart(s, ctx, end, &b)
		if err != nil {
			return nil, err
		}
		e, err = i.lexPart(s, ctx, *e, &b)
		if err != nil {
			return nil, err
		}
		end = *e
	}
	tk.End = end
	tk.Text = b.String()
	return tk, nil
}

func (i identifierNameLexer) lexStart(s scanner, ctx *context, end position, b *bytes.Buffer) (*position, error) {
	n, w, err := s.next()
	if err != nil {
		return nil, err
	}
	end.Column += w
	if isUnicodeIDStart(n) || n == '$' || n == '_' {
		b.WriteRune(n)
	} else if n == reverseSolidus {
		nx, w, err := s.next()
		if err != nil {
			return nil, err
		}
		end.Column += w
		if nx != 'u' {
			return nil, fmt.Errorf(unexpectedTkn, i.name(), end)
		}
		b.WriteRune(nx)

		// wer are lexing a valid UnicodeEscapeSequence
		nx, w, err = s.next()
		if err != nil {
			return nil, err
		}
		end.Column += w
		b.WriteRune(nx)
		if isHexDigit(nx) {
			// four hex digits, we already have one three to go
			for k := 0; k < 3; k++ {
				nx, w, err = s.next()
				if err != nil {
					return nil, err
				}
				end.Column += w
				b.WriteRune(nx)
				if !isHexDigit(nx) {
					return nil, fmt.Errorf(unexpectedTkn, i.name(), end)
				}
			}
			return &end, nil
		}
		if nx == '{' {
			for {
				nx, w, err = s.next()
				if err != nil {
					return nil, err
				}
				end.Column += w
				b.WriteRune(nx)
				if !isHexDigit(nx) {
					if nx == '}' {
						return &end, nil
					}
					return nil, fmt.Errorf(unexpectedTkn, i.name(), end)
				}
			}
		}
	}
	return &end, nil
}
func (i identifierNameLexer) lexPart(s scanner, ctx *context, end position, b *bytes.Buffer) (*position, error) {
	if i.accept(s) {
		return i.lexStart(s, ctx, end, b)
	}
	nx, _, err := s.peek()
	if err != nil {
		return nil, err
	}
	if isUnicodeIDContinue(nx) {
		nx, w, err := s.next()
		if err != nil {
			return nil, err
		}
		end.Column += w
		b.WriteRune(nx)
		return &end, nil
	}
	return &end, nil
}
