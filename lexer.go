package goes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func printToken(tk *token) []byte {
	tk.RKind = tk.Kind.String()
	b, err := json.MarshalIndent(tk, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	return b
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
	Unknown kind = iota
	eof
	comment
	SingleLineComment
	MultiLineComment

	lineTerminator
	LF //LINE FEED
	CR //CARRIAGE RETURN
	LS //LINE SEPARATOR
	PS //PARAGRAPH SEPARATOR

	whiteSpace
	TAB    // CHARACTER TABULATION
	VT     //LINE TABULATION
	FF     //FORM FEED (FF)
	SP     //SPACE
	NBSP   //NO-BREAK SPACE
	ZWNBSP // ZERO WIDTH NO-BREAK SPACE
	USP    //Any other Unicode “Separator, space” code poin
)

func (k kind) String() string {
	switch k {
	case SingleLineComment:
		return "SINGLE_LINE_COMMENT"
	case MultiLineComment:
		return "MULTI_LINE_COMMENT"
	case eof:
		return "EOF"
	case LF:
		return "LINE_FEED"
	case CR:
		return "CARRIAGE_RETURN"
	case LS:
		return "LINE_SEPARATOR"
	case PS:
		return "PARAGRAPH_SEPARATOR"
	case TAB:
		return "CHARACTER_TABULATION"
	case VT:
		return "LINE_TABULATION"
	case FF:
		return "FORM_FEED"
	case SP:
		return "SPACE"
	case NBSP:
		return "NO_BREAK_SPACE"
	case ZWNBSP:
		return "ZERO_WIDTH_NO_BREAK_SPACE"
	case USP:
		return "OTHER_SPACE"
	default:
		return "UNKOWN"
	}
}

func (k kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func getKind(k string) kind {
	switch k {
	case "SINGLE_LINE_COMMENT":
		return SingleLineComment
	case "MULTI_LINE_COMMENT":
		return MultiLineComment
	case "EOF":
		return eof
	case "LINE_FEED":
		return LF
	case "CARRIAGE_RETURN":
		return CR
	case "LINE_SEPARATOR":
		return LS
	case "PARAGRAPH_SEPARATOR":
		return PS
	case "CHARACTER_TABULATION":
		return TAB
	case "LINE_TABULATION":
		return VT
	case "FORM_FEED":
		return FF
	case "SPACE":
		return SP
	case "NO_BREAK_SPACE":
		return NBSP
	case "ZERO_WIDTH_NO_BREAK_SPACE":
		return ZWNBSP
	case "OTHER_SPACE":
		return USP
	default:
		return Unknown
	}
}

type scanner interface {
	next() (rune, int, error)
	peek() (rune, int, error)
	peekAt(n int) (rune, int, error)
}

type bufioScanner struct {
	src *bufio.Reader
}

func newBufioScanner(r io.Reader) *bufioScanner {
	return &bufioScanner{src: bufio.NewReader(r)}
}
func (b *bufioScanner) next() (rune, int, error) {
	return b.src.ReadRune()
}

func (b *bufioScanner) peek() (ch rune, size int, err error) {
	defer func() {
		err = b.src.UnreadByte()
	}()
	return b.src.ReadRune()
}

// reads the nth rune without advancing the reader
func (b *bufioScanner) peekAt(n int) (ch rune, size int, err error) {
	max := n * utf8.UTFMax
	bv, err := b.src.Peek(max)
	if err != nil {
		if err == io.EOF {
			// try reading a small chunk assuming the unicode chars are of size
			// 2
			bv, err = b.src.Peek(2 * n)
			if err != nil {
				if err == io.EOF {
					// try reading a small chunk assuming the unicode chars are
					// of size n
					bv, err = b.src.Peek(n)
					if err != nil {
						return 0, 0, err
					}
				}
			}
		}
		return 0, 0, err
	}
	width := 0
	for i := 0; i < n; i++ {
		ch, size = utf8.DecodeRune(bv[width:])
		width += size
	}
	return
}

func (b *bufioScanner) rewind() error {
	return b.src.UnreadRune()
}

func (b *bufioScanner) rewindN(n int) error {
	for i := 0; i < n; i++ {
		if err := b.rewind(); err != nil {
			return err
		}
	}
	return nil
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
	name() string
	accept(scanner) bool
	lex(scanner, *context) (*token, error)
}

func lex(src io.Reader, lexmes ...lexMe) ([]*token, error) {
	s := &bufioScanner{bufio.NewReader(src)}
	ctx := &context{lexers: make(map[string]lexMe)}
	for _, v := range lexmes {
		ctx.lexers[v.name()] = v
	}
	nextLexer := func() lexMe {
		for i := 0; i < len(lexmes); i++ {
			if lexmes[i].accept(s) {
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
		tk, err := v.lex(s, ctx)
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, tk)
		ctx.lastToken = tk
	}
	return tokens, nil
}
