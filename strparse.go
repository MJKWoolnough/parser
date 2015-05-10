// Package strparse is a simple helper package for parsing strings
package strparse

import (
	"strings"
	"unicode/utf8"
)

// Parser is a helper with aids with the parsing of formatted strings.
type Parser struct {
	Str        string
	pos, width int
}

// New returns a new Parser type containg the given string.
func New(s string) *Parser {
	return &Parser{Str: s}
}

func (p *Parser) next() rune {
	if p.pos == len(p.Str) {
		p.width = 0
		return -1
	}
	r, s := utf8.DecodeRuneInString(p.Str[p.pos:])
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

// Peek returns the next rune without advancing the read position.
func (p *Parser) Peek() rune {
	r := p.next()
	p.backup()
	return r
}

// Get returns a string of everything that has been read so far and resets the
// string for the next round of parsing.
func (p *Parser) Get() string {
	s := p.Str[:p.pos]
	p.Clear()
	return s
}

// Len returns the current length of the read string.
func (p *Parser) Len() int {
	return p.pos
}

// Left returns how much of the string is left in the Parser. This includes
// everything read since the last Get.
func (p *Parser) Left() int {
	return len(p.Str)
}

func (p *Parser) Clear() {
	p.Str = p.Str[p.pos:]
	p.pos = 0
	p.width = 0
}

// Accept returns true if the next character to be read is contained within the
// given string.
// Upon true, it advances the read position, otherwise the position remains the
// same.
func (p *Parser) Accept(chars string) bool {
	if strings.IndexRune(chars, p.next()) < 0 {
		p.backup()
		return false
	}
	return true
}

// AcceptRun reads from the string as long as the read character is in the
// given string.
func (p *Parser) AcceptRun(chars string) {
	for {
		if strings.IndexRune(chars, p.next()) < 0 {
			p.backup()
			break
		}
	}
}

// Except returns true if the next character to be read is not contained within
// the given string.
// Upon true, it advances the read position, otherwise the position remains the
// same.
func (p *Parser) Except(chars string) bool {
	if r := p.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
		p.backup()
		return false
	}
	return true
}

// ExceptRun reads from the string as long as the read character is not in the
// given string.
func (p *Parser) ExceptRun(chars string) {
	for {
		if r := p.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
			p.backup()
			break
		}
	}
}
