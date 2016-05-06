// Package parser is a simple helper package for parsing strings, byte slices and io.Readers
package parser

import (
	"bufio"
	"io"
	"strings"
)

type parser interface {
	backup()
	get() string
	length() int
	next() rune
}

// Parser is the wrapper type for the various different parsers.
type Parser struct {
	parser
}

// NewStringParser returns a Parser which parses a string.
func NewStringParser(str string) Parser {
	return Parser{&strParser{str: str}}
}

// NewByteParser returns a Parser which parses a byte slice.
func NewByteParser(data []byte) Parser {
	return Parser{&byteParser{data: data}}
}

// NewReaderParser returns a Parser which parses a Reader.
func NewReaderParser(reader io.Reader) Parser {
	return Parser{&readerParser{reader: bufio.NewReader(reader)}}
}

// Peek returns the next rune without advancing the read position.
func (p Parser) Peek() rune {
	r := p.next()
	p.backup()
	return r
}

// Get returns a string of everything that has been read so far and resets the
// string for the next round of parsing.
func (p Parser) Get() string {
	return p.get()
}

// Len returns the number of bytes that has been read since the last Get.
func (p Parser) Len() int {
	return p.length()
}

// Accept returns true if the next character to be read is contained within the
// given string.
// Upon true, it advances the read position, otherwise the position remains the
// same.
func (p Parser) Accept(chars string) bool {
	if strings.IndexRune(chars, p.next()) < 0 {
		p.backup()
		return false
	}
	return true
}

// AcceptRun reads from the string as long as the read character is in the
// given string.
func (p Parser) AcceptRun(chars string) {
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
func (p Parser) Except(chars string) bool {
	if r := p.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
		p.backup()
		return false
	}
	return true
}

// ExceptRun reads from the string as long as the read character is not in the
// given string.
func (p Parser) ExceptRun(chars string) {
	for {
		if r := p.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
			p.backup()
			break
		}
	}
}
