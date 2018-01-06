package lexer

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

const backSlash = 0x005C
const singleQuote = 0x2019

type stringLexer struct{}

func (stringLexer) Name() string {
	return "string"
}

func isSingleCharacterEscape(ch rune) bool {
	switch ch {
	case '"', 'b', 'f', 'n', 'r', 't', 'v', backSlash, singleQuote:
		return true
	default:
		return false
	}
}

func isNonEscapeChar(ch rune) bool {
	return utf8.ValidRune(ch) && !isEscapeChar(ch)
}

func isEscapeChar(ch rune) bool {
	return isSingleCharacterEscape(ch) ||
		isDecimalDigit(ch) || ch == 'x' || ch == 'u'
}

func (stringLexer) Accept(s scanner) bool {
	ch, _, er := s.Peek()
	if er != nil {
		return false
	}
	return ch == '"' || string(ch) == "'"
}

func (sl stringLexer) Lex(s scanner, ctx *context) (*token, error) {
	var start position
	if ctx.lastToken != nil {
		start = ctx.lastToken.End
	}
	ch, _, err := s.Next()
	if err != nil {
		return nil, err
	}
	tk := newToken(start)
	tk.AddRune(ch)
	if ch == '"' {
		var gerr error
	main:
		for {
			nx, _, err := s.Next()
			if err != nil {
				if err == io.EOF {
					return nil, errors.New("missing closing string")
				}
				return nil, err
			}
			tk.AddRune(nx)
			if nx == '"' {
				tk.AddRune(nx)
				return tk, nil
			}
			if nx == backSlash {
				nxt, _, err := s.Next()
				if err != nil {
					return nil, err
				}
				tk.AddRune(nxt)
				switch {
				case nxt == '0':
					next, _, err := s.Next()
					if err != nil {
						return nil, err
					}
					tk.AddRune(next)
					if !isDecimalDigit(next) {
						return nil, fmt.Errorf(unexpectedTkn, sl.Name(), tk.End)
					}
					continue
				case nx == 'x':
					next, _, err := s.Next()
					if err != nil {
						return nil, err
					}
					tk.AddRune(next)
					if !isHexDigit(next) {
						return nil, fmt.Errorf(unexpectedTkn, sl.Name(), tk.End)
					}
					continue
				case nxt == 'u':
					next, _, err := s.Next()
					if err != nil {
						return nil, err
					}
					tk.AddRune(next)
					if isHexDigit(next) {
						for i := 0; i < 3; i++ {
							next, _, err = s.Next()
							if err != nil {
								return nil, err
							}
							tk.AddRune(next)
							if !isHexDigit(next) {
								return nil, fmt.Errorf(unexpectedTkn, sl.Name(), tk.End)
							}
						}
						continue
					}
					if next == '{' {
						for {
							next, _, err = s.Next()
							if err != nil {
								return nil, err
							}
							tk.AddRune(next)
							if next == '}' {
								continue
							}
							if !isHexDigit(next) {
								gerr = fmt.Errorf(unexpectedTkn, sl.Name(), tk.End)
								break main
							}
						}
					}
				case isSingleCharacterEscape(nxt) || isNonEscapeChar(nxt):
					continue
				default:
					return nil, fmt.Errorf(unexpectedTkn, sl.Name(), tk.End)
				}
			}
		}
		if gerr != nil {
			return nil, gerr
		}
	}
	return nil, fmt.Errorf(unexpectedTkn, sl.Name(), tk.End)
}
