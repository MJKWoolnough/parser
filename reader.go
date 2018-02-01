package parser

import (
	"bufio"
	"bytes"
	"unicode/utf8"
)

type readerParser struct {
	reader *bufio.Reader
	buf    bytes.Buffer
	width  int
}

func (r *readerParser) next() rune {
	ru, s, err := r.reader.ReadRune()
	if err != nil {
		r.width = 0
		return -1
	}
	if ru == utf8.RuneError && s == 1 {
		_ = r.reader.UnreadRune()
		b, _ := r.reader.ReadByte()
		_, _, _ = r.reader.ReadRune()
		ru = rune(b)
	}
	r.width = s
	r.buf.WriteRune(ru)
	return ru
}

func (r *readerParser) backup() {
	if r.width > 0 {
		r.buf.Truncate(r.buf.Len() - r.width)
		_ = r.reader.UnreadRune()
		r.width = 0
	}
}

func (r *readerParser) get() string {
	s := r.buf.String()
	r.buf.Reset()
	r.width = 0
	return s
}

func (r *readerParser) length() int {
	return r.buf.Len()
}
