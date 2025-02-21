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
