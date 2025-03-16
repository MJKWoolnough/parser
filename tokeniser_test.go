package parser

import (
	"iter"
	"strings"
	"testing"
)

func tokenisers(str string) iter.Seq2[string, Tokeniser] {
	return func(yield func(string, Tokeniser) bool) {
		_ = yield("string", NewStringTokeniser(str)) &&
			yield("bytes", NewByteTokeniser([]byte(str))) &&
			yield("reader", NewReaderTokeniser(strings.NewReader(str))) &&
			yield("rune reader", NewRuneReaderTokeniser(strings.NewReader(str))) &&
			yield("sub (string)", Tokeniser{tokeniser: NewStringTokeniser(str).sub()}) &&
			yield("sub (bytes)", Tokeniser{tokeniser: NewByteTokeniser([]byte(str)).sub()}) &&
			yield("sub (reader)", Tokeniser{tokeniser: NewReaderTokeniser(strings.NewReader(str)).sub()}) &&
			yield("sub (rune reader)", Tokeniser{tokeniser: NewRuneReaderTokeniser(strings.NewReader(str)).sub()})
	}
}

func TestTokeniserNext(t *testing.T) {
	for n, p := range tokenisers("ABCDEFGH") {
		if c := p.Peek(); c != 'A' {
			t.Errorf("test 1 (%s): expecting %q, got %q", n, 'A', c)
		} else if c = p.Peek(); c != 'A' {
			t.Errorf("test 2 (%s): expecting %q, got %q", n, 'A', c)
		} else if c = p.Next(); c != 'A' {
			t.Errorf("test 3 (%s): expecting %q, got %q", n, 'A', c)
		} else if c = p.Next(); c != 'B' {
			t.Errorf("test 4 (%s): expecting %q, got %q", n, 'B', c)
		} else if c = p.Peek(); c != 'C' {
			t.Errorf("test 5 (%s): expecting %q, got %q", n, 'C', c)
		}
	}
}

func TestTokeniserLen(t *testing.T) {
	for n, p := range tokenisers("A…") {
		p.Peek()

		if l := p.Len(); l != 0 {
			t.Errorf("test 1 (%s): expecting to have read 0 bytes, read %d", n, l)
		}

		p.Next()

		if l := p.Len(); l != 1 {
			t.Errorf("test 2 (%s): expecting to have read 1 byte, read %d", n, l)
		}

		p.Next()

		if l := p.Len(); l != 4 {
			t.Errorf("test 3 (%s): expecting to have read 4 bytes, read %d", n, l)
		}

		p.Next()

		if l := p.Len(); l != 4 {
			t.Errorf("test 4 (%s): expecting to have read 4 bytes, read %d", n, l)
		}

		p.Next()

		if l := p.Len(); l != 4 {
			t.Errorf("test 5 (%s): expecting to have read 4 bytes, read %d", n, l)
		}
	}
}

func TestTokeniserAccept(t *testing.T) {
	for n, p := range tokenisers("ABC£") {
		if _, s := p.Accept("ABCD"), p.Get(); s != "A" {
			t.Errorf("test 1 (%s): expecting \"A\", got %q", n, s)
		} else if _, s = p.Accept("ABCD"), p.Get(); s != "B" {
			t.Errorf("test 2 (%s): expecting \"B\", got %q", n, s)
		} else if _, s = p.Accept("ABCD"), p.Get(); s != "C" {
			t.Errorf("test 3 (%s): expecting \"C\", got %q", n, s)
		} else if _, s = p.Accept("ABCD"), p.Get(); s != "" {
			t.Errorf("test 4 (%s): expecting \"\", got %q", n, s)
		} else if _, s = p.Accept("£"), p.Get(); s != "£" {
			t.Errorf("test 5 (%s): expecting \"£\", got %q", n, s)
		}
	}
}

