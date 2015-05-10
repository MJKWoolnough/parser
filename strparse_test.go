package strparse_test

import (
	"fmt"
	"testing"

	"github.com/MJKWoolnough/strparse"
)

func TestAccept(t *testing.T) {
	p := strparse.Parser{Str: "ABC£"}
	p.Accept("ABCD")
	if s := p.Get(); s != "A" {
		t.Errorf("expecting \"A\", got %q", s)
		return
	}
	p.Accept("ABCD")
	if s := p.Get(); s != "B" {
		t.Errorf("expecting \"B\", got %q", s)
		return
	}
	p.Accept("ABCD")
	if s := p.Get(); s != "C" {
		t.Errorf("expecting \"C\", got %q", s)
		return
	}
	p.Accept("ABCD")
	if s := p.Get(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
		return
	}
	p.Accept("£")
	if s := p.Get(); s != "£" {
		t.Errorf("expecting \"£\", got %q", s)
		return
	}
}

func TestAcceptRun(t *testing.T) {
	p := strparse.Parser{Str: "123ABC££$$%%^^\n"}
	p.AcceptRun("0123456789")
	if s := p.Get(); s != "123" {
		t.Errorf("expecting \"123\", got %q", s)
		return
	}
	p.AcceptRun("ABC")
	if s := p.Get(); s != "ABC" {
		t.Errorf("expecting \"ABC\", got %q", s)
		return
	}
	p.AcceptRun("£$%^")
	if s := p.Get(); s != "££$$%%^^" {
		t.Errorf("expecting \"££$$%%^^\", got %q", s)
		return
	}
	p.AcceptRun("\n")
	if s := p.Get(); s != "\n" {
		t.Errorf("expecting \"\\n\", got %q", s)
		return
	}
}

func TestExcept(t *testing.T) {
	p := strparse.Parser{Str: "123"}
	p.Except("1")
	if s := p.Get(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
	}
	p.Except("2")
	if s := p.Get(); s != "1" {
		t.Errorf("expecting \"1\", got %q", s)
	}
	p.Except("2")
	if s := p.Get(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
	}
	p.Except("!")
	if s := p.Get(); s != "2" {
		t.Errorf("expecting \"2\", got %q", s)
	}
	p.Except("!")
	if s := p.Get(); s != "3" {
		t.Errorf("expecting \"3\", got %q", s)
	}
	p.Except("!")
	if s := p.Get(); s != "" {
		t.Errorf("expecting \"\", got %q", s)
	}
}

func TestExceptRun(t *testing.T) {
	p := strparse.Parser{Str: "12345ABC\n67890DEF\nOH MY!"}
	p.ExceptRun("\n")
	if s := p.Get(); s != "12345ABC" {
		t.Errorf("expecting \"12345ABC\", got %q", s)
		return
	}
	p.Except("")
	p.Get()
	p.ExceptRun("\n")
	if s := p.Get(); s != "67890DEF" {
		t.Errorf("expecting \"67890DEF\", got %q", s)
		return
	}
	p.Except("")
	p.Get()
	p.ExceptRun("")
	if s := p.Get(); s != "OH MY!" {
		t.Errorf("expecting \"OH MY!\", got %q", s)
		return
	}
}

func Example() {
	p := strparse.New("Hello, World!")
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
