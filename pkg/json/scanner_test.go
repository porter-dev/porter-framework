package json

import (
	"testing"
)

// define token classes as enum
const (
	literal int = iota
	operator
)

func TestNew(t *testing.T) {
	str := "testing"
	bytes := []byte(str)

	s := NewScanner(bytes)

	if s.srcPos.Line != 1 {
		t.Errorf("Scan source position line number must be 1.")
	}
}

func TestNextWithValidBytes(t *testing.T) {
	str := "abcd"
	bytes := []byte(str)

	s := NewScanner(bytes)

	runes := []rune(str)

	for _, r := range runes {
		got, _ := s.next()

		if got != r {
			t.Errorf("s.next(): expected %v, got %v", r, got)
		}
	}
}

func TestNextWithRuneError(t *testing.T) {
	bytes := []byte("\"\\FFFD")
	s := NewScanner(bytes)

	exp := rune(bytes[0])
	got, _ := s.next()

	if got != exp {
		t.Errorf("s.next() with UTF-8 replacement char: expected %v, got %v", exp, got)
	}
}

type want struct {
	richTok RichToken
	class   int
}

type singleTest struct {
	want want
	in   string
}

func checkScanResult(t *testing.T, richTok RichToken, c singleTest) {
	if richTok.pos.Offset != c.want.richTok.pos.Offset {
		t.Errorf("s(%v).Scan() offset error: expected %v, got %v", c.in, c.want.richTok.pos.Offset, richTok.pos.Offset)
	}

	if richTok.pos.Column != c.want.richTok.pos.Column {
		t.Errorf("s(%v).Scan() column error: expected %v, got %v", c.in, c.want.richTok.pos.Column, richTok.pos.Column)
	}

	if richTok.pos.Line != c.want.richTok.pos.Line {
		t.Errorf("s(%v).Scan() line error: expected %v, got %v", c.in, c.want.richTok.pos.Line, richTok.pos.Line)
	}

	if richTok.tok != c.want.richTok.tok {
		t.Errorf("s(%v).Scan() token error: expected %v (%v), got %v (%v)", c.in, c.want.richTok.tok,
			Tokens[c.want.richTok.tok], richTok.tok, Tokens[richTok.tok])
	}

	if richTok.lit != c.want.richTok.lit {
		t.Errorf("s(%v).Scan() token error: expected %v, got %v", c.in, c.want.richTok.lit, richTok.lit)
	}

	switch c.want.class {
	case literal:
		if !richTok.tok.IsLiteral() {
			t.Errorf("s(%v).Scan() class error: token %v (%v) should be literal", c.in, richTok.tok,
				Tokens[richTok.tok])
		}
	case operator:
		if !richTok.tok.IsOperator() {
			t.Errorf("s(%v).Scan() class error: token %v (%v) should be operator", c.in, richTok.tok,
				Tokens[richTok.tok])
		}
	}
}

var singleTestLiterals = [...]singleTest{
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 1, Line: 1, Column: 2},
				tok: STRING,
				lit: "Arbitrary characters !@#$%^&*()",
			},
			class: literal,
		},
		"\"Arbitrary characters !@#$%^&*()\"",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: NUMBER,
				lit: "0",
			},
			class: literal,
		},
		"0",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: NUMBER,
				lit: "0.1",
			},
			class: literal,
		},
		"0.1",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: NUMBER,
				lit: "-0.1",
			},
			class: literal,
		},
		"-0.1",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: NUMBER,
				lit: "-0.1e123",
			},
			class: literal,
		},
		"-0.1e123",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: NUMBER,
				lit: "-0.1e+12",
			},
			class: literal,
		},
		"-0.1e+12",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: NUMBER,
				lit: "0.1e-12",
			},
			class: literal,
		},
		"0.1e-12",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: TRUE,
				lit: "true",
			},
			class: literal,
		},
		"true",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: FALSE,
				lit: "false",
			},
			class: literal,
		},
		"false",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: NULL,
				lit: "null",
			},
			class: literal,
		},
		"null",
	},
}

