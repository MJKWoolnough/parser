package parser_test

import (
	"fmt"
	"testing"

	"github.com/MJKWoolnough/parser"
)

func getTokeniser(p parser.Parser) *parser.Tokeniser {
	var tk *parser.Tokeniser
	p.TokeniserState(func(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
		tk = t
		return parser.Token{}, nil
	})
	p.GetToken()
	return tk
}

func TestByteAccept(t *testing.T) {
	p := getTokeniser(parser.NewByteParser([]byte("ABC£")))
	p.Accept("ABCD")
	if s := p.Lexeme(); s != "A" {
		t.Errorf("expecting \"A\", got %q", s)
		return
	}
	p.Accept("ABCD")
	if s := p.Lexeme(); s != "B" {
		t.Errorf("expecting \"B\", got %q", s)
		return
	}
	p.Accept("ABCD")
	if s := p.Lexeme(); s != "C" {
		t.Errorf("expecting \"C\", got %q", s)
		return
	}
	p.Accept("ABCD")
	if s := p.Lexeme(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
		return
	}
	p.Accept("£")
	if s := p.Lexeme(); s != "£" {
		t.Errorf("expecting \"£\", got %q", s)
		return
	}
}

func TestByteAcceptRun(t *testing.T) {
	p := getTokeniser(parser.NewByteParser([]byte("123ABC££$$%%^^\n")))
	p.AcceptRun("0123456789")
	if s := p.Lexeme(); s != "123" {
		t.Errorf("expecting \"123\", got %q", s)
		return
	}
	p.AcceptRun("ABC")
	if s := p.Lexeme(); s != "ABC" {
		t.Errorf("expecting \"ABC\", got %q", s)
		return
	}
	p.AcceptRun("£$%^")
	if s := p.Lexeme(); s != "££$$%%^^" {
		t.Errorf("expecting \"££$$%%^^\", got %q", s)
		return
	}
	p.AcceptRun("\n")
	if s := p.Lexeme(); s != "\n" {
		t.Errorf("expecting \"\\n\", got %q", s)
		return
	}
}

func TestByteExcept(t *testing.T) {
	p := getTokeniser(parser.NewByteParser([]byte("123")))
	p.Except("1")
	if s := p.Lexeme(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
	}
	p.Except("2")
	if s := p.Lexeme(); s != "1" {
		t.Errorf("expecting \"1\", got %q", s)
	}
	p.Except("2")
	if s := p.Lexeme(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
	}
	p.Except("!")
	if s := p.Lexeme(); s != "2" {
		t.Errorf("expecting \"2\", got %q", s)
	}
	p.Except("!")
	if s := p.Lexeme(); s != "3" {
		t.Errorf("expecting \"3\", got %q", s)
	}
	p.Except("!")
	if s := p.Lexeme(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
	}
}

func TestByteExceptRun(t *testing.T) {
	p := getTokeniser(parser.NewByteParser([]byte("12345ABC\n67890DEF\nOH MY!")))
	p.ExceptRun("\n")
	if s := p.Lexeme(); s != "12345ABC" {
		t.Errorf("expecting \"12345ABC\", got %q", s)
		return
	}
	p.Except("")
	p.Lexeme()
	p.ExceptRun("\n")
	if s := p.Lexeme(); s != "67890DEF" {
		t.Errorf("expecting \"67890DEF\", got %q", s)
		return
	}
	p.Except("")
	p.Lexeme()
	p.ExceptRun("")
	if s := p.Lexeme(); s != "OH MY!" {
		t.Errorf("expecting \"OH MY!\", got %q", s)
		return
	}
}

func ExampleNewByteParser() {
	p := getTokeniser(parser.NewByteParser([]byte("Hello, World!"))) // needs fixing
	alphaNum := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	p.AcceptRun(alphaNum)
	word := p.Lexeme()
	fmt.Println("got word:", word)

	p.ExceptRun(alphaNum)
	p.Lexeme()

	p.AcceptRun(alphaNum)
	word = p.Lexeme()
	fmt.Println("got word:", word)
	// Output:
	// got word: Hello
	// got word: World
}
