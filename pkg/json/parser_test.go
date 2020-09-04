package json

// NOTE: this is not a "pure parser". It is meant to be converted
// to Porter configuration objects, which have a tighter set of constraints
// than JSON (although the specs are very similar). But notably, it does not
// support:
// - numbers that cause an overflow
// - strings with null characters in them
// - does not accept surrogate code points
// - does not accept escape sequence \/
// - accepts single spaces, because...why not?
//
// The distinction between a JSON Number and a Float/Integer type should also
// be noted.

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

type jsonTest struct {
	name string
	json string
	want v.Value
}

var jsonTestsArrayPass = []jsonTest{
	jsonTest{
		name: "Array: arrays with spaces",
		json: "[[]   ]",
		want: v.Array{
			v.Array{},
		},
	},
	jsonTest{
		name: "Array: empty string",
		json: "[\"\"]",
		want: v.Array{
			v.String(""),
		},
	},
	jsonTest{
		name: "Array: empty array",
		json: "[]",
		want: v.Array{},
	},
	jsonTest{
		name: "Array: array ending with newline",
		json: "[\"a\"]\n",
		want: v.Array{
			v.String("a"),
		},
	},
	jsonTest{
		name: "Array: array false",
		json: "[false]",
		want: v.Array{
			v.Boolean(false),
		},
	},
	jsonTest{
		name: "Array: array heterogeneous",
		json: "[null, 1, \"1\", {}]",
		want: v.Array{
			nil,
			v.Integer(1),
			v.String("1"),
			v.Object{},
		},
	},
	jsonTest{
		name: "Array: array null",
		json: "[null]",
		want: v.Array{
			nil,
		},
	},
	jsonTest{
		name: "Array: array with 1 and newline",
		json: "[1\n]",
		want: v.Array{
			v.Integer(1),
		},
	},
	jsonTest{
		name: "Array: array with leading space",
		json: " [1]",
		want: v.Array{
			v.Integer(1),
		},
	},
	jsonTest{
		name: "Array: array with several null",
		json: "[1,null,null,null,2]",
		want: v.Array{
			v.Integer(1),
			nil,
			nil,
			nil,
			v.Integer(2),
		},
	},
	jsonTest{
		name: "Array: array with trailing space",
		json: "[2] ",
		want: v.Array{
			v.Integer(2),
		},
	},
}

func TestJSONPassArray(t *testing.T) {
	for _, c := range jsonTestsArrayPass {
		src := []byte(c.json)
		p := NewParser(src)

		res, _ := p.Parse()

		if !v.IsEqual(res, c.want) {
			t.Errorf("Failed on: %s, %v, %v", c.name, res, c.want)
		}
	}
}

var jsonTestsNumberPass = []jsonTest{
	jsonTest{
		name: "Number: basic number",
		json: "[123e4]",
		want: v.Array{
			v.Integer(123e4),
		},
	},
	jsonTest{
		name: "Number: 0e1",
		json: "[0e1]",
		want: v.Array{
			v.Integer(0),
		},
	},
	jsonTest{
		name: "Number: number after space",
		json: "[ 4]",
		want: v.Array{
			v.Integer(4),
		},
	},
	jsonTest{
		name: "Number: number double close to zero",
		json: "[-0.000000000000000000000000000000000000000000000000000000000000000000000000000001]",
		want: v.Array{
			v.Float(-1e-78),
		},
	},
	jsonTest{
		name: "Number: number int with exponent",
		json: "[20e1]",
		want: v.Array{
			v.Integer(200),
		},
	},
	jsonTest{
		name: "Number: number negative 0",
		json: "[-0]",
		want: v.Array{
			v.Integer(0),
		},
	},
	jsonTest{
		name: "Number: number negative int",
		json: "[-123]",
		want: v.Array{
			v.Integer(-123),
		},
	},
	jsonTest{
		name: "Number: number negative one",
		json: "[-1]",
		want: v.Array{
			v.Integer(-1),
		},
	},
	jsonTest{
		name: "Number: number capital e",
		json: "[1E10]",
		want: v.Array{
			v.Integer(1e10),
		},
	},
	jsonTest{
		name: "Number: number capital e negative exponent",
		json: "[1E-2]",
		want: v.Array{
			v.Float(1e-2),
		},
	},
	jsonTest{
		name: "Number: number capital e positive exponent",
		json: "[1E+2]",
		want: v.Array{
			v.Integer(1e+2),
		},
	},
	jsonTest{
		name: "Number: number real fraction with exponent",
		json: "[123.4e1]",
		want: v.Array{
			v.Integer(123.4 * 10),
		},
	},
	jsonTest{
		name: "Number: number real negative exponent",
		json: "[1e-2]",
		want: v.Array{
			v.Float(1e-2),
		},
	},
	jsonTest{
		name: "Number: number simple int",
		json: "[123]",
		want: v.Array{
			v.Integer(123),
		},
	},
	jsonTest{
		name: "Number: number simple float",
		json: "[123.456789]",
		want: v.Array{
			v.Float(123.456789),
		},
	},
}

