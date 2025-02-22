package parser_test

import (
	"fmt"
	"strings"

	"vimagination.zapto.org/parser"
)

func ExampleNewByteTokeniser() {
	p := parser.NewByteTokeniser([]byte("Hello, World!"))
	alphaNum := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	p.AcceptRun(alphaNum)
	word := p.Get()
	fmt.Println("got word:", word)

	p.ExceptRun(alphaNum)
	p.Get()

	p.AcceptRun(alphaNum)
	word = p.Get()
	fmt.Println("got word:", word)
	// Output:
	// got word: Hello
	// got word: World
}

func ExampleNewStringTokeniser() {
	p := parser.NewStringTokeniser("Hello, World!")
	alphaNum := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	p.AcceptRun(alphaNum)
	word := p.Get()
	fmt.Println("got word:", word)

	p.ExceptRun(alphaNum)
	p.Get()

	p.AcceptRun(alphaNum)
	word = p.Get()
	fmt.Println("got word:", word)
	// Output:
	// got word: Hello
	// got word: World
}

func ExampleNewReaderTokeniser() {
	p := parser.NewReaderTokeniser(strings.NewReader("Hello, World!"))
	alphaNum := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	p.AcceptRun(alphaNum)
	word := p.Get()
	fmt.Println("got word:", word)

	p.ExceptRun(alphaNum)
	p.Get()

	p.AcceptRun(alphaNum)
	word = p.Get()
	fmt.Println("got word:", word)
	// Output:
	// got word: Hello
	// got word: World
}
