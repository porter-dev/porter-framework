package json

import "testing"

func TestIsLiteral(t *testing.T) {
	cases := []struct {
		in   Token
		want bool
	}{
		{in: NUMBER, want: true},
		{in: LBRACE, want: false},
	}

	for _, c := range cases {
		got := c.in.IsLiteral()
		if got != c.want {
			t.Errorf("%q.IsLiteral() == %t, want %t", c.in, got, c.want)
		}
	}
}

func TestIsOperator(t *testing.T) {
	cases := []struct {
		in   Token
		want bool
	}{
		{in: RBRACK, want: true},
		{in: NUMBER, want: false},
	}

	for _, c := range cases {
		got := c.in.IsOperator()
		if got != c.want {
			t.Errorf("%q.IsOperator() == %t, want %t", c.in, got, c.want)
		}
	}
}
