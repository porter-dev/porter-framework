package json

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

type encoderTest struct {
	name string
	val  v.Value
	want string
}

var encoderTestsLiteralsPass = []encoderTest{
	encoderTest{
		name: "Literal: integer",
		val:  v.Integer(1),
		want: "1",
	},
	encoderTest{
		name: "Literal: float",
		val:  v.Float(0.1),
		want: "1e-01",
	},
	encoderTest{
		name: "Literal: true",
		val:  v.Boolean(true),
		want: "true",
	},
	encoderTest{
		name: "Literal: false",
		val:  v.Boolean(false),
		want: "false",
	},
	encoderTest{
		name: "Literal: null",
		val:  nil,
		want: "null",
	},
	encoderTest{
		name: "Literal: string",
		val:  v.String("hello"),
		want: "\"hello\"",
	},
}

func TestEncoderLiteralPass(t *testing.T) {
	for _, c := range encoderTestsLiteralsPass {
		got, _ := ToJSON(c.val)

		if got != c.want {
			t.Errorf("Failed on: %v, %s, %s", c.val, got, c.want)
		}
	}
}

var encoderTestsArrayPass = []encoderTest{
	encoderTest{
		name: "Array: empty array",
		val:  v.Array{},
		want: "[]",
	},
	encoderTest{
		name: "Array: with integers",
		val: v.Array{
			v.Integer(0),
			v.Integer(1),
			v.Integer(2),
		},
		want: "[0,1,2]",
	},
	encoderTest{
		name: "Array: with floats",
		val: v.Array{
			v.Float(0.1),
			v.Float(0.2),
			v.Float(0.3),
		},
		want: "[1e-01,2e-01,3e-01]",
	},
	encoderTest{
		name: "Array: with named literals",
		val: v.Array{
			v.Boolean(true),
			v.Boolean(false),
			nil,
		},
		want: "[true,false,null]",
	},
	encoderTest{
		name: "Array: nested arrays",
		val: v.Array{
			v.Array{
				v.Array{},
			},
			v.Array{},
		},
		want: "[[[]],[]]",
	},
}

func TestEncoderArrayPass(t *testing.T) {
	for _, c := range encoderTestsArrayPass {
		got, _ := ToJSON(c.val)

		if got != c.want {
			t.Errorf("Failed on: %v, %s, %s", c.val, got, c.want)
		}
	}
}

var encoderTestsObjectPass = []encoderTest{
	encoderTest{
		name: "Object: empty object",
		val:  v.Object{},
		want: "{}",
	},
	encoderTest{
		name: "Object: basic object",
		val: v.Object{
			v.String("foo"): v.String("bar"),
		},
		want: "{\"foo\":\"bar\"}",
	},
	encoderTest{
		name: "Object: nested object",
		val: v.Object{
			v.String("hello"): v.Object{
				v.String("there"): v.Object{
					v.String("general"): v.String("kenobi"),
				},
			},
			v.String("!"): v.Object{},
		},
		want: "{\"hello\":{\"there\":{\"general\":\"kenobi\"}},\"!\":{}}",
	},
}

func TestEncoderObjectPass(t *testing.T) {
	for _, c := range encoderTestsObjectPass {
		got, _ := ToJSON(c.val)

		if got != c.want {
			t.Errorf("Failed on: %s, %v, %s, %s", c.name, c.val, got, c.want)
		}
	}
}

var encoderTestsStructurePass = []encoderTest{
	encoderTest{
		name: "Structure: mixed array with nested object",
		val: v.Array{
			v.String("hello"),
			v.String("there"),
			v.Object{
				v.String("general"): v.String("kenobi"),
			},
		},
		want: "[\"hello\",\"there\",{\"general\":\"kenobi\"}]",
	},
}

func TestEncoderStructurePass(t *testing.T) {
	for _, c := range encoderTestsStructurePass {
		got, _ := ToJSON(c.val)

		if got != c.want {
			t.Errorf("Failed on: %v, %s, %s", c.val, got, c.want)
		}
	}
}
