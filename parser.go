// Package parser is a simple helper package for parsing strings, byte slices and io.Readers
package parser

import (
	"bufio"
	"errors"
	"io"
)

// TokenType represents the type of token being read.
//
// Negative values are reserved for this package.
type TokenType int

// Constants TokenError (-2) and TokenDone (-1)
const (
	TokenDone TokenType = -1 - iota
	TokenError
	TokenEmpty
)

// Token represents data parsed from the stream.
type Token struct {
	Type TokenType
	Data string
}

// StateFn is the type that the worker funcs implement in order to be used by
// the parser.
type StateFn func() (Token, StateFn)

type parser interface {
	backup()
	get() string
	length() int
	next() rune
}

// Parser is the wrapper type for the various different parsers.
type Parser struct {
	parser
	State       StateFn
	Parser      ParseFn
	Err         error
	peekedToken Token
	tokens      []Token
}

// NewStringParser returns a Parser which parses a string.
func NewStringParser(str string) Parser {
	return Parser{parser: &strParser{str: str}, peekedToken: Token{Type: TokenEmpty}}
}

// NewByteParser returns a Parser which parses a byte slice.
func NewByteParser(data []byte) Parser {
	return Parser{parser: &byteParser{data: data}, peekedToken: Token{Type: TokenEmpty}}
}

// NewReaderParser returns a Parser which parses a Reader.
func NewReaderParser(reader io.Reader) Parser {
	return Parser{parser: &readerParser{reader: bufio.NewReader(reader)}, peekedToken: Token{Type: TokenEmpty}}
}

// GetToken reads the next token in the stream, and returns the token and any
// error that occurred.
func (p *Parser) GetToken() (Token, error) {
	if p.peekedToken.Type != TokenEmpty {
		tk := p.peekedToken
		p.peekedToken.Type = TokenEmpty
		return tk, p.Err
	}
	if p.Err == io.EOF {
		return Token{
			Type: TokenDone,
			Data: "",
		}, io.EOF
	}
	if p.State == nil {
		p.Err = ErrNoState
		p.State = p.Error
	}
	var tk Token
	tk, p.State = p.State()
	if p.Err == io.EOF {
		if tk.Type == TokenError {
			p.Err = io.ErrUnexpectedEOF
		} else {
			return tk, nil
		}
	}
	return tk, p.Err
}

// BufferToken puts the given token in the Peek buffer.
func (p *Parser) BufferToken(t Token) {
	p.peekedToken = t
}

// Done is a StateFn that is used to indicate that there are no more tokens to
// parse.
func (p *Parser) Done() (Token, StateFn) {
	p.Err = io.EOF
	return Token{
		Type: TokenDone,
		Data: "",
	}, p.Done
}

// Error represents an error state for the parser.
//
// Should be called from other StateFn's that detect an error. The error value
// should be set to Parser.Err and then this func should be called.
func (p *Parser) Error() (Token, StateFn) {
	return Token{
		Type: TokenError,
		Data: p.Err.Error(),
	}, p.Error
}

// Errors
var (
	ErrNoState = errors.New("no state")
)
