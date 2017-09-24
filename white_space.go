package goes

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

func (w whiteSpaceLexer) lext(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.end, start
	}
	n, size, err := s.next()
	if err != nil {
		return nil, err
	}
	end.column += size
	if isWhiteSpace(n) {
		tk := &token{
			start: start,
			end:   end,
			text:  string(n),
		}
		switch n {
		case 0x0009:
			tk.kind = TAB
		case 0x000B:
			tk.kind = VT
		case 0x000C:
			tk.kind = FF
		case 0x0020:
			tk.kind = SP
		case 0x00A0:
			tk.kind = NBSP
		case 0xFEFF:
			tk.kind = ZWNBSP
		default:
			if unicode.IsSpace(n) {
				tk.kind = USP
			}
		}
		return tk, nil
	}
	return nil, fmt.Errorf(unexpectedTkn, w.name(), end)
}
