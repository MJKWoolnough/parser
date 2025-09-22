# parser

[![CI](https://github.com/MJKWoolnough/parser/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/parser/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/parser.svg)](https://pkg.go.dev/vimagination.zapto.org/parser)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/parser)](https://goreportcard.com/report/vimagination.zapto.org/parser)

--
    import "vimagination.zapto.org/parser"

Package parser is a simple helper package for parsing strings, byte slices and io.Readers.

## Highlights

 - Methods to accept a character from a string (`Accept`), or a run of such characters (`AcceptRun`).
 - Methods to accept a character not in a string (`Except`), or a run of such characters (`ExceptRun`).
 - Methods to accept whole strings, or one of many strings.
 - SubTokenisers and state storing to allow forward checking before concluding token type.

## Usage

```go
package main

import (
	"fmt"

	"vimagination.zapto.org/parser"
)

const (
	TokenWhitespace parser.TokenType = iota
	TokenWord
)

func whitespace (t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.AcceptRun(" ")

	if t.Len() == 0 {
		return t.Done()
	}

	return t.Return(TokenWhitespace, word)
}

func word(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.ExceptRun(" ")

	if t.Len() == 0 {
		return t.Done()
	}

	return t.Return(TokenWord, whitespace)
}

func start(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if t.Peek() == ' ' {
		return whitespace(t)
	}

	return word(t)
}

func main() {
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
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/parser
