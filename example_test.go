package parser_test

import (
	"fmt"

	"vimagination.zapto.org/parser"
)

func Example() {
	const (
		TokenWhitespace parser.TokenType = iota
		TokenWord
	)

	var start, word, whitespace parser.TokenFunc

	whitespace = func(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
		t.AcceptRun(" ")

		if t.Len() == 0 {
			return t.Done()
		}

		return t.Return(TokenWhitespace, word)
	}

	word = func(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
		t.ExceptRun(" ")

		if t.Len() == 0 {
			return t.Done()
		}

		return t.Return(TokenWord, whitespace)
	}

	start = func(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
		if t.Peek() == ' ' {
			return whitespace(t)
		}

		return word(t)
	}

	p := parser.New(parser.NewStringTokeniser("Hello World    Foo Bar"))

	p.TokeniserState(start)

	for p.Peek().Type != parser.TokenDone {
		tk := p.Next()
		typ := "word"

		if tk.Type == TokenWhitespace {
			typ = "whitespace"
		}

		fmt.Printf("got token (%s): %q\n", typ, tk.Data)
	}

	// Output:
	// got token (word): "Hello"
	// got token (whitespace): " "
	// got token (word): "World"
	// got token (whitespace): "    "
	// got token (word): "Foo"
	// got token (whitespace): " "
	// got token (word): "Bar"
}
