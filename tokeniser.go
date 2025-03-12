package parser

import (
	"errors"
	"io"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TokenType represents the type of token being read.
//
// Negative values are reserved for this package.
type TokenType int

// Constants TokenError (-2) and TokenDone (-1).
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

// State represents a position in the byte stream of the Tokeniser.
type State interface {
	// Reset returns the byte stream to the position it was in when this
	// object was created.
	//
	// Only valid until Tokeniser.Get is called.
	Reset() bool
}

type tokeniser interface {
	backup()
	get() string
	length() int
	next() rune
	reset()
	state() State
	sub() tokeniser
	slice(int, int) (string, int)
}

// Tokeniser is a state machine to generate tokens from an input.
type Tokeniser struct {
	tokeniser
	Err   error
	state TokenFunc
}

// GetToken runs the state machine and retrieves a single token and possible an
// error.
func (t *Tokeniser) GetToken() (Token, error) {
	tk := t.get()

	if tk.Type == TokenError {
		return tk, t.Err
	}

	return tk, nil
}

// Iter yields each token as it's returned, stopping after yielding a TokenDone
// or TokenError Token.
func (t *Tokeniser) Iter(yield func(Token) bool) {
	for {
		if tk := t.get(); !yield(tk) || tk.Type == TokenDone || tk.Type == TokenError {
			break
		}
	}
}

// GetError returns any error that has been generated by the Tokeniser.
func (t *Tokeniser) GetError() error {
	return t.Err
}

// TokeniserState allows the internal state of the Tokeniser to be set.
func (t *Tokeniser) TokeniserState(tf TokenFunc) {
	t.state = tf
}

func (t *Tokeniser) get() Token {
	if errors.Is(t.Err, io.EOF) {
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

	if tk.Type == TokenError && errors.Is(t.Err, io.EOF) {
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
	if !strings.ContainsRune(chars, t.next()) {
		t.backup()

		return false
	}

	return true
}

// Next returns the next rune and advances the read position.
func (t *Tokeniser) Next() rune {
	return t.next()
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
		if c := t.next(); !strings.ContainsRune(chars, c) {
			t.backup()

			return c
		}
	}
}

// AcceptString attempts to accept each character from the given string, in
// order, returning the number of characters accepted before a failure.
func (t *Tokeniser) AcceptString(str string, caseInsensitive bool) int {
	for n, r := range str {
		if p := t.Peek(); p < 0 || !runeComparison(p, r, caseInsensitive) {
			return n
		}

		t.Next()
	}

	return len(str)
}

func runeComparison(a, b rune, caseInsensitive bool) bool {
	if caseInsensitive {
		al := unicode.SimpleFold(a)
		bl := unicode.SimpleFold(b)

		return a == b || a == bl || bl == a || al == bl
	}

	return a == b
}

// AcceptWord attempts to parse one of the words (string of characters)
// provided in the slice.
//
// Returns the longest word parsed, or empty string if no words matched.
func (t *Tokeniser) AcceptWord(words []string, caseInsensitive bool) string {
	words = slices.Clone(words)

	return t.acceptWord(words, caseInsensitive)
}

func (t *Tokeniser) acceptWord(words []string, caseInsensitive bool) string {
	s := t.State()

	var sb strings.Builder

	for len(words) > 0 {
		char := t.Next()

		sb.WriteRune(char)

		if char < 0 {
			break
		}

		var found bool

		newWords := words[:0]

		for _, word := range words {
			if len(word) > 0 {
				r, s := utf8.DecodeRuneInString(word)
				if r == utf8.RuneError && s == 1 {
					r = rune(word[0])
				}

				if runeComparison(char, r, caseInsensitive) {
					word = word[s:]
					found = found || word == ""
					newWords = append(newWords, word)
				}
			}
		}

		words = newWords

		if found {
			if len(words) > 0 {
				sb.WriteString(t.acceptWord(words, caseInsensitive))
			}

			return sb.String()
		}
	}

	s.Reset()

	return ""
}

// Except returns true if the next character to be read is not contained within
// the given string.
// Upon true, it advances the read position, otherwise the position remains the
// same.
func (t *Tokeniser) Except(chars string) bool {
	if r := t.next(); r == -1 || strings.ContainsRune(chars, r) {
		t.backup()

		return false
	}

	return true
}

// Reset restores the state to after the last Get() call (or init, it Get() has
// not been called).
func (t *Tokeniser) Reset() {
	t.reset()
}

// Retrieve the current Tokeniser state that allows you to reset to that point.
// State is only valid until next 'Get' call.
func (t *Tokeniser) State() State {
	return t.tokeniser.state()
}

// SubTokeniser create a new Tokeniser that uses this existing tokeniser as its
// source.
//
// This allows the sub-tokenisers Get method to be called without calling it on
// its parent.
func (t *Tokeniser) SubTokeniser() *Tokeniser {
	return &Tokeniser{
		tokeniser: t.tokeniser.sub(),
	}
}

// ExceptRun reads from the string as long as the read character is not in the
// given string.
//
// Returns the rune that stopped the run.
func (t *Tokeniser) ExceptRun(chars string) rune {
	for {
		if r := t.next(); r == -1 || strings.ContainsRune(chars, r) {
			t.backup()

			return r
		}
	}
}

// Return simplifies the returning from TokenFns, taking a TokenType and a next
// TokenFn, default to Done.
//
// The returned token is of the type specified with the data set to the output
// of t.Get().
func (t *Tokeniser) Return(typ TokenType, fn TokenFunc) (Token, TokenFunc) {
	if fn == nil {
		fn = (*Tokeniser).Done
	}

	return Token{
		Type: typ,
		Data: t.Get(),
	}, fn
}

// ReturnError simplifies the handling of errors, setting the error and calling
// Tokeniser.Error().
func (t *Tokeniser) ReturnError(err error) (Token, TokenFunc) {
	t.Err = err

	return t.Error()
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
	if t.Err == nil {
		t.Err = ErrUnknownError
	}

	return Token{
		Type: TokenError,
		Data: t.Err.Error(),
	}, (*Tokeniser).Error
}

type sub struct {
	tokeniser
	tState, start int
}

func (s *sub) get() string {
	if s.start < 0 {
		return ""
	}

	var str string

	str, s.start = s.slice(s.tState, s.start)

	return str
}

// Errors.
var (
	ErrNoState      = errors.New("no state")
	ErrUnknownError = errors.New("unknown error")
)
