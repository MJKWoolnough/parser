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

func (p *strParser) reset() {
	p.pos = 0
	p.width = 0
}

func (p *strParser) sub() tokeniser {
	return &sub{
		tokeniser: p,
		tState:    len(p.str),
		start:     p.pos,
	}
}

func (p *strParser) slice(state, start int) (string, int) {
	if len(p.str) != state || start > p.pos {
		return "", -1
	}

	return p.str[start:p.pos], p.pos
}

type strState struct {
	s          *strParser
	stateID    int
	pos, width int
}

func (p *strParser) state() State {
	return &strState{
		s:       p,
		stateID: len(p.str),
		pos:     p.pos,
		width:   p.width,
	}
}

func (s *strState) Reset() bool {
	if len(s.s.str) != s.stateID {
		return false
	}

	s.s.pos = s.pos
	s.s.width = s.width

	return true
}
