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
	case 0x00A, 0x000D, 0x02028, 0x2029:
		return true
	default:
		return false
	}
}

func (t terminatorLexer) lex(s scanner, ctx *context) (*token, error) {
	n, _, err := s.next()
	if err != nil {
		return nil, err
	}
	if isLineTerminator(n) {
		return &token{kind: lineTerminator, text: string(n)}, nil
	}
	return nil, fmt.Errorf(unexpectedTkn, t.name(), "line terminator", string(n))
}
