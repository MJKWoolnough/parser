package parser

import (
	"io"
	"unicode/utf8"
)

type runeSourceParser struct {
	source   io.RuneReader
	buf      []rune
	pos      int
	stateNum int
}

func (r *runeSourceParser) next() rune {
	if r.pos < len(r.buf) {
		ru := r.buf[r.pos]
		r.pos++

		return ru
	}

	ru, _, err := r.source.ReadRune()
	if err != nil {
		ru = -1
	}

	r.buf = append(r.buf, ru)
	r.pos++

	return ru
}

func (r *runeSourceParser) backup() {
	if r.pos > 0 {
		r.pos--
	}
}

func (r *runeSourceParser) get() string {
	s := string(r.buf[:r.pos])
	r.buf = r.buf[r.pos:]
	r.pos = 0
	r.stateNum++

	return s
}

func (r *runeSourceParser) length() int {
	var l int

	for _, r := range r.buf[:r.pos] {
		if r != -1 {
			l += utf8.RuneLen(r)
		}
	}

	return l
}

func (r *runeSourceParser) reset() {
	r.pos = 0
}

func (r *runeSourceParser) sub() tokeniser {
	return &sub{
		tokeniser: r,
		tState:    r.stateNum,
		start:     r.pos,
	}
}

func (r *runeSourceParser) slice(state, start int) (string, int) {
	if r.stateNum != state || start > r.pos {
		return "", -1
	}

	return string(r.buf[start:r.pos]), r.pos
}

type runeState struct {
	r        *runeSourceParser
	stateNum int
	pos      int
}

func (r *runeSourceParser) state() State {
	return &runeState{
		r:        r,
		stateNum: r.stateNum,
		pos:      r.pos,
	}
}

func (r *runeState) Reset() bool {
	if r.r.stateNum != r.stateNum {
		return false
	}

	r.r.pos = r.pos

	return true
}
