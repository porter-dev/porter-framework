package translator

import (
	"testing"
)

type transTest struct {
	in, out string
}

var testToStrings = []transTest{
	{
		in:  "<<test\nfoo:{{bar.foo}}\ntest>>",
		out: "`\nfoo:` + bar.foo + `\n`",
	},
	{
		in:  "var test1 int = 2;\nvar test2 int = 2",
		out: "var test1 int = 2;\nvar test2 int = 2",
	},
	{
		in:  "var test1 string = \"foo\";\nvar test2 string = \"bar\";\n<<test\n{{test1}}:{{test2}}\ntest>>",
		out: "var test1 string = \"foo\";\nvar test2 string = \"bar\";\n`\n` + test1 + `:` + test2 + `\n`",
	},
}

func TestTranslateToString(t *testing.T) {
	for _, c := range testToStrings {
		bytes := []byte(c.in)

		trans := NewTranslator(bytes)

		res := trans.TranslateToString()

		if string(res) != c.out {
			t.Errorf("(%s).Translate() expected %s, got %s", c.in, c.out, res)
		}
	}
}

var testToJSONs = []transTest{
	{
		in:  "<<json\n{\nfoo:{{bar.foo}}\n}\njson>>",
		out: "json.Inject(`\n{\nfoo:{{bar.foo}}\n}\n`,bar.foo)",
	},
}

func TestTranslateToJSON(t *testing.T) {
	for _, c := range testToJSONs {
		bytes := []byte(c.in)

		trans := NewTranslator(bytes)

		res := trans.TranslateToJSON()

		if string(res) != c.out {
			t.Errorf("(%s).Translate() expected %s, got %s", c.in, c.out, res)
		}
	}
}
