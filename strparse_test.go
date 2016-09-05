package parser_test

import (
	"fmt"
	"testing"

	"github.com/MJKWoolnough/parser"
)

func TestStrAccept(t *testing.T) {
	testTokeniserAccept(t, parser.NewStringTokeniser("ABC£"))
}

func TestStrAcceptRun(t *testing.T) {
	testTokeniserAcceptRun(t, parser.NewStringTokeniser("123ABC££$$%%^^\n"))
}

func TestStrExcept(t *testing.T) {
	testTokeniserExcept(t, parser.NewStringTokeniser("123"))
}

func TestStrExceptRun(t *testing.T) {
	testTokeniserExceptRun(t, parser.NewStringTokeniser("12345ABC\n67890DEF\nOH MY!"))
}

func ExampleNewStringTokeniserParser() {
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
