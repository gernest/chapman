package goes

import "fmt"

type numeralLexer struct{}

func (numeralLexer) name() string {
	return "numeral"
}
func (numeralLexer) accept(s scanner) bool {
	ch, _, err := s.peek()
	if err != nil {
		return false
	}
	return isDecimalDigit(ch)
}

func isDecimalDigit(ch rune) bool {
	switch ch {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func isNonZeroDigit(ch rune) bool {
	switch ch {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func (n numeralLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	nx, w, err := s.next()
	if err != nil {
		return nil, err
	}
	if isDecimalDigit(nx) {

	}
	return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
}
