package parser

import (
	"testing"
)

func TestTokeniserAcceptString(t *testing.T) {
	p := NewStringTokeniser("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
			t.Errorf("test %d: expecting to parse %d chars, parsed %d", n+1, test.Read, read)
		}
	}
}

func TestTokeniserAcceptWord(t *testing.T) {
	p := NewStringTokeniser("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
			t.Errorf("test %d: expecting to parse %q, parsed %q", n+1, test.Read, read)
		}
	}
}
