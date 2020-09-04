package value

import (
	"testing"
)

type isEqualTest struct {
	name string
	v1   Value
	v2   Value
	want bool
}

var isEqualTests = []isEqualTest{
	isEqualTest{
		name: "Test v1 not Porter type",
		v1:   "test",
		v2:   Integer(0),
		want: false,
	},
	isEqualTest{
		name: "Test v2 not Porter type",
		v1:   Integer(0),
		v2:   "test",
		want: false,
	},
	isEqualTest{
		name: "Test nil true",
		v1:   nil,
		v2:   nil,
		want: true,
	},
	isEqualTest{
		name: "Test nil false",
		v1:   nil,
		v2:   Float(1.23),
		want: false,
	},
	isEqualTest{
		name: "Test simple floats false",
		v1:   Float(1.23),
		v2:   Float(1.234),
		want: false,
	},
	isEqualTest{
		name: "Test simple floats true",
		v1:   Float(1.23),
		v2:   Float(1.23),
		want: true,
	},
	isEqualTest{
		name: "Test simple ints false",
		v1:   Integer(1),
		v2:   Integer(2),
		want: false,
	},
	isEqualTest{
		name: "Test simple ints true",
		v1:   Integer(5),
		v2:   Integer(5),
		want: true,
	},
	isEqualTest{
		name: "Test simple bools false",
		v1:   Boolean(false),
		v2:   Boolean(true),
		want: false,
	},
	isEqualTest{
		name: "Test simple bools true",
		v1:   Boolean(false),
		v2:   Boolean(false),
		want: true,
	},
	isEqualTest{
		name: "Test simple strings false",
		v1:   String("foo1"),
		v2:   String("foo2"),
		want: false,
	},
	isEqualTest{
		name: "Test simple strings true",
		v1:   String("foo"),
		v2:   String("foo"),
		want: true,
	},
	isEqualTest{
		name: "Test simple strings false",
		v1:   String("foo1"),
		v2:   String("foo2"),
		want: false,
	},
	isEqualTest{
		name: "Test simple strings true",
		v1:   String("foo"),
		v2:   String("foo"),
		want: true,
	},
	isEqualTest{
		name: "Test simple objects false",
		v1: Object{
			String("foo"): String("bar"),
		},
		v2:   Object{},
		want: false,
	},
	isEqualTest{
		name: "Test simple objects true",
		v1: Object{
			String("foo"): String("bar"),
		},
		v2: Object{
			String("foo"): String("bar"),
		},
		want: true,
	},
	isEqualTest{
		name: "Test simple objects false",
		v1: Array{
			String("foo"),
		},
		v2:   Array{},
		want: false,
	},
	isEqualTest{
		name: "Test simple arrays true",
		v1: Array{
			String("foo"),
		},
		v2: Array{
			String("foo"),
		},
		want: true,
	},
	isEqualTest{
		name: "Test object same length, different value type",
		v1: Object{
			String("foo"): String("bar"),
		},
		v2: Object{
			String("foo"): Integer(4),
		},
		want: false,
	},
	isEqualTest{
		name: "Test array same length, different value type",
		v1: Array{
			String("foo"),
		},
		v2: Array{
			Integer(4),
		},
		want: false,
	},
	isEqualTest{
		name: "Test simple type comparison array",
		v1: Array{
			String("foo"),
		},
		v2:   Integer(4),
		want: false,
	},
	isEqualTest{
		name: "Test simple type comparison object",
		v1: Object{
			String("bar"): String("foo"),
		},
		v2:   Integer(4),
		want: false,
	},
	isEqualTest{
		name: "Test simple type comparison literals",
		v1:   Float(5),
		v2:   Integer(4),
		want: false,
	},
	isEqualTest{
		name: "Test deep object and array comparison",
		v1: Object{
			String("foo"): Object{
				String("foo"): String("bar"),
			},
			String("bar"): Array{
				Integer(0),
				Integer(1),
			},
		},
		v2: Object{
			String("foo"): Object{
				String("foo"): String("bar"),
			},
			String("bar"): Array{
				Integer(0),
				Integer(1),
			},
		},
		want: true,
	},
}

func TestIsEqual(t *testing.T) {
	for _, c := range isEqualTests {
		got := IsEqual(c.v1, c.v2)

		if got != c.want {
			t.Errorf("Failed on: %s", c.name)
		}
	}
}

type getTest struct {
	name string
	v    Value
	path string
	want Value
}

var getTests = []getTest{
	getTest{
		name: "Empty string test",
		v:    Integer(1),
		path: "",
		want: Integer(1),
	},
	getTest{
		name: "Object nesting",
		v: Object{
			"foo": Object{
				"bar": Integer(1),
			},
		},
		path: "foo.bar",
		want: Integer(1),
	},
	getTest{
		name: "Object nesting with bracket syntax",
		v: Object{
			"foo": Object{
				"bar": Integer(1),
			},
		},
		path: "foo[bar]",
		want: Integer(1),
	},
	getTest{
		name: "Array indexing",
		v: Object{
			"foo": Array{
				Integer(0),
				Integer(1),
			},
		},
		path: "foo[1]",
		want: Integer(1),
	},
	getTest{
		name: "Complex object with mixed indexing",
		v: Object{
			"foo": Array{
				Object{
					"bar": Object{
						"foo1": Array{
							String("hello"),
							String("darkness"),
						},
					},
				},
			},
		},
		path: "foo[0].bar.foo1[1]",
		want: String("darkness"),
	},
	getTest{
		name: "Complex object with bracket indexing",
		v: Object{
			"foo": Array{
				Object{
					"bar": Object{
						"foo1": Array{
							String("hello"),
							String("darkness"),
						},
					},
				},
			},
		},
		path: "foo[0][bar][foo1][1]",
		want: String("darkness"),
	},
}

func TestGet(t *testing.T) {
	for _, c := range getTests {
		got, _ := Get(c.v, c.path)

		if !IsEqual(c.want, got) {
			t.Errorf("Failed on: %s", c.name)
		}
	}
}
