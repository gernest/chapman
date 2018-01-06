package lexer

import (
	"bytes"
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
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	ch, w, err := s.Next()
	if err != nil {
		return nil, err
	}
	end.Column += w
	var buf bytes.Buffer
	buf.WriteRune(ch)
	if ch == '"' {
		var gerr error
	main:
		for {
			nx, w, err := s.Next()
			if err != nil {
				if err == io.EOF {
					return nil, errors.New("missing closing string")
				}
				return nil, err
			}
			end.Column += w
			buf.WriteRune(nx)
			if nx == '"' {
				tk := &token{Start: start, Text: buf.String(), End: end}
				return tk, nil
			}
			if nx == backSlash {
				nxt, w, err := s.Next()
				if err != nil {
					return nil, err
				}
				end.Column += w
				buf.WriteRune(nxt)
				switch {
				case nxt == '0':
					next, w, err := s.Next()
					if err != nil {
						return nil, err
					}
					end.Column += w
					buf.WriteRune(next)
					if !isDecimalDigit(next) {
						return nil, fmt.Errorf(unexpectedTkn, sl.Name(), end)
					}
					continue
				case nx == 'x':
					next, w, err := s.Next()
					if err != nil {
						return nil, err
					}
					end.Column += w
					buf.WriteRune(next)
					if !isHexDigit(next) {
						return nil, fmt.Errorf(unexpectedTkn, sl.Name(), end)
					}
					continue
				case nxt == 'u':
					next, w, err := s.Next()
					if err != nil {
						return nil, err
					}
					end.Column += w
					buf.WriteRune(next)
					if isHexDigit(next) {
						for i := 0; i < 3; i++ {
							next, w, err = s.Next()
							if err != nil {
								return nil, err
							}
							end.Column += w
							buf.WriteRune(next)
							if !isHexDigit(next) {
								return nil, fmt.Errorf(unexpectedTkn, sl.Name(), end)
							}
						}
						continue
					}
					if next == '{' {
						for {
							next, w, err = s.Next()
							if err != nil {
								return nil, err
							}
							end.Column += w
							buf.WriteRune(next)
							if next == '}' {
								continue
							}
							if !isHexDigit(next) {
								gerr = fmt.Errorf(unexpectedTkn, sl.Name(), end)
								break main
							}
						}
					}
				case isSingleCharacterEscape(nxt) || isNonEscapeChar(nxt):
					continue
				default:
					return nil, fmt.Errorf(unexpectedTkn, sl.Name(), end)
				}
			}
		}
		if gerr != nil {
			return nil, gerr
		}
	}
	return nil, fmt.Errorf(unexpectedTkn, sl.Name(), end)
}
