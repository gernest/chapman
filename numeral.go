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
	if isDecimalDigit(nx) || isFloat(s, nx) {
		end.Column += w
		tk := &token{Start: start, Text: string(nx)}
		nxt, w, err := s.next()
		if err != nil {
			if err == io.EOF {
				tk.Kind = INT
				return tk, nil
			}
			return nil, err
		}
		if isTokenSep(nxt) {
			s.rewind()
			tk.Kind = INT
			return tk, nil
		}
		end.Column += w
		tk.Text += string(nxt)

		if nxt == '.' {
			ch, w, err := s.next()
			if err != nil {
				return nil, err
			}
			if isDecimalDigit(ch) {
				tk.Kind = FLOAT
				end.Column += w
				tk.Text += string(ch)
				for {
					ch, w, err = s.next()
					if err != nil {
						if err == io.EOF {
							break
						} else {
							return nil, err
						}
					}
					if isDecimalDigit(ch) {
						end.Column += w
						tk.Text += string(ch)
						continue
					}
					if ch == 'e' || ch == 'E' {
						tk.Text += string(ch)
						end.Column += w
						ch, w, err = s.next()
						if err != nil {
							return nil, err
						}
						if isDecimalDigit(ch) || ch == '+' || ch == '-' {
							tk.Text += string(ch)
							end.Column += w
							continue
						}
						return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
					}
					if isTokenSep(ch) {
						break
					}
					return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
				}
				return tk, nil
			}
		}
		if nx == '0' && nxt == 'x' || nxt == 'X' {
			tk.Kind = HEX
			for {
				ch, w, err := s.next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isHexDigit(ch) {
					tk.Text += string(ch)
					end.Column += w
					continue
				}
				if isTokenSep(ch) {
					s.rewind()
					break
				}
				return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
			}
			return tk, nil
		}
		if nx == '0' && nxt == 'b' || nxt == 'B' {
			tk.Kind = BINARY
			for {
				ch, w, err := s.next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isBinaryDigit(ch) {
					tk.Text += string(ch)
					end.Column += w
					continue
				}
				if isTokenSep(ch) {
					s.rewind()
					break
				}
				return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
			}
			return tk, nil
		}
		if nx == '0' && nxt == 'o' || nxt == 'O' {
			tk.Kind = OCTAL
			for {
				ch, w, err := s.next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isOctalDigit(ch) {
					tk.Text += string(ch)
					end.Column += w
					continue
				}
				if isTokenSep(ch) {
					s.rewind()
					break
				}
				return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
			}
			return tk, nil
		}
		if isDecimalDigit(nxt) {
			tk.Kind = INT
			for {
				ch, w, err := s.next()
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return nil, err
					}
				}
				if isDecimalDigit(ch) {
					tk.Text += string(ch)
					end.Column += w
					continue
				}
				if ch == 'e' || ch == 'E' {
					tk.Text += string(ch)
					end.Column += w
					ch, w, err = s.next()
					if err != nil {
						return nil, err
					}
					if isDecimalDigit(ch) || ch == '+' || ch == '-' {
						tk.Text += string(ch)
						end.Column += w
						continue
					}
					return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
				}
				if isTokenSep(ch) {
					break
				}
				return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
			}
			return tk, nil
		}

	}

	return nil, fmt.Errorf(unexpectedTkn, n.name(), end)
}