func TestTokeniserAcceptRun(t *testing.T) {
	for n, p := range tokenisers("123ABC££$$%%^^\n") {
		if _, s := p.AcceptRun("0123456789"), p.Get(); s != "123" {
			t.Errorf("test 1 (%s): expecting \"123\", got %q", n, s)
		} else if _, s = p.AcceptRun("ABC"), p.Get(); s != "ABC" {
			t.Errorf("test 2 (%s): expecting \"ABC\", got %q", n, s)
		} else if _, s = p.AcceptRun("£$%^"), p.Get(); s != "££$$%%^^" {
			t.Errorf("test 3 (%s): expecting \"££$$%%^^\", got %q", n, s)
		} else if _, s = p.AcceptRun("\n"), p.Get(); s != "\n" {
			t.Errorf("test 4 (%s): expecting \"\\n\", got %q", n, s)
		}
	}
}

func TestTokeniserExcept(t *testing.T) {
	for n, p := range tokenisers("123") {
		if _, s := p.Except("1"), p.Get(); s != "" {
			t.Errorf("test 1 (%s): expecting \"\", got %q", n, s)
		} else if _, s = p.Except("2"), p.Get(); s != "1" {
			t.Errorf("test 2 (%s): expecting \"1\", got %q", n, s)
		} else if _, s = p.Except("2"), p.Get(); s != "" {
			t.Errorf("test 3 (%s): expecting \"\", got %q", n, s)
		} else if _, s = p.Except("!"), p.Get(); s != "2" {
			t.Errorf("test 4 (%s): expecting \"2\", got %q", n, s)
		} else if _, s = p.Except("!"), p.Get(); s != "3" {
			t.Errorf("test 5 (%s): expecting \"3\", got %q", n, s)
		} else if _, s = p.Except("!"), p.Get(); s != "" {
			t.Errorf("test 6 (%s): expecting \"\", got %q", n, s)
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

func TestTokeniserReset(t *testing.T) {
	for n, p := range tokenisers("ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		p.ExceptRun("E")
		p.Reset()

		if got := p.Get(); got != "" {
			t.Errorf("test 1 (%s): expecting to get %q, got %q", n, "", got)
		} else if _, got = p.ExceptRun("E"), p.Get(); got != "ABCD" {
			t.Errorf("test 2 (%s): expecting to get %q, got %q", n, "ABCD", got)
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
			{
				Words: []string{"ZYX", "ST", "STZ"},
				Read:  "ST",
			},
		} {
			if read := p.AcceptWord(test.Words, test.CaseInsensitive); read != test.Read {
				t.Errorf("test %d (%s): expecting to parse %q, parsed %q", n+1, m, test.Read, read)
			}
		}
	}
}

func TestTokeniserSub(t *testing.T) {
	for n, p := range tokenisers("ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		if _, ok := p.tokeniser.(*sub); ok {
			break
		} else if c := p.Next(); c != 'A' {
			t.Errorf("test 1 (%s): expecting to read %q, got %q", n, 'A', c)

			continue
		}

		q := p.SubTokeniser()

		if c := q.Next(); c != 'B' {
			t.Errorf("test 2 (%s): expecting to read %q, got %q", n, 'B', c)

			continue
		}

		if c := q.ExceptRun("H"); c != 'H' {
			t.Errorf("test 3 (%s): expecting to read %q, got %q", n, 'H', c)

			continue
		} else if got := q.Get(); got != "BCDEFG" {
			t.Errorf("test 4 (%s): expecting to read %q, got %q", n, "BCDEFG", got)

			continue
		}

		q.Next()

		if got := q.Get(); got != "H" {
			t.Errorf("test 5 (%s): expecting to read %q, got %q", n, "H", got)

			continue
		} else if got := p.Get(); got != "ABCDEFGH" {
			t.Errorf("test 6 (%s): expecting to read %q, got %q", n, "ABCDEFGH", got)

			continue
		}

		q.Next()

		if got := q.Get(); got != "" {
			t.Errorf("test 7 (%s): expecting to read %q, got %q", n, "", got)

			continue
		}

		p.Next()

		q = p.SubTokeniser()

		q.Next()

		r := q.SubTokeniser()

		r.Next()

		if got := r.Get(); got != "L" {
			t.Errorf("test 8 (%s): expecting to read %q, got %q", n, "L", got)
		} else if got := q.Get(); got != "KL" {
			t.Errorf("test 9 (%s): expecting to read %q, got %q", n, "KL", got)
		} else if got := p.Get(); got != "IJKL" {
			t.Errorf("test 10 (%s): expecting to read %q, got %q", n, "HIJKL", got)
		}
	}
}
