package parser

import (
	"iter"
	"strings"
	"testing"
)

func tokenisers(str string) iter.Seq2[string, Tokeniser] {
	return func(yield func(string, Tokeniser) bool) {
		_ = yield("string", NewStringTokeniser(str)) && yield("bytes", NewByteTokeniser([]byte(str))) && yield("reader", NewReaderTokeniser(strings.NewReader(str)))
	}
}

func TestTokeniserAccept(t *testing.T) {
	for n, p := range tokenisers("ABC£") {
		p.Accept("ABCD")
		if s := p.Get(); s != "A" {
			t.Errorf("test 1 (%s): expecting \"A\", got %q", n, s)
		}

		p.Accept("ABCD")
		if s := p.Get(); s != "B" {
			t.Errorf("test 2 (%s): expecting \"B\", got %q", n, s)

			continue
		}

		p.Accept("ABCD")
		if s := p.Get(); s != "C" {
			t.Errorf("test 3 (%s): expecting \"C\", got %q", n, s)

			continue
		}

		p.Accept("ABCD")
		if s := p.Get(); s != "" {
			t.Errorf("test 4 (%s): expecting \"\", got %q", n, s)

			continue
		}

		p.Accept("£")
		if s := p.Get(); s != "£" {
			t.Errorf("test 5 (%s): expecting \"£\", got %q", n, s)

			continue
		}
	}
}

func TestTokeniserAcceptRun(t *testing.T) {
	for n, p := range tokenisers("123ABC££$$%%^^\n") {
		p.AcceptRun("0123456789")
		if s := p.Get(); s != "123" {
			t.Errorf("test 1 (%s): expecting \"123\", got %q", n, s)

			continue
		}

		p.AcceptRun("ABC")
		if s := p.Get(); s != "ABC" {
			t.Errorf("test 2 (%s): expecting \"ABC\", got %q", n, s)

			continue
		}

		p.AcceptRun("£$%^")
		if s := p.Get(); s != "££$$%%^^" {
			t.Errorf("test 3 (%s): expecting \"££$$%%^^\", got %q", n, s)

			continue
		}

		p.AcceptRun("\n")
		if s := p.Get(); s != "\n" {
			t.Errorf("test 4 (%s): expecting \"\\n\", got %q", n, s)

			continue
		}
	}
}

func TestTokeniserExcept(t *testing.T) {
	for n, p := range tokenisers("123") {
		p.Except("1")
		if s := p.Get(); s != "" {
			t.Errorf("test 1 (%s): expecting \"\", got %q", n, s)

			continue
		}

		p.Except("2")
		if s := p.Get(); s != "1" {
			t.Errorf("test 2 (%s): expecting \"1\", got %q", n, s)

			continue
		}

		p.Except("2")
		if s := p.Get(); s != "" {
			t.Errorf("test 3 (%s): expecting \"\", got %q", n, s)

			continue
		}

		p.Except("!")
		if s := p.Get(); s != "2" {
			t.Errorf("test 4 (%s): expecting \"2\", got %q", n, s)

			continue
		}

		p.Except("!")
		if s := p.Get(); s != "3" {
			t.Errorf("test 5 (%s): expecting \"3\", got %q", n, s)

			continue
		}

		p.Except("!")
		if s := p.Get(); s != "" {
			t.Errorf("test 6 (%s): expecting \"\", got %q", n, s)

			continue
		}
	}
}

func TestTokeniserExceptRun(t *testing.T) {
	for n, p := range tokenisers("12345ABC\n67890DEF\nOH MY!") {
		p.ExceptRun("\n")
		if s := p.Get(); s != "12345ABC" {
			t.Errorf("test 1 (%s): expecting \"12345ABC\", got %q", n, s)

			continue
		}

		p.Except("")
		p.Get()
		p.ExceptRun("\n")

		if s := p.Get(); s != "67890DEF" {
			t.Errorf("test 2 (%s): expecting \"67890DEF\", got %q", n, s)

			continue
		}

		p.Except("")
		p.Get()
		p.ExceptRun("")

		if s := p.Get(); s != "OH MY!" {
			t.Errorf("test 3 (%s): expecting \"OH MY!\", got %q", n, s)

			continue
		}
	}
}

func TestTokeniserState(t *testing.T) {
	for n, p := range tokenisers("12345ABC\n67890DEF\nOH MY!") {
		state := p.State()

		a := p.Next()
		b := p.Next()
		c := p.Next()
		d := p.Next()
		l := p.Len()

		state.Reset()

		if p.Next() != a || p.Next() != b || p.Next() != c || p.Next() != d || p.Len() != l {
			t.Errorf("test 1 (%s): failed to reset state correctly", n)

			continue
		}

		state = p.State()

		a = p.Next()
		b = p.Next()
		c = p.Next()
		d = p.Next()
		l = p.Len()

		state.Reset()

		if p.Next() != a || p.Next() != b || p.Next() != c || p.Next() != d || p.Len() != l {
			t.Errorf("test 2 (%s): failed to reset state correctly", n)

			continue
		}
	}
}

func TestTokeniserAcceptString(t *testing.T) {
	for m, p := range tokenisers("ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		for n, test := range [...]struct {
			Str             string
			Read            int
			CaseInsensitive bool
		}{
			{
				Str: "Z",
			},
			{
				Str:  "A",
				Read: 1,
			},
			{
				Str:  "BCD",
				Read: 3,
			},
			{
				Str:  "EFGZ",
				Read: 3,
			},
			{
				Str:  "hij",
				Read: 0,
			},
			{
				Str:             "hij",
				Read:            3,
				CaseInsensitive: true,
			},
		} {
			if read := p.AcceptString(test.Str, test.CaseInsensitive); read != test.Read {
				t.Errorf("test %d (%s): expecting to parse %d chars, parsed %d", n+1, m, test.Read, read)
			}
		}
	}
}

func TestTokeniserAcceptWord(t *testing.T) {
	for m, p := range tokenisers("ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		for n, test := range [...]struct {
			Words           []string
			Read            string
			CaseInsensitive bool
		}{
			{},
			{
				Words: []string{"Z"},
			},
			{
				Words: []string{"Z", "Y"},
			},
			{
				Words: []string{"A"},
				Read:  "A",
			},
			{
				Words: []string{"BD"},
			},
			{
				Words: []string{"BD", "BE"},
			},
			{
				Words: []string{"BCD", "BCE"},
				Read:  "BCD",
			},
			{
				Words: []string{"EFH", "EFG"},
				Read:  "EFG",
			},
			{
				Words: []string{"HIJ", "HIJK"},
				Read:  "HIJK",
			},
			{
				Words:           []string{"LMNOP", "LMOPQ", "LmNoPqR"},
				Read:            "LMNOPQR",
				CaseInsensitive: true,
			},
		} {
			if read := p.AcceptWord(test.Words, test.CaseInsensitive); read != test.Read {
				t.Errorf("test %d (%s): expecting to parse %q, parsed %q", n+1, m, test.Read, read)
			}
		}
	}
}
