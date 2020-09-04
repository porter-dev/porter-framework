package plan

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

func TestIsOpEqual(t *testing.T) {
	op1 := Operation{
		Op:   UPDATE,
		Old:  v.Boolean(true),
		New:  v.Boolean(false),
		Path: "",
	}

	op2 := op1

	if !op1.IsEqual(&op2) {
		t.Errorf("Expected IsEqualOp(op1, op2) to be true, got false")
	}

	op3 := op1
	op3.Op = READ

	if op1.IsEqual(&op3) {
		t.Errorf("Expected IsEqualOp(op1, op3) to be false, got true")
	}

	op4 := op1
	op4.Path = "Hello"

	if op1.IsEqual(&op4) {
		t.Errorf("Expected IsEqualOp(op1, op4) to be false, got true")
	}

	op5 := op1
	op5.New = v.Boolean(true)

	if op1.IsEqual(&op5) {
		t.Errorf("Expected IsEqualOp(op1, op5) to be false, got true")
	}
}
