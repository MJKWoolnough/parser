// Package parser is a simple helper package for parsing strings, byte slices and io.Readers
package parser

import (
	"bufio"
	"io"
)

// Parser is the wrapper type for the various different parsers.
type Parser struct {
	phraser Phraser
}

// New wraps a Tokeniser to provide additional functions useful to parsing
func New(t Tokeniser) Parser {
	return Parser{
		phraser: Phraser{
			tokeniser: t,
		},
	}
}

// NewStringTokeniser returns a Tokeniser which uses a string.
func NewStringTokeniser(str string) Tokeniser {
	return Tokeniser{
		tokeniser: &strParser{
			str: str,
		},
	}
}

// NewByteTokeniser returns a Tokeniser which uses a byte slice.
func NewByteTokeniser(data []byte) Tokeniser {
	return Tokeniser{
		tokeniser: &byteParser{
			data: data,
		},
	}
}

// NewReaderTokeniser returns a Tokeniser which uses an io.Reader
func NewReaderTokeniser(reader io.Reader) Tokeniser {
	return Tokeniser{
		tokeniser: &readerParser{
			reader: bufio.NewReader(reader),
		},
	}
}

func (p *Parser) GetToken() (Token, error) {
	tk := p.phraser.tokeniser.get()
	if tk.Type == TokenError {
		return tk, p.phraser.tokeniser.Err
	}
	return tk, nil
}

func (p *Parser) GetPhrase() (Phrase, error) {
	if p.phraser.tokeniser.Err == io.EOF {
		return Phrase{
			Type: PhraseDone,
			Data: make([]Token, 0),
		}, io.EOF
	}
	if p.phraser.state == nil {
		p.phraser.tokeniser.Err = ErrNoState
		p.phraser.state = (*Phraser).Error
	}
	var ph Phrase
	ph, p.phraser.state = p.phraser.state(&p.phraser)
	if ph.Type == PhraseError {
		if p.phraser.tokeniser.Err == io.EOF {
			p.phraser.tokeniser.Err = io.ErrUnexpectedEOF
		}
		return ph, p.phraser.tokeniser.Err
	}
	return ph, nil
}

func (p *Parser) TokeniserState(t TokenFunc) {
	p.phraser.tokeniser.state = t
}

func (p *Parser) PhraserState(pf PhraseFunc) {
	p.phraser.state = pf
}
