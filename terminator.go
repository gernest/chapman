package goes

import "fmt"

type terminatorLexer struct{}

func (terminatorLexer) name() string {
	return "terminator"
}

func (terminatorLexer) accept(s scanner) bool {
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

func (t terminatorLexer) lex(s scanner, ctx *context) (*token, error) {
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
			Kind:  lineTerminator,
			Text:  string(n),
			Start: start,
			End:   end}
		switch n {
		case 0x0000A:
			tk.Kind = LF
		case 0x000D:
			tk.Kind = CR
		case 0x02028:
			tk.Kind = LS
		case 0x2029:
			tk.Kind = PS
		}
		return tk, nil
	}
	return nil, fmt.Errorf(unexpectedTkn, t.name(), end)
}
