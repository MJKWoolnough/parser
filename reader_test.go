package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"vimagination.zapto.org/parser"
)

func TestReaderAccept(t *testing.T) {
	testTokeniserAccept(t, parser.NewReaderTokeniser(strings.NewReader("ABC£")))
}

func TestReaderAcceptRun(t *testing.T) {
	testTokeniserAcceptRun(t, parser.NewReaderTokeniser(strings.NewReader("123ABC££$$%%^^\n")))
}

func TestReaderExcept(t *testing.T) {
	testTokeniserExcept(t, parser.NewReaderTokeniser(strings.NewReader("123")))
}

func TestReaderExceptRun(t *testing.T) {
	testTokeniserExceptRun(t, parser.NewReaderTokeniser(strings.NewReader("12345ABC\n67890DEF\nOH MY!")))
}

func TestReaderState(t *testing.T) {
	testTokeniserState(t, parser.NewReaderTokeniser(strings.NewReader("12345ABC\n67890DEF\nOH MY!")))
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
