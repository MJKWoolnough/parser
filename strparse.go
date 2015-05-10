// Package strparse is a simple helper package for parsing strings
package strparse

import (
	"strings"
	"unicode/utf8"
)

type Parser struct {
	Str        string
	pos, width int
}

func New(s string) *Parser {
	return &Parser{Str: s}
}

func (p *Parser) next() rune {
	if p.pos == len(p.Str) {
		p.width = 0
		return -1
	}
	r, s := utf8.DecodeRuneInString(p.Str)
	p.pos += s
	p.width = s
	return r
}

func (p *Parser) backup() {
	if p.width > 0 {
		p.pos -= p.width
		p.width = 0
	}
}

func (p *Parser) Peek() rune {
	r := p.next()
	p.backup()
	return r
}

func (p *Parser) Get() string {
	s := p.Str[:p.pos]
	p.Clear()
	return s
}

func (p *Parser) Len() int {
	return p.pos
}

func (p *Parser) Cap() int {
	return len(p.Str)
}

func (p *Parser) Clear() {
	p.Str = p.Str[p.pos:]
	p.pos = 0
	p.width = 0
}

func (p *Parser) Accept(chars string) bool {
	if strings.IndexRune(chars, p.next()) < 0 {
		p.backup()
		return false
	}
	return true
}

func (p *Parser) AcceptRun(chars string) {
	for {
		if strings.IndexRune(chars, p.next()) < 0 {
			p.backup()
			break
		}
	}
}

func (p *Parser) Except(chars string) bool {
	if r := p.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
		p.backup()
		return false
	}
	return true
}

func (p *Parser) ExceptRun(chars string) {
	for {
		if r := p.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
			p.backup()
			break
		}
	}
}
