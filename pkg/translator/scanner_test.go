package translator

import (
	"testing"
)

type tokWant struct {
	pos Position
	tok int
	lit string
}

type test struct {
	tokWant []tokWant
	in      string
}

func checkScanResult(t *testing.T, pos Position, tok int, lit string, c test, i int) {
	if pos.Offset != c.tokWant[i].pos.Offset {
		t.Errorf("s(%v).Scan() offset error: expected %v, got %v", c.in, c.tokWant[i].pos.Offset, pos.Offset)
	}

	if pos.Column != c.tokWant[i].pos.Column {
		t.Errorf("s(%v).Scan() column error: expected %v, got %v", c.in, c.tokWant[i].pos.Column, pos.Column)
	}

	if pos.Line != c.tokWant[i].pos.Line {
		t.Errorf("s(%v).Scan() line error: expected %v, got %v", c.in, c.tokWant[i].pos.Line, pos.Line)
	}

	if tok != c.tokWant[i].tok {
		t.Errorf("s(%v).Scan() token error: expected %v (%v), got %v (%v)", c.in, c.tokWant[i].tok,
			Tokens[c.tokWant[i].tok], tok, Tokens[tok])
	}

	if lit != c.tokWant[i].lit {
		t.Errorf("s(%v).Scan() token error: expected %v, got %v", c.in, c.tokWant[i].lit, lit)
	}
}

var testHeredoc = [...]test{
	{
		[]tokWant{
			{
				Position{Offset: 0, Line: 1, Column: 1},
				LHEREDOC,
				"<<test.yaml",
			},
			{
				Position{Offset: 20, Line: 1, Column: 21},
				RHEREDOC,
				"test.yaml>>",
			},
		},
		"<<test.yaml foo:bar test.yaml>>",
	},
	{
		[]tokWant{
			{
				Position{Offset: 0, Line: 1, Column: 1},
				LHEREDOC,
				"<<test.yaml",
			},
			{
				Position{Offset: 24, Line: 3, Column: 2},
				RHEREDOC,
				"test.yaml>>",
			},
		},
		"<<test.yaml \n foo:bar \n test.yaml>>",
	},
	{
		[]tokWant{
			{
				Position{Offset: 0, Line: 1, Column: 1},
				LHEREDOC,
				"<<test.yaml",
			},
			{
				Position{Offset: 18, Line: 2, Column: 6},
				LINJECT,
				"{{",
			},
			{
				Position{Offset: 20, Line: 2, Column: 8},
				CODE,
				"var bar := 2",
			},
			{
				Position{Offset: 32, Line: 2, Column: 20},
				RINJECT,
				"}}",
			},
			{
				Position{Offset: 37, Line: 3, Column: 2},
				RHEREDOC,
				"test.yaml>>",
			},
		},
		"<<test.yaml \n foo:{{var bar := 2}} \n test.yaml>>",
	},
}

func TestScanHeredoc(t *testing.T) {
	for _, c := range testHeredoc {
		bytes := []byte(c.in)
		s := NewScanner(bytes)
		tokens := s.Scan()

		for i := range c.tokWant {
			checkScanResult(t, tokens[i].pos, tokens[i].tok, tokens[i].lit, c, i)
		}
	}
}
