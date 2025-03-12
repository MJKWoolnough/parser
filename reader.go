package parser

import (
	"bufio"
	"unicode/utf8"
)

type readerParser struct {
	reader   *bufio.Reader
	buf      []rune
	pos      int
	stateNum int
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
	} else if ru == utf8.RuneError && s == 1 {
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
	r.stateNum++

	return s
}

func (r *readerParser) length() int {
	var l int

	for _, r := range r.buf {
		s := utf8.RuneLen(r)
		l += s
	}

	return l
}

func (r *readerParser) reset() {
	r.pos = 0
}

func (r *readerParser) sub() tokeniser {
	return &sub{
		tokeniser: r,
		tState:    r.stateNum,
		start:     r.pos,
	}
}

func (r *readerParser) slice(state, start int) (string, int) {
	if r.stateNum != state || start > r.pos {
		return "", -1
	}

	return string(r.buf[start:r.pos]), r.pos
}

type readerState struct {
	r        *readerParser
	stateNum int
	pos      int
}

func (r *readerParser) state() State {
	return &readerState{
		r:        r,
		stateNum: r.stateNum,
		pos:      r.pos,
	}
}

func (r *readerState) Reset() bool {
	if r.r.stateNum != r.stateNum {
		return false
	}

	r.r.pos = r.pos

	return true
}
