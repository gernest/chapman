package lexer

import (
	"fmt"
	"io"
)

var puncs = map[string]bool{
	"{":    true,
	"(":    true,
	")":    true,
	"[":    true,
	"]":    true,
	".":    true,
	"...":  true,
	";":    true,
	",":    true,
	"<":    true,
	">":    true,
	"<=":   true,
	">=":   true,
	"==":   true,
	"!=":   true,
	"===":  true,
	"!==":  true,
	"+":    true,
	"-":    true,
	"*":    true,
	"%":    true,
	"**":   true,
	"++":   true,
	"--":   true,
	"<<":   true,
	">>":   true,
	">>>":  true,
	"&":    true,
	"|":    true,
	"^":    true,
	"!":    true,
	"~":    true,
	"&&":   true,
	"||":   true,
	"?":    true,
	":":    true,
	"=":    true,
	"+=":   true,
	"-=":   true,
	"*=":   true,
	"%=":   true,
	"**=":  true,
	"<<=":  true,
	">>=":  true,
	">>>=": true,
	"&=":   true,
	"|=":   true,
	"^=":   true,
	"=>":   true,
	"/":    true,
	"}":    true,
}

type punctuatorLexer struct{}

func (punctuatorLexer) Name() string {
	return "punctuator"
}

func (punctuatorLexer) Accept(s scanner) bool {
	ch, _, err := s.Peek()
	if err != nil {
		return false
	}
	return isPunctuator(string(ch))
}

func isPunctuator(ch string) bool {
	return puncs[ch]
}

