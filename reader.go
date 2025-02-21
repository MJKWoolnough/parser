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
	if r.pos < len(r.buf) {
		ru := r.buf[r.pos]
		r.pos++

		return ru
	}

	ru, s, err := r.reader.ReadRune()
	if err != nil {
		ru = -1
	}

	if ru == utf8.RuneError && s == 1 {
		b, _ := r.reader.ReadByte()
		ru = rune(b)
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
	l := 0

	for _, r := range r.buf[:r.pos] {
		rl := utf8.RuneLen(r)
		if rl > 0 {
			l += rl
		} else {
			l++
		}
	}

	return l
}