func TestScanLiterals(t *testing.T) {
	for _, c := range singleTestLiterals {
		bytes := []byte(c.in)
		s := NewScanner(bytes)
		richTok, _ := s.Scan()

		checkScanResult(t, richTok, c)
	}
}

var singleTestOperators = [...]singleTest{
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: LBRACE,
				lit: "{",
			},
			class: operator,
		},
		"{",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: LBRACK,
				lit: "[",
			},
			class: operator,
		},
		"[",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: RBRACE,
				lit: "}",
			},
			class: operator,
		},
		"}",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: RBRACK,
				lit: "]",
			},
			class: operator,
		},
		"]",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: COMMA,
				lit: ",",
			},
			class: operator,
		},
		",",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 0, Line: 1, Column: 1},
				tok: COLON,
				lit: ":",
			},
			class: operator,
		},
		":",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 1, Line: 1, Column: 2},
				tok: LINJECT,
				lit: "{{",
			},
			class: operator,
		},
		"{{",
	},
	{
		want{
			richTok: RichToken{
				pos: Position{Offset: 1, Line: 1, Column: 2},
				tok: RINJECT,
				lit: "}}",
			},
			class: operator,
		},
		"}}",
	},
}

func TestScanOperators(t *testing.T) {
	for _, c := range singleTestOperators {
		bytes := []byte(c.in)
		s := NewScanner(bytes)
		richTok, _ := s.Scan()

		checkScanResult(t, richTok, c)
	}
}

type multiTest struct {
	want []want
	in   string
}

var multiTests = [...]multiTest{
	{
		[]want{
			want{
				richTok: RichToken{
					pos: Position{Offset: 0, Line: 1, Column: 1},
					tok: LBRACE,
					lit: "{",
				},
				class: operator,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 4, Line: 2, Column: 2},
					tok: STRING,
					lit: "foo",
				},
				class: literal,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 8, Line: 2, Column: 6},
					tok: COLON,
					lit: ":",
				},
				class: operator,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 11, Line: 2, Column: 9},
					tok: STRING,
					lit: "bar",
				},
				class: literal,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 17, Line: 3, Column: 2},
					tok: RBRACE,
					lit: "}",
				},
				class: operator,
			},
		},
		"{ \n\"foo\": \"bar\"\n }",
	},
	{
		[]want{
			want{
				richTok: RichToken{
					pos: Position{Offset: 0, Line: 1, Column: 1},
					tok: LBRACE,
					lit: "{",
				},
				class: operator,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 4, Line: 2, Column: 2},
					tok: STRING,
					lit: "foo",
				},
				class: literal,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 8, Line: 2, Column: 6},
					tok: COLON,
					lit: ":",
				},
				class: operator,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 10, Line: 2, Column: 8},
					tok: LBRACK,
					lit: "[",
				},
				class: operator,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 11, Line: 2, Column: 9},
					tok: NUMBER,
					lit: "0",
				},
				class: literal,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 12, Line: 2, Column: 10},
					tok: COMMA,
					lit: ",",
				},
				class: operator,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 14, Line: 2, Column: 12},
					tok: NUMBER,
					lit: "1",
				},
				class: literal,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 15, Line: 2, Column: 13},
					tok: RBRACK,
					lit: "]",
				},
				class: operator,
			},
			want{
				richTok: RichToken{
					pos: Position{Offset: 18, Line: 3, Column: 2},
					tok: RBRACE,
					lit: "}",
				},
				class: operator,
			},
		},
		"{ \n\"foo\": [0, 1]\n }",
	},
}

func TestScanMulti(t *testing.T) {
	for _, c := range multiTests {
		bytes := []byte(c.in)
		s := NewScanner(bytes)

		for _, w := range c.want {
			richTok, _ := s.Scan()

			checkScanResult(t, richTok, singleTest{
				want: w,
				in:   c.in,
			})
		}
	}
}
