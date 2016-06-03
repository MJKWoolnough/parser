package parser

import (
	"errors"
	"io"
	"strings"
)

// TokenType represents the type of token being read.
//
// Negative values are reserved for this package.
type TokenType int

// Constants TokenError (-2) and TokenDone (-1)
const (
	TokenDone TokenType = -1 - iota
	TokenError
)

// Token represents data parsed from the stream.
type Token struct {
	Type TokenType
	Data string
}

// TokenFunc is the type that the worker funcs implement in order to be used by
// the tokeniser.
type TokenFunc func(*Tokeniser) (Token, TokenFunc)

type tokeniser interface {
	backup()
	get() string
	length() int
	next() rune
}

// Tokeniser is a state machine to generate tokens from an input
type Tokeniser struct {
	tokeniser
	Err   error
	state TokenFunc
}

// GetToken runs the state machine and retrieves a single token and possible an
// error
func (t *Tokeniser) GetToken() (Token, error) {
	tk := t.get()
	if tk.Type == TokenError {
		return tk, t.Err
	}
	return tk, nil
}

// TokeniserState allows the internal state of the Tokeniser to be set
func (t *Tokeniser) TokeniserState(tf TokenFunc) {
	t.state = tf
}

func (t *Tokeniser) get() Token {
	if t.Err == io.EOF {
		return Token{
			Type: TokenDone,
			Data: "",
		}
	}
	if t.state == nil {
		t.Err = ErrNoState
		t.state = (*Tokeniser).Error
	}
	var tk Token
	tk, t.state = t.state(t)
	if tk.Type == TokenError && t.Err == io.EOF {
		t.Err = io.ErrUnexpectedEOF
	}
	return tk
}

// Accept returns true if the next character to be read is contained within the
// given string.
//
// Upon true, it advances the read position, otherwise the position remains the
// same.
func (t *Tokeniser) Accept(chars string) bool {
	if strings.IndexRune(chars, t.next()) < 0 {
		t.backup()
		return false
	}
	return true
}

// Peek returns the next rune without advancing the read position.
func (t *Tokeniser) Peek() rune {
	r := t.next()
	t.backup()
	return r
}

// Get returns a string of everything that has been read so far and resets
// the string for the next round of parsing.
func (t *Tokeniser) Get() string {
	return t.tokeniser.get()
}

// Len returns the number of bytes that has been read since the last Get.
func (t *Tokeniser) Len() int {
	return t.length()
}

// AcceptRun reads from the string as long as the read character is in the
// given string.
//
// Returns the rune that stopped the run.
func (t *Tokeniser) AcceptRun(chars string) rune {
	for {
		if c := t.next(); strings.IndexRune(chars, c) < 0 {
			t.backup()
			return c
		}
	}
}

// Except returns true if the next character to be read is not contained within
// the given string.
// Upon true, it advances the read position, otherwise the position remains the
// same.
func (t *Tokeniser) Except(chars string) bool {
	if r := t.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
		t.backup()
		return false
	}
	return true
}

// ExceptRun reads from the string as long as the read character is not in the
// given string.
//
// Returns the rune that stopped the run.
func (t *Tokeniser) ExceptRun(chars string) rune {
	for {
		if r := t.next(); r == -1 || strings.IndexRune(chars, r) >= 0 {
			t.backup()
			return r
		}
	}
}

// Done is a TokenFunc that is used to indicate that there are no more tokens to
// parse.
func (t *Tokeniser) Done() (Token, TokenFunc) {
	t.Err = io.EOF
	return Token{
		Type: TokenDone,
		Data: "",
	}, (*Tokeniser).Done
}

// Error represents an error state for the parser.
//
// The error value should be set in Tokeniser.Err and then this func should be
// called.
func (t *Tokeniser) Error() (Token, TokenFunc) {
	return Token{
		Type: TokenError,
		Data: t.Err.Error(),
	}, (*Tokeniser).Error
}

// Errors
var (
	ErrNoState = errors.New("no state")
)
