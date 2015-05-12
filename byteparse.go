// Package parser is a simple helper package for parsing strings, byte slices and Readers
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
