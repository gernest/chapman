package chapman

import (
	"fmt"
	"io"
)

type numeralLexer struct{}

func (numeralLexer) name() string {
	return "numeral"
}
func (numeralLexer) accept(s scanner) bool {
	ch, _, err := s.peek()
	if err != nil {
		return false
	}
	return isDecimalDigit(ch) || isFloat(s, ch)
}

func isFloat(s scanner, ch rune) bool {
	if ch != '.' {
		return false
	}
	p, _, err := s.peekAt(2)
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

func (n numeralLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	nx, w, err := s.next()
	if err != nil {
		return nil, err
	}
	end.Column += w
	tk := &token{Start: start, Text: string(nx)}
	if nx == '0' {
		nxt, w, err := s.next()
		if err != nil {
			if err == io.EOF {
				//decimal
				tk.End = end
				tk.Kind = INT
				return tk, nil
			}
			return nil, err
		}
		if isLineTerminator(nxt) || isWhiteSpace(nxt) {
			if err := s.rewind(); err != nil {
				return nil, err
			}
			tk.Kind = INT
			tk.End = end
			return tk, nil
		}
		end.Column += w
		tk.Text += string(nxt)
		switch nxt {
		case 'b', 'B': //binary
			txt, err := digits(n, s, &end, isBinaryDigit)
			if err != nil {
				return nil, err
			}
			tk.Text += txt
			tk.End = end
			tk.Kind = BINARY
			return tk, nil
		case 'o', 'O':
			txt, err := digits(n, s, &end, isOctalDigit)
			if err != nil {
				return nil, err
			}
			tk.Text += txt
			tk.End = end
			tk.Kind = OCTAL
			return tk, nil

		case 'x', 'X': //hexadecimal
			txt, err := digits(n, s, &end, isHexDigit)
			if err != nil {
				return nil, err
			}
			tk.Text += txt
			tk.End = end
			tk.Kind = OCTAL
			return tk, nil
		case '.':
			p, _, err := s.peek()
			if err != nil {
				return nil, err
			}
			if isDecimalDigit(p) {
				txt, err := digits(n, s, &end, isDecimalDigit)
				if err != nil {
					return nil, err
				}
				tk.Text += txt
			}
			if p == 'e' {
				txt, err := n.exponent(s, &end)
				if err != nil {
					return nil, err
				}
				tk.Text += txt
			}
			tk.Kind = FLOAT
			tk.End = end
			return tk, nil
		default:
			return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
		}
	}
	if isNonZeroDigit(nx) {
		// txt, err := digits(n, s, &end, isDecimalDigit)
		// if err != nil {
		// 	return nil, err
		// }
		// tk.Text += txt

	}
	return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
}

func digits(l lexMe, s scanner, end *position, accCond func(rune) bool) (string, error) {
	digits := ""
	for {
		next, w, err := s.next()
		if err != nil {
			if err == io.EOF {
				if digits == "" {
					return "", fmt.Errorf(unexpectedTkn, l.name(), end)
				}
				break
			}
			return "", err
		}
		if accCond(next) {
			end.Column += w
			digits += string(next)
			continue
		}
		if isWhiteSpace(next) || isLineTerminator(next) {
			if err := s.rewind(); err != nil {
				return "", err
			}
			break
		}
		return "", fmt.Errorf(unexpectedTkn, l.name(), end)
	}
	return digits, nil
}

func (n numeralLexer) exponent(s scanner, end *position) (string, error) {
	e := ""
	nx, w, err := s.next()
	if err != nil {
		return "", err
	}
	end.Column += w
	e += string(nx)
	if nx == 'e' || nx == 'E' {
		p, w, err := s.peek()
		if err != nil {
			return "", err
		}
		if p == '+' || p == '-' {
			s.next()
			e += string(p)
			end.Column += w
		}
		txt, err := digits(n, s, end, isDecimalDigit)
		if err != nil {
			return "", err
		}
		e += txt
		return e, nil
	}
	return "", fmt.Errorf(unexpectedTkn, n.name(), end)
}
