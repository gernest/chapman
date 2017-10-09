package chapman

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

type stringLexer struct{}

func (stringLexer) name() string {
	return "string"
}

func isSingleCharacterEscape(ch rune) bool {
	switch ch {
	case '"', 'b', 'f', 'n', 'r', 't', 'v':
		return true
	default:
		s := string(ch)
		return s == "'" || s == "\\"
	}
}

func isNonEscapeChar(ch rune) bool {
	return utf8.ValidRune(ch) && !isEscapeChar(ch)
}

func isEscapeChar(ch rune) bool {
	return isSingleCharacterEscape(ch) ||
		isDecimalDigit(ch) || ch == 'x' || ch == 'u'
}

func (stringLexer) accept(s scanner) bool {
	ch, _, er := s.peek()
	if er != nil {
		return false
	}
	return ch == '"' || string(ch) == "'"
}

func (sl stringLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	ch, w, err := s.next()
	if err != nil {
		return nil, err
	}
	end.Column += w
	var buf bytes.Buffer
	buf.WriteRune(ch)
	if ch == '"' {
		for {
			nx, w, err := s.next()
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
			if string(nx) == "\\" {
				nxt, w, err := s.next()
				if err != nil {
					return nil, err
				}
				end.Column += w
				buf.WriteRune(nxt)
				if isSingleCharacterEscape(nxt) || isNonEscapeChar(nxt) {
					continue
				}
				if nxt == '0' {
					next, w, err := s.next()
					if err != nil {
						return nil, err
					}
					end.Column += w
					buf.WriteRune(next)
					if !isDecimalDigit(next) {
						return nil, fmt.Errorf(unexpectedTkn, sl.name(), end)
					}
				}
				if nxt == 'x' {
					next, w, err := s.next()
					if err != nil {
						return nil, err
					}
					end.Column += w
					buf.WriteRune(next)
					if !isHexDigit(next) {
						return nil, fmt.Errorf(unexpectedTkn, sl.name(), end)
					}
					continue
				}
				if nxt == 'u' {
					next, w, err := s.next()
					if err != nil {
						return nil, err
					}
					end.Column += w
					buf.WriteRune(next)
					if !isHexDigit(next) || next != '{' {
						return nil, fmt.Errorf(unexpectedTkn, sl.name(), end)
					}
					continue
				}
			}
		}
	}
	return nil, fmt.Errorf(unexpectedTkn, sl.name(), end)
}
