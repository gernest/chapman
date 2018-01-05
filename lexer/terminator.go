package lexer

import "fmt"

type lineTerminatorLexer struct{}

func (lineTerminatorLexer) name() string {
	return "LineTerminator"
}

func (lineTerminatorLexer) accept(s scanner) bool {
	n, _, err := s.peek()
	if err != nil {
		return false
	}
	return isLineTerminator(n)
}

func isLineTerminator(ch rune) bool {
	switch ch {
	case 0x0000A, 0x000D, 0x02028, 0x2029:
		return true
	default:
		return false
	}
}

func (t lineTerminatorLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	n, _, err := s.next()
	if err != nil {
		return nil, err
	}
	if isLineTerminator(n) {
		end.Line++
		end.Column = 0
		tk := &token{
			Text:  string(n),
			Start: start,
			End:   end}
		switch n {
		case 0x0000A:
			tk.Kind = LF
		case 0x000D:
			tk.Kind = CR
			nxt, _, err := s.peek()
			if err == nil {
				//http://ecma-international.org/ecma-262/#sec-line-terminators
				//
				// Treat <CR><LF> as <CR>.
				if nxt == 0x0000A {
					s.next()
				}
			}
		case 0x02028:
			tk.Kind = LS
		case 0x2029:
			tk.Kind = PS
		}
		return tk, nil
	}
	return nil, fmt.Errorf(unexpectedTkn, t.name(), end)
}
