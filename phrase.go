package parser

import (
	"errors"
	"io"
)

// PhraseType represents the type of phrase being read.
//
// Negative values are reserved for this package.
type PhraseType int

// Constants PhraseError (-2) and PhraseDone (-1).
const (
	PhraseDone PhraseType = -1 - iota
	PhraseError
)

// PhraseFunc is the type that the worker types implement in order to be used
// by the Phraser.
type PhraseFunc func(*Parser) (Phrase, PhraseFunc)

// Phrase represents a collection of tokens that have meaning together.
type Phrase struct {
	Type PhraseType
	Data []Token
}

// Parser is a type used to get tokens or phrases (collection of token) from an
// input.
type Parser struct {
	Tokeniser
	state       PhraseFunc
	tokens      []Token
	peekedToken bool
}

// GetPhrase runs the state machine and retrieves a single Phrase and possibly
// an error.
func (p *Parser) GetPhrase() (Phrase, error) {
	if errors.Is(p.Err, io.EOF) {
		return Phrase{
			Type: PhraseDone,
			Data: make([]Token, 0),
		}, io.EOF
	}

	if p.state == nil {
		p.Err = ErrNoState
		p.state = (*Parser).Error
	}

	var ph Phrase

	ph, p.state = p.state(p)

	if ph.Type == PhraseError {
		if errors.Is(p.Err, io.EOF) {
			p.Err = io.ErrUnexpectedEOF
		}

		return ph, p.Err
	}

	return ph, nil
}

// GetToken runs the state machine and retrieves a single Token and possibly
// an error.
//
// If a Token has already been 'peek'ed, that token will be returned without
// running the state machine.
func (p *Parser) GetToken() (Token, error) {
	tk := p.get()

	if tk.Type == TokenError {
		return tk, p.Err
	}

	return tk, nil
}

// PhraserState allows the internal state of the Phraser to be set.
func (p *Parser) PhraserState(pf PhraseFunc) {
	p.state = pf
}

func (p *Parser) get() Token {
	if p.peekedToken {
		p.peekedToken = false

		return p.tokens[len(p.tokens)-1]
	} else if len(p.tokens) > 0 && p.tokens[len(p.tokens)-1].Type < 0 {
		return p.tokens[len(p.tokens)-1]
	}

	tk := p.Tokeniser.get()
	p.tokens = append(p.tokens, tk)

	return tk
}

func (p *Parser) backup() {
	p.peekedToken = true
}

// Accept will accept a token with one of the given types, returning true if
// one is read and false otherwise.
func (p *Parser) Accept(types ...TokenType) bool {
	tk := p.get()

	for _, t := range types {
		if tk.Type == t {
			return true
		}
	}

	p.backup()

	return false
}

// Next returns the next Token and advances the read position.
func (p *Parser) Next() Token {
	return p.get()
}

// Peek takes a look at the upcoming Token and returns it.
func (p *Parser) Peek() Token {
	tk := p.get()
	p.backup()

	return tk
}

// Get retrieves a slice of the Tokens that have been read so far.
func (p *Parser) Get() []Token {
	var toRet []Token

	if p.peekedToken {
		tk := p.tokens[len(p.tokens)-1]
		toRet = make([]Token, len(p.tokens)-1)
		copy(toRet, p.tokens)
		p.tokens[0] = tk
		p.tokens = p.tokens[:1]
	} else {
		toRet = make([]Token, len(p.tokens))
		copy(toRet, p.tokens)
		p.tokens = p.tokens[:0]
	}

	return toRet
}

// Len returns how many tokens have been read.
func (p *Parser) Len() int {
	l := len(p.tokens)

	if p.peekedToken {
		l--
	}

	return l
}

// AcceptRun will keep Accepting tokens as long as they match one of the
// given types.
//
// It will return the type of the token that made it stop.
func (p *Parser) AcceptRun(types ...TokenType) TokenType {
Loop:
	for {
		tk := p.get()

		for _, t := range types {
			if tk.Type == t {
				continue Loop
			}
		}

		p.backup()

		return tk.Type
	}
}

// AcceptToken will accept a token matching one of the ones provided exactly,
// returning true if one is read and false otherwise.
func (p *Parser) AcceptToken(tokens ...Token) bool {
	tk := p.get()

	for _, ttk := range tokens {
		if ttk == tk {
			return true
		}
	}

	p.backup()

	return false
}

// Except will Accept a token that is not one of the types given. Returns true
// if it Accepted a token.
func (p *Parser) Except(types ...TokenType) bool {
	tk := p.get()

	for _, t := range types {
		if tk.Type == t {
			p.backup()

			return false
		}
	}

	return true
}

// ExceptRun will keep Accepting tokens as long as they do not match one of the
// given types.
//
// It will return the type of the token that made it stop.
func (p *Parser) ExceptRun(types ...TokenType) TokenType {
	for {
		tk := p.get()

		for _, t := range types {
			if tk.Type == t || tk.Type < 0 {
				p.backup()

				return tk.Type
			}
		}
	}
}

// Return simplifies the returning from PhraseFns, taking a PhraseType and a
// next TokenFn, default to Done.
//
// The returned phrase is of the type specified with the data set to the output
// of p.Get().
func (p *Parser) Return(typ PhraseType, fn PhraseFunc) (Phrase, PhraseFunc) {
	if fn == nil {
		fn = (*Parser).Done
	}

	return Phrase{
		Type: typ,
		Data: p.Get(),
	}, fn
}

// Done is a PhraseFunc that is used to indicate that there are no more phrases
// to parse.
func (p *Parser) Done() (Phrase, PhraseFunc) {
	p.Err = io.EOF

	return Phrase{
		Type: PhraseDone,
		Data: make([]Token, 0),
	}, (*Parser).Done
}

// Error represents an error state for the phraser.
//
// The error value should be set in Parser.Err and then this func should be
// called.
func (p *Parser) Error() (Phrase, PhraseFunc) {
	if p.Err == nil {
		p.Err = ErrUnknownError
	}

	return Phrase{
		Type: PhraseError,
		Data: []Token{
			{Type: TokenError, Data: p.Err.Error()},
		},
	}, (*Parser).Error
}
