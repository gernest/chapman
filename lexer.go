package goes

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

const unexpectedTkn = `%s : unexpected token at [post] expected %s found %s`

type kind uint

func (k kind) String() string {
	switch k {
	case comment:
		return "COMMENT"
	case lineTerminator:
		return "LINE_TERMINATOR"
	default:
		return "UNKOWN"
	}
}

type token struct {
	pos    int
	offset int
	text   string
	line   int
	kind   kind
}

func (t *token) String() string {
	return fmt.Sprintf(" <%s>", t.kind.String())
}

const (
	comment kind = iota
	eof

	lineTerminator
)

type scanner interface {
	next() (rune, int, error)
	peek() (rune, int, error)
	peekAt(n int) (rune, int, error)
}

type bufioScanner struct {
	src *bufio.Reader
}

func (b *bufioScanner) next() (rune, int, error) {
	return b.src.ReadRune()
}

func (b *bufioScanner) peek() (ch rune, size int, err error) {
	return b.peekAt(1)
}

// reads the nth rune without advancing the reader
func (b *bufioScanner) peekAt(n int) (ch rune, size int, err error) {
	max := n * utf8.UTFMax
	bv, err := b.src.Peek(max)
	if err != nil {
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
	lexers map[string]lexMe
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
	}
	return tokens, nil
}
