// Package parser is a simple helper package for parsing strings, byte slices and io.Readers.
package parser // import "vimagination.zapto.org/parser"

import (
	"bufio"
	"io"
)

// New creates a new Parser from the given Tokeniser.
func New(t Tokeniser) Parser {
	return Parser{Tokeniser: t}
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

// NewReaderTokeniser returns a Tokeniser which uses an io.Reader.
func NewReaderTokeniser(reader io.Reader) Tokeniser {
	return Tokeniser{
		tokeniser: &readerParser{
			reader: bufio.NewReader(reader),
		},
	}
}

// NewRuneReaderTokeniser returns a Tokeniser which uses an io.RuneReader.
//
// Any rune errors will result in EOF.
func NewRuneReaderTokeniser(source io.RuneReader) Tokeniser {
	return Tokeniser{
		tokeniser: &runeSourceParser{
			source: source,
		},
	}
}
