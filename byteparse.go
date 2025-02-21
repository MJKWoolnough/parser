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

type byteState struct {
	b     *byteParser
	state byteParser
}

func (p *byteParser) state() State {
	return &byteState{
		b:     p,
		state: *p,
	}
}

func (b *byteState) Reset() bool {
	if string(b.b.data) != string(b.state.data) {
		return false
	}

	b.b.pos = b.state.pos
	b.b.width = b.state.width

	return true
}