func (p punctuatorLexer) Lex(s scanner, ctx *context) (*token, error) {
	var start position
	if ctx.lastToken != nil {
		start = ctx.lastToken.End
	}
	nx, _, err := s.Next()
	if err != nil {
		return nil, err
	}
	tk := newToken(start)
	tk.AddText(string(nx))

	switch nx {
	case '{':
		tk.Kind = LBRACE
		return tk, nil
	case '}':
		tk.Kind = RBRACE
		return tk, nil
	case '(':
		tk.Kind = LPAREN
		return tk, nil
	case ')':
		tk.Kind = RPAREN
		return tk, nil
	case '[':
		tk.Kind = LBRACK
		return tk, nil
	case ']':
		tk.Kind = RBRACK
		return tk, nil
	case '.':
		tk.Kind = PERIOD

		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '.':
				s.Next()
				tk.AddText(string(nxt))

				nxt, _, err = s.Next()
				if err != nil {
					return nil, err
				}
				if nxt != '.' {
					return nil, fmt.Errorf(unexpectedTkn, p.Name(), tk.End)
				}
				tk.AddText(string(nxt))
				tk.Kind = ELLIPSIS
			}
		}
		return tk, nil
	// case "...":
	// 	tk.Kind = ELLIPSIS
	// 	return tk, nil
	case ';':
		tk.Kind = SEMICOLON
		return tk, nil
	case ',':
		tk.Kind = COMMA
		return tk, nil
	case '<':
		tk.Kind = LSS
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '<':
				s.Next()
				tk.Kind = SHL
				tk.AddText(string(nxt))

				if p.Accept(s) {
					nxt, _, err = s.Next()
					if err == io.EOF {
						return tk, nil
					}
					if err != nil {
						return nil, err
					}
					if nxt == '=' {
						s.Next()
						tk.Kind = SHLAssign
						tk.AddText(string(nxt))
					}
				}

			case '=':
				// We advance the cursor since we already read the rune thorugh peek.
				s.Next()
				tk.Kind = LEQ
				tk.AddText(string(nxt))
			}
		}
		return tk, nil
	case '>':
		tk.Kind = GTR
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '=': //>=

				// We advance the cursor since we already read the rune thorugh peek.
				s.Next()
				tk.Kind = GEQ
				tk.AddText(string(nxt))
			case '>': //>>
				tk.Kind = SHR
				s.Next()
				tk.Kind = GEQ
				tk.AddText(string(nxt))

				if p.Accept(s) {
					nxt, _, err = s.Peek()
					if err == io.EOF {
						return tk, nil
					}
					if err != nil {
						return nil, err
					}
					switch nxt {

					case '=': //>>=
						tk.Kind = SHRAssign
						s.Next()
						tk.Kind = GEQ
						tk.AddText(string(nxt))
					case '>': //>>>
						tk.Kind = USHR
						s.Next()
						tk.AddText(string(nxt))
						if p.Accept(s) {
							nxt, _, err = s.Peek()
							if err == io.EOF {
								return tk, nil
							}
							if err != nil {
								return nil, err
							}
							if nxt == '=' { ///>>>=
								s.Next()
								tk.Kind = USHRAssign
								tk.AddText(string(nxt))
							}
						}
					}
				}
			}
		}
		return tk, nil
	case '+':
		tk.Kind = ADD
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '=':
				s.Next()
				tk.Kind = AddAssign
				tk.AddText(string(nxt))
			case '+':
				s.Next()
				tk.Kind = INC
				tk.AddText(string(nxt))
			}
		}
		return tk, nil
	case '-':
		tk.Kind = SUB
		if p.Accept(s) {
			nxt, _, err := s.Next()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '-':
				s.Next()
				tk.Kind = DEC
				tk.AddText(string(nxt))
			case '=':
				s.Next()
				tk.Kind = SubAssign
				tk.AddText(string(nxt))
			}
		}
		return tk, nil
	case '*':
		tk.Kind = MUL
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '=':
				s.Next()
				tk.Kind = MulAssign
				tk.AddText(string(nxt))
			case '*':
				s.Next()
				tk.Kind = EXP
				tk.AddText(string(nxt))

				if p.Accept(s) {
					nxt, _, err := s.Peek()
					if err == io.EOF {
						return tk, nil
					}
					if err != nil {
						return nil, err
					}
					switch nxt {
					case '=':
						s.Next()
						tk.Kind = EXP
						tk.AddText(string(nxt))
					}
				}
			}
		}
		return tk, nil
	case '%':
		tk.Kind = REM
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '=':
				s.Next()
				tk.Kind = RemAssign
				tk.AddText(string(nxt))
			}
		}
		return tk, nil
	// case "**":
	// 	tk.Kind = EXP
	// 	return tk, nil
	// case "++":
	// 	tk.Kind = INC
	// 	return tk, nil
	// case "--":
	// 	tk.Kind = DEC
	// 	return tk, nil
	// case "<<":
	// 	tk.Kind = SHL
	// 	return tk, nil
	// case ">>":
	// 	tk.Kind = SHR
	// 	return tk, nil
	// case ">>>":
	// 	tk.Kind = USHR
	// 	return tk, nil
	case '&':
		tk.Kind = AND
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '=':
				s.Next()
				tk.Kind = AndAssign
				tk.AddText(string(nxt))
			case '&':
				s.Next()
				tk.Kind = LAND
				tk.AddText(string(nxt))
			}
		}
		return tk, nil
	case '|':
		tk.Kind = OR
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '|':
				s.Next()
				tk.Kind = LOR
				tk.AddText(string(nxt))
			case '=':
				s.Next()
				tk.Kind = OrAssign
				tk.AddText(string(nxt))
			}
		}
		return tk, nil
	case '^':
		tk.Kind = XOR
		if p.Accept(s) {
			nxt, _, err := s.Next()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '=':
				s.Next()
				tk.Kind = XorAssign
				tk.AddText(string(nxt))
			}
		}
		return tk, nil
	case '!':
		tk.Kind = NOT
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '=':
				// We advance the cursor since we already read the rune thorugh peek.
				s.Next()
				tk.Kind = NEQ
				tk.AddText(string(nxt))
				if p.Accept(s) {
					nxt, _, err = s.Peek()
					if err == io.EOF {
						return tk, nil
					}
					if err != nil {
						return nil, err
					}
					if nxt == '=' {
						s.Next()
						tk.Kind = SNEQ
						tk.AddText(string(nxt))
					}
				}
			}
		}
		return tk, nil
	case '~':
		tk.Kind = TILDE
		return tk, nil
	// case "||":
	// tk.Kind = LOR
	// return tk, nil
	case '?':
		tk.Kind = QN
		return tk, nil
	case ':':
		tk.Kind = COLON
		return tk, nil
	case '=':
		tk.Kind = ASSIGN
		if p.Accept(s) {
			nxt, _, err := s.Peek()
			if err == io.EOF {
				return tk, nil
			}
			if err != nil {
				return nil, err
			}
			switch nxt {
			case '>':
				s.Next()
				tk.Kind = ARROW
				tk.AddText(string(nxt))
			case '=':
				// We advance the cursor since we already read the rune thorugh peek.
				s.Next()
				tk.Kind = EQL
				tk.AddText(string(nxt))
				if p.Accept(s) {
					nxt, _, err = s.Peek()
					if err == io.EOF {
						return tk, nil
					}
					if err != nil {
						return nil, err
					}
					if nxt == '=' {
						s.Next()
						tk.Kind = SEQL
						tk.AddText(string(nxt))
					}
				}
			}
		}
		return tk, nil
	// case "+=":
	// 	tk.Kind = AddAssign
	// 	return tk, nil
	// case "-=":
	// 	tk.Kind = SubAssign
	// 	return tk, nil
	// case "*=":
	// 	tk.Kind = MulAssign
	// 	return tk, nil
	// case "%=":
	// 	tk.Kind = RemAssign
	// 	return tk, nil
	// case "**=":
	// 	tk.Kind = ExpAssign
	// 	return tk, nil
	// case "<<=":
	// 	tk.Kind = SHLAssign
	// 	return tk, nil
	// case ">>=":
	// 	tk.Kind = SHRAssign
	// 	return tk, nil
	// case ">>>=":
	// 	tk.Kind = SHRAssign
	// 	return tk, nil
	// case "|=":
	// 	tk.Kind = OrAssign
	// 	return tk, nil
	// case "^=":
	// 	tk.Kind = XorAssign
	// 	return tk, nil
	// case "=>":
	// 	tk.Kind = ARROW
	// 	return tk, nil
	case '/':
		tk.Kind = QUO
		return tk, nil
	// case "/=":
	// 	tk.Kind = QuoAssign
	// 	return tk, nil
	default:
		return nil, fmt.Errorf(unexpectedTkn, p.Name(), tk.End)
	}
}
