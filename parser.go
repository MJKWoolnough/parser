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

// NewStringParser returns a Parser which parses a string.
func NewStringParser(str string) Parser {
	return Parser{
		phraser: Phraser{
			tokeniser: Tokeniser{
				tokeniser: &strParser{
					str: str,
				},
			},
		},
	}
}

// NewByteParser returns a Parser which parses a byte slice.
func NewByteParser(data []byte) Parser {
	return Parser{
		phraser: Phraser{
			tokeniser: Tokeniser{
				tokeniser: &byteParser{
					data: data,
				},
			},
		},
	}
}

// NewReaderParser returns a Parser which parses a Reader.
func NewReaderParser(reader io.Reader) Parser {
	return Parser{
		phraser: Phraser{
			tokeniser: Tokeniser{
				tokeniser: &readerParser{
					reader: bufio.NewReader(reader),
				},
			},
		},
	}
}

func (p *Parser) GetToken() (Token, error) {
	tk := p.phraser.tokeniser.get()
	if tk.Type == TokenError {
		return tk, p.phraser.tokeniser.err
	}
	return tk, nil
}

func (p *Parser) GetPhrase() (Phrase, error) {
	if p.phraser.tokeniser.err == io.EOF {
		return Phrase{
			Type: PhraseDone,
			Data: make([]Token, 0),
		}, io.EOF
	}
	if p.phraser.state == nil {
		p.phraser.tokeniser.err = ErrNoState
		p.phraser.state = (*Phraser).Error
	}
	var ph Phrase
	ph, p.phraser.state = p.phraser.state(&p.phraser)
	if ph.Type == PhraseError {
		if p.phraser.tokeniser.err == io.EOF {
			p.phraser.tokeniser.err = io.ErrUnexpectedEOF
		}
		return ph, p.phraser.tokeniser.err
	}
	return ph, nil
}

func (p *Parser) TokeniserState(t TokenFunc) {
	p.phraser.tokeniser.state = t
}

func (p *Parser) PhraserState(pf PhraseFunc) {
	p.phraser.state = pf
}
