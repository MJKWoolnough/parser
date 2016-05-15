package parser

import "io"

type PhraseType int

const (
	PhraseDone PhraseType = -1 - iota
	PhraseError
)

type PhraseFunc func(*Parser) (Phrase, PhraseFunc)

type Phrase struct {
	Type PhraseType
	Data []Token
}

type Parser struct {
	Tokeniser
	state       PhraseFunc
	tokens      []Token
	peekedToken bool
}

func (p *Parser) GetPhrase() (Phrase, error) {
	if p.Err == io.EOF {
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
		if p.Err == io.EOF {
			p.Err = io.ErrUnexpectedEOF
		}
		return ph, p.Err
	}
	return ph, nil
}

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

func (p *Parser) Peek() Token {
	tk := p.get()
	p.backup()
	return tk
}

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

func (p *Parser) Len() int {
	l := len(p.tokens)
	if p.peekedToken {
		l--
	}
	return l
}

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

func (p *Parser) Done() (Phrase, PhraseFunc) {
	p.Err = io.EOF
	return Phrase{
		Type: PhraseDone,
		Data: make([]Token, 0),
	}, (*Parser).Done
}

func (p *Parser) Error() (Phrase, PhraseFunc) {
	return Phrase{
		Type: PhraseError,
		Data: []Token{
			{Type: TokenError, Data: p.Err.Error()},
		},
	}, (*Parser).Error
}
