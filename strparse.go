// Package parser is a simple helper package for parsing strings, byte slices and Readers
package parser

import "unicode/utf8"

type strParser struct {
	str        string
	pos, width int
}

func (p *strParser) next() rune {
	if p.pos == len(p.str) {
		p.width = 0
		return -1
	}
	r, s := utf8.DecodeRuneInString(p.str[p.pos:])
	if r == utf8.RuneError && s == 1 {
		r = rune(p.str[p.pos])
	}
	p.pos += s
	p.width = s
	return r
}

func (p *strParser) backup() {
	if p.width > 0 {
		p.pos -= p.width
		p.width = 0
	}
}

func (p *strParser) get() string {
	s := p.str[:p.pos]
	p.str = p.str[p.pos:]
	p.pos = 0
	p.width = 0
	return s
}

func (p *strParser) length() int {
	return p.pos
}
