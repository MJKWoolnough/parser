package parser

import (
	"bufio"
	"bytes"
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
	r.width = s
	r.buf.WriteRune(ru)
	return ru
}

func (r *readerParser) backup() {
	if r.width > 0 {
		r.buf.Truncate(r.buf.Len() - r.width)
		r.reader.UnreadRune()
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
