package lexer

import (
	"fmt"
	"io"
)

type numeralLexer struct{}

func (numeralLexer) Name() string {
	return "numeral"
}

func (numeralLexer) Accept(s scanner) bool {
	ch, _, err := s.Peek()
	if err != nil {
		return false
	}
	return isDecimalDigit(ch) || isFloat(s, ch)
}

func isFloat(s scanner, ch rune) bool {
	if ch != '.' {
		return false
	}
	p, _, err := s.PeekAt(2)
	if err != nil {
		return false
	}
	return isDecimalDigit(p)
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

func isBinaryDigit(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isOctalDigit(ch rune) bool {
	switch ch {
	case '0', '1', '2', '3', '4', '5', '6', '7':
		return true
	default:
		return false
	}
}

func isHexDigit(ch rune) bool {
	switch ch {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F':
		return true
	default:
		return false
	}
}

func (n numeralLexer) Lex(s scanner, ctx *context) (*token, error) {
	var start position
	if ctx.lastToken != nil {
		start = ctx.lastToken.End
	}
	nx, _, err := s.Next()
	if err != nil {
		return nil, err
	}
	tk := newToken(start)
	if isDecimalDigit(nx) || isFloat(s, nx) {
		tk.AddRune(string(nx))
		nxt, _, err := s.Peek()
		if err != nil {
			if err == io.EOF {
				tk.Kind = INT
				return tk, nil
			}
			return nil, err
		}
		switch {
		case isTokenSep(nxt):
			s.Next()
			tk.Kind = INT
			return tk, nil
		case nxt == '.':
			s.Next()
			tk.AddRune(string(nxt))
			ch, _, err := s.Peek()
			if err != nil {
				return nil, err
			}
			switch {
			case isDecimalDigit(ch):
				s.Next()
				tk.AddRune(string(ch))
				for {
					ch, _, err = s.Peek()
					if err != nil {
						if err == io.EOF {
							break
						} else {
							return nil, err
						}
					}
					if isDecimalDigit(ch) {
						s.Next()
						tk.AddRune(string(ch))
						continue
					}
					if ch == 'e' || ch == 'E' {
						s.Next()
						tk.AddRune(string(ch))
						ch, _, err = s.Next()
						if err != nil {
							return nil, err
						}
						if isDecimalDigit(ch) || ch == '+' || ch == '-' {
							tk.AddRune(string(ch))
							continue
						}
						return nil, fmt.Errorf(unexpectedTkn, n.Name(), tk.End)
					}
					if isTokenSep(ch) {
						break
					}
					return nil, fmt.Errorf(unexpectedTkn, n.Name(), tk.End)
				}
				return tk, nil
			}
		case nx == '0' && nxt == 'x' || nx == 0 && nxt == 'X':
			tk.Kind = HEX
			s.Next()
			tk.AddRune(string(nxt))
			for {
				ch, _, err := s.Peek()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isHexDigit(ch) {
					s.Next()
					tk.AddRune(string(ch))
					continue
				}
				break
			}
			return tk, nil
		case nx == '0' && nxt == 'b' || nxt == 'B':
			tk.Kind = BINARY
			for {
				ch, _, err := s.Next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isBinaryDigit(ch) {
					tk.AddRune(string(ch))
					continue
				}
				if isTokenSep(ch) {
					s.Rewind()
					break
				}
				return nil, fmt.Errorf(unexpectedTkn, n.Name(), tk.End)
			}
			return tk, nil
		case nx == '0' && nxt == 'o' || nxt == 'O':
			tk.Kind = OCTAL
			for {
				ch, _, err := s.Next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isOctalDigit(ch) {
					tk.AddRune(string(ch))
					continue
				}
				if isTokenSep(ch) {
					s.Rewind()
					break
				}
				return nil, fmt.Errorf(unexpectedTkn, n.Name(), tk.End)
			}
			return tk, nil
		case isDecimalDigit(nxt):
			tk.Kind = INT
			for {
				ch, _, err := s.Next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isDecimalDigit(ch) {
					tk.AddRune(string(ch))
					continue
				}
				if ch == 'e' || ch == 'E' {
					tk.AddRune(string(ch))
					ch, _, err = s.Next()
					if err != nil {
						return nil, err
					}
					if isDecimalDigit(ch) || ch == '+' || ch == '-' {
						tk.AddRune(string(ch))
						continue
					}
					return nil, fmt.Errorf(unexpectedTkn, n.Name(), tk.End)
				}
				if isTokenSep(ch) {
					break
				}
				return nil, fmt.Errorf(unexpectedTkn, n.Name(), tk.End)
			}
			return tk, nil
		default:
			tk.Kind = INT
			return tk, nil
		}
	}
	return nil, fmt.Errorf(unexpectedTkn, n.Name(), tk.End)
}
