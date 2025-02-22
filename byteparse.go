package parser

import "unicode/utf8"

type byteParser struct {
	data       []byte
	pos, width int
}

func (p *byteParser) next() rune {
	if p.pos == len(p.data) {
		p.width = 0

		return -1
	}

	r, s := utf8.DecodeRune(p.data[p.pos:])
	if r == utf8.RuneError && s == 1 {
		r = rune(p.data[p.pos])
	}

	p.pos += s
	p.width = s

	return r
}

func (p *byteParser) backup() {
	if p.width > 0 {
		p.pos -= p.width
		p.width = 0
	}
}

func (p *byteParser) get() string {
	s := p.data[:p.pos]
	p.data = p.data[p.pos:]
	p.pos = 0
	p.width = 0

	return string(s)
}

func (p *byteParser) length() int {
	return p.pos
}

func (p *byteParser) reset() {
	p.pos = 0
	p.width = 0
}

func (p *byteParser) sub() tokeniser {
	return &sub{
		tokeniser: p,
		tState:    len(p.data),
		start:     p.pos,
	}
}

func (p *byteParser) slice(state, start int) (string, int) {
	if len(p.data) != state || start > p.pos {
		return "", -1
	}

	return string(p.data[start:p.pos]), p.pos
}

type byteState struct {
	b          *byteParser
	stateID    int
	pos, width int
}

func (p *byteParser) state() State {
	return &byteState{
		b:       p,
		stateID: len(p.data),
		pos:     p.pos,
		width:   p.width,
	}
}

func (b *byteState) Reset() bool {
	if len(b.b.data) != b.stateID {
		return false
	}

	b.b.pos = b.pos
	b.b.width = b.width

	return true
}
