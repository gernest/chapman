package chapman

import (
	"fmt"
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

func (punctuatorLexer) name() string {
	return "punctuator"
}

func (punctuatorLexer) accept(s scanner) bool {
	ch, _, err := s.peek()
	if err != nil {
		return false
	}
	return isPunctuator(string(ch))
}

func isPunctuator(ch string) bool {
	return puncs[ch]
}

func (p punctuatorLexer) lex(s scanner, ctx *context) (*token, error) {
	var start, end position
	if ctx.lastToken != nil {
		start, end = ctx.lastToken.End, start
	}
	nx, w, err := s.next()
	if err != nil {
		return nil, err
	}
	chrs := string(nx)
	end.Column += w
	for p.accept(s) {
		nxt, w, err := s.next()
		if err != nil {
			return nil, err
		}
		end.Column += w
		chrs += string(nxt)
	}
	tk := &token{Start: start, Text: chrs, End: end}
	switch chrs {
	case "{":
		tk.Kind = LBRACE
		return tk, nil
	case "}":
		tk.Kind = RBRACE
		return tk, nil
	case "(":
		tk.Kind = LPAREN
		return tk, nil
	case ")":
		tk.Kind = RPAREN
		return tk, nil
	case "[":
		tk.Kind = LBRACK
		return tk, nil
	case "]":
		tk.Kind = RBRACK
		return tk, nil
	case ".":
		tk.Kind = PERIOD
		return tk, nil
	case "...":
		tk.Kind = ELLIPSIS
		return tk, nil
	case ";":
		tk.Kind = SEMICOLON
		return tk, nil
	case ",":
		tk.Kind = COMMA
		return tk, nil
	case "<":
		tk.Kind = LSS
		return tk, nil
	case ">":
		tk.Kind = GTR
		return tk, nil
	case "<=":
		tk.Kind = LEQ
		return tk, nil
	case ">=":
		tk.Kind = GEQ
		return tk, nil
	case "==":
		tk.Kind = EQL
		return tk, nil
	case "!=":
		tk.Kind = NEQ
		return tk, nil
	case "===":
		tk.Kind = SEQL
		return tk, nil
	case "!==":
		tk.Kind = SNEQ
		return tk, nil
	case "+":
		tk.Kind = ADD
		return tk, nil
	case "-":
		tk.Kind = SUB
		return tk, nil
	case "*":
		tk.Kind = MUL
		return tk, nil
	case "%":
		tk.Kind = REM
		return tk, nil
	case "**":
		tk.Kind = EXP
		return tk, nil
	case "++":
		tk.Kind = INC
		return tk, nil
	case "--":
		tk.Kind = DEC
		return tk, nil
	case "<<":
		tk.Kind = SHL
		return tk, nil
	case ">>":
		tk.Kind = SHR
		return tk, nil
	case ">>>":
		tk.Kind = USHR
		return tk, nil
	case "&":
		tk.Kind = AND
		return tk, nil
	case "|":
		tk.Kind = OR
		return tk, nil
	case "^":
		tk.Kind = XOR
		return tk, nil
	case "!":
		tk.Kind = NOT
		return tk, nil
	case "~":
		tk.Kind = TILDE
		return tk, nil
	case "&&":
		tk.Kind = LAND
		return tk, nil
	case "||":
		tk.Kind = LOR
		return tk, nil
	case "?":
		tk.Kind = QN
		return tk, nil
	case ":":
		tk.Kind = COLON
		return tk, nil
	case "=":
		tk.Kind = ASSIGN
		return tk, nil
	case "+=":
		tk.Kind = AddAssign
		return tk, nil
	case "-=":
		tk.Kind = SubAssign
		return tk, nil
	case "*=":
		tk.Kind = MulAssign
		return tk, nil
	case "%=":
		tk.Kind = RemAssign
		return tk, nil
	case "**=":
		tk.Kind = ExpAssign
		return tk, nil
	case "<<=":
		tk.Kind = SHLAssign
		return tk, nil
	case ">>=":
		tk.Kind = SHRAssign
		return tk, nil
	case ">>>=":
		tk.Kind = SHRAssign
		return tk, nil
	case "&=":
		tk.Kind = AndAssign
		return tk, nil
	case "|=":
		tk.Kind = OrAssign
		return tk, nil
	case "^=":
		tk.Kind = XorAssign
		return tk, nil
	case "=>":
		tk.Kind = ARROW
		return tk, nil
	case "/":
		tk.Kind = QUO
		return tk, nil
	case "/=":
		tk.Kind = QuoAssign
		return tk, nil
	default:
		return nil, fmt.Errorf(unexpectedTkn, p.name(), end)
	}
}
