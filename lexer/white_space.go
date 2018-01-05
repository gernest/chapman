package lexer

import (
	"fmt"
	"unicode"
)

type whiteSpaceLexer struct{}

func (whiteSpaceLexer) name() string {
	return "whiteSpaceLexer"
}

func (whiteSpaceLexer) accept(s scanner) bool {
	n, _, err := s.peek()
	if err != nil {
		return false
	}
	return isWhiteSpace(n)
}

func isWhiteSpace(ch rune) bool {
	switch ch {
	case 0x0009, 0x000B, 0x000C, 0x0020, 0x00A0, 0xFEFF:
		return true
	default:
		return unicode.IsSpace(ch)
	}
}

func (w whiteSpaceLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	n, size, err := s.next()
	if err != nil {
		return nil, err
	}
	end.Column += size
	if isWhiteSpace(n) {
		tk := &token{
			Start: start,
			End:   end,
			Text:  string(n),
		}
		switch n {
		case 0x0009:
			tk.Kind = TAB
		case 0x000B:
			tk.Kind = VT
		case 0x000C:
			tk.Kind = FF
		case 0x0020:
			tk.Kind = SP
		case 0x00A0:
			tk.Kind = NBSP
		case 0xFEFF:
			tk.Kind = ZWNBSP
		default:
			if unicode.IsSpace(n) {
				tk.Kind = USP
			}
		}
		return tk, nil
	}
	return nil, fmt.Errorf(unexpectedTkn, w.name(), end)
}
