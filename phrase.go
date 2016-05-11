package parser

import "io"

type PhraseType int

const (
	PhraseDone PhraseType = -1 - iota
	PhraseError
)

type PhraseFunc func(*Phraser) (Phrase, PhraseFunc)

type Phrase struct {
	Type PhraseType
	Data []Token
}

type Phraser struct {
	tokeniser   Tokeniser
	state       PhraseFunc
	tokens      []Token
	peekedToken bool
}

func (p *Phraser) get() Token {
	if p.peekedToken {
		p.peekedToken = false
		return p.tokens[len(p.tokens)-1]
	}
	tk := p.tokeniser.get()
	if tk.Type >= 0 {
		p.tokens = append(p.tokens, tk)
	}
	return tk
}

func (p *Phraser) backup() {
	p.peekedToken = true
}

func (p *Phraser) Accept(types ...TokenType) bool {
	tk := p.get()
	for _, t := range types {
		if tk.Type == t {
			return true
		}
	}
	p.backup()
	return false
}

func (p *Phraser) Peek() Token {
	tk := p.get()
	p.backup()
	return tk
}

func (p *Phraser) Phrase() []Token {
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

func (p *Phraser) Len() int {
	l := len(p.tokens)
	if p.peekedToken {
		l--
	}
	return l
}

func (p *Phraser) AcceptRun(types ...TokenType) TokenType {
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

func (p *Phraser) Except(types ...TokenType) bool {
	tk := p.get()
	for _, t := range types {
		if tk.Type == t {
			p.backup()
			return false
		}
	}
	return true
}

func (p *Phraser) ExceptRun(types ...TokenType) TokenType {
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

func (p *Phraser) Done() (Phrase, PhraseFunc) {
	p.tokeniser.err = io.EOF
	return Phrase{
		Type: PhraseDone,
		Data: make([]Token, 0),
	}, (*Phraser).Done
}

func (p *Phraser) Error() (Phrase, PhraseFunc) {
	return Phrase{
		Type: PhraseError,
		Data: []Token{
			{Type: TokenError, Data: p.tokeniser.err.Error()},
		},
	}, (*Phraser).Error
}
