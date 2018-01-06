package lexer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

const unexpectedTkn = `%s : unexpected token at %v`

type kind uint

type token struct {
	Text  string
	Kind  kind   `json:"-"`
	RKind string `json:"Kind"`
	Start position
	End   position
}

func newToken(start position) *token {
	return &token{Start: start}
}

func (t *token) AddRune(ch rune) {
	t.AddString(string(ch))
}
func (t *token) AddString(txt string) {
	t.Text += txt
	t.End.Column += len(txt)
}

func printToken(tk *token) []byte {
	tk.RKind = tk.Kind.String()
	b, err := json.MarshalIndent(tk, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func printTokens(tks []*token) ([]byte, error) {
	for i := 0; i < len(tks); i++ {
		tk := tks[i]
		tk.RKind = tk.Kind.String()
	}
	b, err := json.MarshalIndent(tks, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	return b, nil
}

func decodeToken(b []byte) (*token, error) {
	t := &token{}
	err := json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}
	t.Kind = getKind(t.RKind)
	return t, nil
}

// lexical token types
const (
	ILLEGAL kind = iota
	EOF

	SingleLineComment
	MultiLineComment

	LF //LINE FEED
	CR //CARRIAGE RETURN
	LS //LINE SEPARATOR
	PS //PARAGRAPH SEPARATOR

	TAB    // CHARACTER TABULATION
	VT     //LINE TABULATION
	FF     //FORM FEED (FF)
	SP     //SPACE
	NBSP   //NO-BREAK SPACE
	ZWNBSP // ZERO WIDTH NO-BREAK SPACE
	USP    //Any other Unicode “Separator, space” code poin

	IdentifierName

	punctuator
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND    // &
	OR     // |
	XOR    // ^
	SHL    // <<
	SHR    // >>
	USHR   // >>>
	AndNot // &^

	AddAssign // +=
	SubAssign // -=
	MulAssign // *=
	QuoAssign // /=
	RemAssign // %=
	ExpAssign //**=

	AndAssign    // &=
	OrAssign     // |=
	XorAssign    // ^=
	SHLAssign    // <<=
	SHRAssign    // >>=
	USHRAssign   // >>>=
	AndNotAssign // &^=

	LAND // &&
	LOR  // ||
	INC  // ++
	DEC  // --
	EXP  //**

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !
	SEQL   // ===

	NEQ      // !=
	SNEQ     //!===
	QUOEQ    // /=
	LEQ      // <=
	GEQ      // >=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	QN    // ?
	TILDE //~
	ARROW // =>

	Template

	NULL  // null
	TRUE  // true
	FALSE //false

	INT //integer
	BINARY
	OCTAL
	FLOAT
	HEX
	STRING
)

var kindMap = map[kind]string{
	ILLEGAL:           "ILLEGAL",
	SingleLineComment: "SINGLE_LINE_COMMENT",
	MultiLineComment:  "MULTI_LINE_COMMENT",
	EOF:               "EOF",
	LF:                "LINE_FEED",
	CR:                "CARRIAGE_RETURN",
	LS:                "LINE_SEPARATOR",
	PS:                "PARAGRAPH_SEPARATOR",
	TAB:               "CHARACTER_TABULATION",
	VT:                "LINE_TABULATION",
	FF:                "FORM_FEED",
	SP:                "SPACE",
	NBSP:              "NO_BREAK_SPACE",
	ZWNBSP:            "ZERO_WIDTH_NO_BREAK_SPACE",
	USP:               "OTHER_SPACE",
	IdentifierName:    "IDENTIFIER_NAME",
	Template:          "TEMPLATE",
	ADD:               "ADD",
	SUB:               "SUB",
	MUL:               "MULTIPLY",
	QUO:               "QUOTIENT",
	REM:               "REMAINDER",
	AND:               "AND",
	OR:                "OR",
	XOR:               "XOR",
	SHL:               "LEFT_SHIFT",
	SHR:               "RIGHT_SHIFT",
	USHR:              "UNSIGNED_RIGHT_SHIFT",
	AndNot:            "AND_NOT",
	AddAssign:         "ADD_ASSIGN",
	SubAssign:         "SUB_ASSING",
	MulAssign:         "MUL_ASSIGN",
	QuoAssign:         "QUO_ASSIGN",
	RemAssign:         "REM_ASSIGN",
	ExpAssign:         "EXPONENT_ASSIGN",
	AndAssign:         "AND_ASSIGN",
	OrAssign:          "OR_ASSIGN",
	XorAssign:         "XOR_ASSIGN",
	SHLAssign:         "LEFT_SHIFT_ASSIGN",
	SHRAssign:         "RIGHT_SHIFT_ASSIGN",
	USHRAssign:        "UNSIGNED_RIGHT_SHIFT_ASSIGN",
	AndNotAssign:      "AND_NOT_ASSIGN",
	LAND:              "LOGICAL_AND",
	LOR:               "LOGICAL_OR",
	INC:               "INCREMENT",
	DEC:               "DECREMENT",
	EXP:               "EXPONENT",
	EQL:               "EQUAL",
	LSS:               "LESS_THAN",
	GTR:               "GREATER_THAN",
	ASSIGN:            "ASSIGN",
	NOT:               "NOT",
	SEQL:              "STRICT_EQUAL",
	NEQ:               "NOT_EQUAL",
	SNEQ:              "STRICT_NOT_EQUAL",
	QUOEQ:             "QUOTIENT_ASSIGN",
	LEQ:               "LESS_THAN_OR_EQUAL",
	GEQ:               "GREATER_THAN_OR_EQUAL",
	ELLIPSIS:          "ELLIPSIS",
	LPAREN:            "LEFT_PAREN",
	LBRACK:            "LEFT_BRACKET",
	LBRACE:            "LEFT_BRACE",
	COMMA:             "COMMA",
	PERIOD:            "PERIOD",
	RPAREN:            "RIGHT_PAREN",
	RBRACK:            "RIGHT_BRACKET",
	RBRACE:            "RIGHT_BRACE",
	SEMICOLON:         "SEMICOLON",
	COLON:             "COLON",
	QN:                "QUESTION_MARK",
	TILDE:             "TILDE",
	ARROW:             "ARROW",
	NULL:              "NULL",
	TRUE:              "TRUE",
	FALSE:             "FALSE",
	INT:               "INT",
	BINARY:            "BINARY",
	OCTAL:             "OCTAL",
	FLOAT:             "FLOAT",
	HEX:               "HEX",
	STRING:            "STRING",
}

var reverseKindMap map[string]kind

func init() {
	reverseKindMap = make(map[string]kind)
	for k, v := range kindMap {
		reverseKindMap[v] = k
	}
}

func (k kind) String() string {
	return kindMap[k]
}

func (k kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func getKind(k string) kind {
	return reverseKindMap[k]
}

// scanner is an interface for reading one token at a time from UTF text.
type scanner interface {
	Next() (rune, int, error)
	Peek() (rune, int, error)
	PeekAt(n int) (rune, int, error)
	Rewind() error
}

type bufioScanner struct {
	src *bufio.Reader
}

func newBufioScanner(r io.Reader) *bufioScanner {
	return &bufioScanner{src: bufio.NewReader(r)}
}
func (b *bufioScanner) Next() (rune, int, error) {
	return b.src.ReadRune()
}

func (b *bufioScanner) Peek() (ch rune, size int, err error) {
	return b.PeekAt(1)
}

// reads the nth rune without advancing the reader
func (b *bufioScanner) PeekAt(n int) (ch rune, size int, err error) {
	bv, err := b.peekChunck(n)
	width := 0
	for i := 0; i < n; i++ {
		ch, size = utf8.DecodeRune(bv[width:])
		width += size
	}
	return
}

func (b *bufioScanner) Rewind() error {
	return b.src.UnreadRune()
}

func (b *bufioScanner) peekChunck(n int) ([]byte, error) {
	return b.src.Peek(n)
}

type context struct {
	lexers    map[string]lexMe
	lastToken *token
}

type position struct {
	Line   int
	Column int
}

func (p position) String() string {
	return fmt.Sprintf("line %d: column %d", p.Line, p.Column)
}

type lexMe interface {
	Name() string
	Accept(scanner) bool
	Lex(scanner, *context) (*token, error)
}

// make sure all lexers implement lexMe interface
var (
	_ lexMe = singleLineCommentLexer{}
	_ lexMe = multiLineCommentLexer{}
	_ lexMe = lineTerminatorLexer{}
	_ lexMe = identifierNameLexer{}
	_ lexMe = punctuatorLexer{}
	_ lexMe = boolLexer{}
	_ lexMe = nullLexer{}
	_ lexMe = numeralLexer{}
	_ lexMe = stringLexer{}
)

// defaultLexMe returns a list of all available lexers.
func defaultLexMe() []lexMe {
	return []lexMe{
		singleLineCommentLexer{},
		multiLineCommentLexer{},
		lineTerminatorLexer{},
		identifierNameLexer{},
		punctuatorLexer{},
		boolLexer{},
		nullLexer{},
		numeralLexer{},
		stringLexer{},
	}
}

func lex(src io.Reader, lexmes ...lexMe) ([]*token, error) {
	s := &bufioScanner{bufio.NewReader(src)}
	ctx := &context{lexers: make(map[string]lexMe)}
	for _, v := range lexmes {
		ctx.lexers[v.Name()] = v
	}
	nextLexer := func() lexMe {
		for i := 0; i < len(lexmes); i++ {
			if lexmes[i].Accept(s) {
				return lexmes[i]
			}
		}
		return nil
	}
	var tokens []*token
	for {
		v := nextLexer()
		if v == nil {
			break
		}
		tk, err := v.Lex(s, ctx)
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, tk)
		ctx.lastToken = tk
	}
	return tokens, nil
}

// # Derived Property: ID_Start
// #  Characters that can start an identifier.
// #  Generated from:
// #      Lu + Ll + Lt + Lm + Lo + Nl
// #    + Other_ID_Start
// #    - Pattern_Syntax
// #    - Pattern_White_Space
// http://unicode.org/reports/tr44/#Simple_Derived
func isUnicodeIDStart(ch rune) bool {
	if unicode.In(ch, unicode.Lu, unicode.Ll,
		unicode.Lt, unicode.Lm, unicode.Lo,
		unicode.Nl, unicode.Other_ID_Start) {
		return !unicode.In(ch, unicode.Pattern_Syntax,
			unicode.Pattern_White_Space)
	}
	return false
}

func isUnicodeIDContinue(ch rune) bool {
	if isUnicodeIDStart(ch) && unicode.In(ch, unicode.Mn, unicode.Mc,
		unicode.Nd, unicode.Pc, unicode.Other_ID_Continue) {
		return !unicode.In(ch, unicode.Pattern_Syntax,
			unicode.Pattern_White_Space)
	}
	return false
}

func isTokenSep(ch rune) bool {
	return isWhiteSpace(ch) || isLineTerminator(ch)
}

func isTokenSepNoTerm(ch rune) bool {
	return isWhiteSpace(ch)
}
