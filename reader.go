package parser

import (
	"bufio"
	"unicode/utf8"
)

type readerParser struct {
	reader *bufio.Reader
	buf    []rune
	pos    int
}

func (r *readerParser) next() rune {
	ru, s, err := r.reader.ReadRune()
	if err != nil {
		return -1
	}

	if ru == utf8.RuneError && s == 1 {
		r.reader.UnreadRune()

		b, _ := r.reader.ReadByte()
		ru = rune(b)

		r.reader.ReadRune()
	}

	r.buf = append(r.buf, ru)
	r.pos++

	return ru
}

func (r *readerParser) backup() {
	if r.pos > 0 {
		r.pos--
	}
}

func (r *readerParser) get() string {
	s := string(r.buf[:r.pos])
	r.buf = r.buf[r.pos:]
	r.pos = 0

	return s
}

func (r *readerParser) length() int {
	return r.pos
}

func (r *readerParser) reset() {
	r.pos = 0
}
