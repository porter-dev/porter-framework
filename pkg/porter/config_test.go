package porter

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

func TestSimpleConfigJustInput(t *testing.T) {
	conf := CreateDefaultConfig("12345", "./", "./", 2)

	input := v.Object{
		v.String("hello"): v.String("there"),
	}

	conf.Apply(input)
}