func TestJSONPassNumber(t *testing.T) {
	for _, c := range jsonTestsNumberPass {
		src := []byte(c.json)
		p := NewParser(src)

		res, _ := p.Parse()

		if !v.IsEqual(res, c.want) {
			t.Errorf("Failed on: %s, %v, %v", c.name, res, c.want)
		}
	}
}

var jsonTestsObjectPass = []jsonTest{
	jsonTest{
		name: "Object: basic object",
		json: "{\"asd\":\"sdf\"}",
		want: v.Object{
			v.String("asd"): v.String("sdf"),
		},
	},
	jsonTest{
		name: "Object: two keys object",
		json: "{\"asd\":\"sdf\", \"dfg\":\"fgh\"}",
		want: v.Object{
			v.String("asd"): v.String("sdf"),
			v.String("dfg"): v.String("fgh"),
		},
	},
	jsonTest{
		name: "Object: duplicate keys",
		json: "{\"a\":\"b\", \"a\":\"c\"}",
		want: v.Object{
			v.String("a"): v.String("c"),
		},
	},
	jsonTest{
		name: "Object: duplicate key and value",
		json: "{\"a\":\"b\", \"a\":\"b\"}",
		want: v.Object{
			v.String("a"): v.String("b"),
		},
	},
	jsonTest{
		name: "Object: empty",
		json: "{}",
		want: v.Object{},
	},
	jsonTest{
		name: "Object: empty",
		json: "{\"\":0}",
		want: v.Object{
			v.String(""): v.Integer(0),
		},
	},
	jsonTest{
		name: "Object: empty",
		json: "{\"\":0}",
		want: v.Object{
			v.String(""): v.Integer(0),
		},
	},
	jsonTest{
		name: "Object: long strings",
		json: "{\"x\":[{\"id\": \"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\"}], \"id\": \"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\"}",
		want: v.Object{
			v.String("x"): v.Array{
				v.Object{
					v.String("id"): v.String("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
				},
			},
			v.String("id"): v.String("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		},
	},
	jsonTest{
		name: "Object: simple array",
		json: "{\"title\":\"\u041f\u043e\u043b\u0442\u043e\u0440\u0430 \u0417\u0435\u043c\u043b\u0435\u043a\u043e\u043f\u0430\" }",
		want: v.Object{
			v.String("title"): v.String("Полтора Землекопа"),
		},
	},
	jsonTest{
		name: "Object: with newlines",
		json: "{\n\"a\": \"b\"\n}",
		want: v.Object{
			v.String("a"): v.String("b"),
		},
	},
}

func TestJSONPassObject(t *testing.T) {
	for _, c := range jsonTestsObjectPass {
		src := []byte(c.json)
		p := NewParser(src)

		res, _ := p.Parse()

		if !v.IsEqual(res, c.want) {
			t.Errorf("Failed on: %s, %v, %v", c.name, res, c.want)
		}
	}
}

var jsonTestsStringPass = []jsonTest{
	jsonTest{
		name: "String: utf8 sequences",
		json: "[\"\u0060\u012a\u12AB\"]",
		want: v.Array{
			v.String("\u0060\u012a\u12AB"),
		},
	},
	jsonTest{
		name: "String: double escape",
		json: "[\"\\a\"]",
		want: v.Array{
			v.String("\\a"),
		},
	},
	jsonTest{
		name: "String: basic in array",
		json: "[\"hello\"]",
		want: v.Array{
			v.String("hello"),
		},
	},
}

func TestJSONPassString(t *testing.T) {
	for _, c := range jsonTestsStringPass {
		src := []byte(c.json)
		p := NewParser(src)

		res, _ := p.Parse()

		if !v.IsEqual(res, c.want) {
			t.Errorf("Failed on: %s, %v, %v", c.name, res, c.want)
		}
	}
}

var jsonTestsStructurePass = []jsonTest{
	jsonTest{
		name: "Structure: lonely false",
		json: "false",
		want: v.Boolean(false),
	},
	jsonTest{
		name: "Structure: lonely int",
		json: "42",
		want: v.Integer(42),
	},
	jsonTest{
		name: "Structure: lonely negative real",
		json: "-0.1",
		want: v.Float(-0.1),
	},
	jsonTest{
		name: "Structure: lonely null",
		json: "null",
		want: nil,
	},
	jsonTest{
		name: "Structure: lonely string",
		json: "\"asd\"",
		want: v.String("asd"),
	},
	jsonTest{
		name: "Structure: lonely true",
		json: "true",
		want: v.Boolean(true),
	},
	jsonTest{
		name: "Structure: empty string",
		json: "\"\"",
		want: v.String(""),
	},
	jsonTest{
		name: "Structure: trailing newline",
		json: "[\"a\"]\n",
		want: v.Array{
			v.String("a"),
		},
	},
	jsonTest{
		name: "Structure: whitespace array",
		json: " [] 	",
		want: v.Array{},
	},
}

func TestJSONPassStructure(t *testing.T) {
	for _, c := range jsonTestsStructurePass {
		src := []byte(c.json)
		p := NewParser(src)

		res, _ := p.Parse()

		if !v.IsEqual(res, c.want) {
			t.Errorf("Failed on: %s, %v, %v", c.name, res, c.want)
		}
	}
}

type jsonTestFail struct {
	name string
	json string
	msg  string
}

var jsonTestsArrayFail = []jsonTestFail{
	jsonTestFail{
		name: "Array: 1 true without comma",
		json: "[1 true]",
		msg:  "Column 8, Line 1: Values must be separated by a comma",
	},
	jsonTestFail{
		name: "Array: invalid UTF-8",
		json: "[a�]",
		msg:  "Column 3, Line 1: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
	},
	jsonTestFail{
		name: "Array: colon instead of comma",
		json: "[\"\": 1]",
		msg:  "Column 5, Line 1: Values must be separated by a comma",
	},
	jsonTestFail{
		name: "Array: colon instead of comma",
		json: "[\"\"],",
		msg:  "Column 6, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Array: comma then number",
		json: "[,1]",
		msg:  "Column 3, Line 1: Must have a value between array elements",
	},
	jsonTestFail{
		name: "Array: double comma",
		json: "[1,,2]",
		msg:  "Column 5, Line 1: Must have a value between array elements",
	},
	jsonTestFail{
		name: "Array: extra closing bracket",
		json: "[\"x\"]]",
		msg:  "Column 7, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Array: extra comma",
		json: "[\"\",]",
		msg:  "Column 6, Line 1: Commas must be followed by a value",
	},
	jsonTestFail{
		name: "Array: no closing bracket",
		json: "[\"x\"",
		msg:  "Column 5, Line 1: No closing bracket in array",
	},
	jsonTestFail{
		name: "Array: incomplete, invalid value",
		json: "[x",
		msg:  "Column 3, Line 1: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
	},
	jsonTestFail{
		name: "Array: inner array no comma",
		json: "[3[4]]",
		msg:  "Column 4, Line 1: Values must be separated by a comma",
	},
	jsonTestFail{
		name: "Array: array just comma",
		json: "[,]",
		msg:  "Column 3, Line 1: Must have a value between array elements",
	},
	jsonTestFail{
		name: "Array: array just minus",
		json: "[-]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Array: unclosed with newlines",
		json: "[\"a\",\n4\n,1,",
		msg:  "Column 4, Line 3: No closing bracket in array",
	},
	jsonTestFail{
		name: "Array: spaces vertical tab formfeed",
		json: "[\"a\"\f]",
		msg: "Column 6, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Array: star inside",
		json: "[*]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Array: unclosed trailing comma",
		json: "[1,",
		msg:  "Column 4, Line 1: No closing bracket in array",
	},
	jsonTestFail{
		name: "Array: unclosed with object inside",
		json: "[{}",
		msg:  "Column 4, Line 1: No closing bracket in array",
	},
	jsonTestFail{
		name: "Array: array with broken subarray",
		json: "[[x]]",
		msg:  "Column 4, Line 1: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
	},
	jsonTestFail{
		name: "Array: no closing bracket",
		json: "]",
		msg:  "Column 2, Line 1: Right bracket ] not preceded by left bracket [",
	},
}

func TestJSONFailArray(t *testing.T) {
	for _, c := range jsonTestsArrayFail {
		src := []byte(c.json)
		p := NewParser(src)

		_, err := p.Parse()

		if err == nil || err.Error() != c.msg {
			t.Errorf("Failed on: %s, input %v, expected %v, got %v", c.name, c.json, c.msg, err)
		}
	}
}

var jsonTestsIncompleteFail = []jsonTestFail{
	jsonTestFail{
		name: "Incomplete: true",
		json: "[tru]",
		msg:  "Column 5, Line 1: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
	},
	jsonTestFail{
		name: "Incomplete: null",
		json: "[nul]",
		msg:  "Column 5, Line 1: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
	},
	jsonTestFail{
		name: "Incomplete: false",
		json: "[fals]",
		msg:  "Column 6, Line 1: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
	},
}

func TestJSONFailIncomplete(t *testing.T) {
	for _, c := range jsonTestsIncompleteFail {
		src := []byte(c.json)
		p := NewParser(src)

		_, err := p.Parse()

		if err == nil || err.Error() != c.msg {
			t.Errorf("Failed on: %s, input %v, expected %v, got %v", c.name, c.json, c.msg, err)
		}
	}
}

var jsonTestsNumberFail = []jsonTestFail{
	jsonTestFail{
		name: "Number: ++",
		json: "[++1234]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: +1",
		json: "[+1]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: +Inf",
		json: "[+Inf]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: -01",
		json: "[-01]",
		msg:  "Column 3, Line 1: 0 cannot be followed by digit without . or exponent",
	},
	jsonTestFail{
		name: "Number: illegal period placed",
		json: "[-1.0.]",
		msg:  "Column 6, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: period without following decimal",
		json: "[-2.]",
		msg:  "Column 5, Line 1: . must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: NaN",
		json: "[-NaN]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: .-1",
		json: "[.-1]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: .2e-3",
		json: "[.2e-3]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: 0.1.2",
		json: "[0.1.2]",
		msg:  "Column 5, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: 0.3e+",
		json: "[0.3e+]",
		msg:  "Column 7, Line 1: e or E must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: 0.3e",
		json: "[0.3e]",
		msg:  "Column 6, Line 1: e or E must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: 0.e1",
		json: "[0.e1]",
		msg:  "Column 4, Line 1: . must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: 0E+",
		json: "[0E+]",
		msg:  "Column 5, Line 1: e or E must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: 0E",
		json: "[0E]",
		msg:  "Column 4, Line 1: e or E must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: 1.0e-",
		json: "[1.0e-]",
		msg:  "Column 7, Line 1: e or E must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: 1 000.0",
		json: "[1 000.0]",
		msg:  "Column 4, Line 1: 0 cannot be followed by digit without . or exponent",
	},
	jsonTestFail{
		name: "Number: fullwidth digit one",
		json: "[１]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: expression",
		json: "[1+2]",
		msg:  "Column 3, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: +- after exponent",
		json: "[0e+-1]",
		msg:  "Column 5, Line 1: e or E must be followed by at least one digit",
	},
	jsonTestFail{
		name: "Number: invalid negative real",
		json: "[-123.123foo]",
		msg:  "Column 13, Line 1: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
	},
	jsonTestFail{
		name: "Number: minus space number",
		json: "[- 1]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: negative int starting with zero",
		json: "[-012]",
		msg:  "Column 3, Line 1: 0 cannot be followed by digit without . or exponent",
	},
	jsonTestFail{
		name: "Number: negative real with int part",
		json: "[-.123]",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Number: leading zero",
		json: "[012]",
		msg:  "Column 2, Line 1: 0 cannot be followed by digit without . or exponent",
	},
	jsonTestFail{
		name: "Number: int with overflow",
		json: "[10000000000000100000000000001000000000000010000000000000100000000000001000000000000010000000000000100000000000001000000000000010000000000000100000000000001000000000000010000000000000e128]",
		msg:  "Column 188, Line 1: Unable to parse number into int or float",
	},
}

func TestJSONFailNumber(t *testing.T) {
	for _, c := range jsonTestsNumberFail {
		src := []byte(c.json)
		p := NewParser(src)

		_, err := p.Parse()

		if err == nil || err.Error() != c.msg {
			t.Errorf("Failed on: %s, input %v, expected %v, got %v", c.name, c.json, c.msg, err)
		}
	}
}

var jsonTestsObjectFail = []jsonTestFail{
	jsonTestFail{
		name: "Object: bracket as key",
		json: "{[: \"x\"}",
		msg:  "Column 3, Line 1: Left brace { not followed by comma-separated string: value pairs or right brace }",
	},
	jsonTestFail{
		name: "Object: comma instead of colon",
		json: "{\"x\", null}",
		msg:  "Column 6, Line 1: Must use colon : to define a string: value pair",
	},
	jsonTestFail{
		name: "Object: double colon",
		json: "{\"x\"::\"b\"}",
		msg:  "Column 7, Line 1: Colons can only be used for string: value pairs inside objects",
	},
	jsonTestFail{
		name: "Object: emoji",
		json: "{🇨🇭}",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Object: garbage at end",
		json: "{\"a\":\"a\" 123}",
		msg:  "Column 13, Line 1: Left brace { not followed by comma-separated string: value pairs or right brace }",
	},
	jsonTestFail{
		name: "Object: key with single quotes",
		json: "{'key': \"value\"}",
		msg:  "Column 2, Line 1: Illegal token",
	},
	jsonTestFail{
		name: "Object: trailing comma",
		json: "{\"0\":\"0\",}",
		msg:  "Column 11, Line 1: Commas must be followed by a string: value pair",
	},
	jsonTestFail{
		name: "Object: missing colon",
		json: "{\"a\" 1}",
		msg:  "Column 7, Line 1: Must use colon : to define a string: value pair",
	},
	jsonTestFail{
		name: "Object: missing key",
		json: "{:1}",
		msg:  "Column 3, Line 1: Left brace { not followed by comma-separated string: value pairs or right brace }",
	},
	jsonTestFail{
		name: "Object: missing value",
		json: "{\"a\":",
		msg:  "Column 6, Line 1: Left brace { not followed by comma-separated string: value pairs or right brace }",
	},
	jsonTestFail{
		name: "Object: missing colon then EOF",
		json: "{\"a\"",
		msg:  "Column 5, Line 1: Must use colon : to define a string: value pair",
	},
	jsonTestFail{
		name: "Object: non-string key",
		json: "{1:1}",
		msg:  "Column 3, Line 1: Left brace { not followed by comma-separated string: value pairs or right brace }",
	},
	jsonTestFail{
		name: "Object: null key",
		json: "{null:null,null:null}",
		msg:  "Column 6, Line 1: Left brace { not followed by comma-separated string: value pairs or right brace }",
	},
	jsonTestFail{
		name: "Object: several trailing commas",
		json: "{\"id\":0,,,,,}",
		msg:  "Column 10, Line 1: Commas must be followed by a string: value pair",
	},
	jsonTestFail{
		name: "Object: single trailing comma",
		json: "{\"id\":0,}",
		msg:  "Column 10, Line 1: Commas must be followed by a string: value pair",
	},
	jsonTestFail{
		name: "Object: two commas in a row",
		json: "{\"a\":\"b\",,\"c\":\"d\"}",
		msg:  "Column 11, Line 1: Commas must be followed by a string: value pair",
	},
	jsonTestFail{
		name: "Object: unterminated value",
		json: "{\"a\":\"a",
		msg:  "Column 8, Line 1: String not terminated",
	},
	jsonTestFail{
		name: "Object: single string",
		json: "{ \"foo\" : \"bar\", \"a\" }",
		msg:  "Column 23, Line 1: Must use colon : to define a string: value pair",
	},
	jsonTestFail{
		name: "Object: single ending bracket",
		json: "}",
		msg:  "Column 2, Line 1: Right brace } not preceded by left brace {",
	},
}

func TestJSONFailObject(t *testing.T) {
	for _, c := range jsonTestsObjectFail {
		src := []byte(c.json)
		p := NewParser(src)

		_, err := p.Parse()

		if err == nil || err.Error() != c.msg {
			t.Errorf("Failed on: %s, input %v, expected %v, got %v", c.name, c.json, c.msg, err)
		}
	}
}

var jsonTestsStructureFail = []jsonTestFail{
	jsonTestFail{
		name: "Structure: single comma",
		json: ",",
		msg:  "Column 2, Line 1: Commas can only be present when separating array or object values",
	},
	jsonTestFail{
		name: "Structure: single colon",
		json: ":",
		msg:  "Column 2, Line 1: Colons can only be used for string: value pairs inside objects",
	},
}

func TestJSONFailStructure(t *testing.T) {
	for _, c := range jsonTestsStructureFail {
		src := []byte(c.json)
		p := NewParser(src)

		_, err := p.Parse()

		if err == nil || err.Error() != c.msg {
			t.Errorf("Failed on: %s, input %v, expected %v, got %v", c.name, c.json, c.msg, err)
		}
	}
}

type injectTest struct {
	name string
	json string
	vals v.Array
	want v.Value
}

type injectTestFail struct {
	name string
	json string
	vals v.Array
	msg  string
}

var injectTestsPass = []injectTest{
	injectTest{
		name: "Inject: value test",
		json: "{ \"id\": {{hello}} }",
		vals: v.Array{
			v.Object{
				"foo": v.String("bar"),
			},
		},
		want: v.Object{
			"id": v.Object{
				"foo": v.String("bar"),
			},
		},
	},
	injectTest{
		name: "Inject: key and  value test",
		json: "{ {{test1}}: {{test2}} }",
		vals: []v.Value{
			v.String("foo"),
			v.String("bar"),
		},
		want: v.Object{
			v.String("foo"): v.String("bar"),
		},
	},
}

func TestInjectPass(t *testing.T) {
	for _, c := range injectTestsPass {
		src := c.json

		res, _ := Inject(src, c.vals...)

		if !v.IsEqual(res, c.want) {
			t.Errorf("Failed on: %s, %v, %v", c.name, res, c.want)
		}
	}
}

var injectTestsFail = []injectTestFail{
	injectTestFail{
		name: "Inject: fail",
		json: "{{a ",
		vals: v.Array{
			v.Object{
				"foo": v.String("bar"),
			},
		},
		msg: "Column 5, Line 1: No closing injection }} found",
	},
}

func TestInjectFail(t *testing.T) {
	for _, c := range injectTestsFail {
		src := c.json

		_, err := Inject(src, c.vals...)

		if err == nil || err.Error() != c.msg {
			t.Errorf("Failed on: %s, input %v, expected %v, got %v", c.name, c.json, c.msg, err)
		}
	}
}
