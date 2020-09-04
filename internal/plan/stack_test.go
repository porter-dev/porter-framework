package plan

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

func TestStackEmpty(t *testing.T) {
	// create an empty stack
	stack := NewOpStack()

	// check that it's empty and that read operations return <nil>
	if isEmpty := stack.IsEmpty(); !isEmpty {
		t.Errorf("Expected stack.IsEmpty to be true, was %v", isEmpty)
	}

	if peek := stack.Peek(); peek != nil {
		t.Errorf("Expected stack.Peek to return <nil>, was %v", peek)
	}

	if pop := stack.Pop(); pop != nil {
		t.Errorf("Expected stack.Pop to return <nil>, was %v", pop)
	}
}

// initStack is a helper method to create a simple OpStack.
func initStack(s *OpStack) (*Operation, *Operation) {
	op1 := &Operation{
		Op:   CREATE,
		Path: "",
		Old:  v.String("hello"),
		New:  v.Value("there"),
	}

	s.Push(op1)

	op2 := &Operation{
		Op:   READ,
		Path: ".hello",
		Old:  v.String("beep"),
		New:  v.Value("beep"),
	}

	s.Push(op2)

	return op1, op2
}

func TestStack(t *testing.T) {
	s := NewOpStack()

	op1, op2 := initStack(s)

	val1 := s.Peek()

	// should reference the exact same struct
	if *op2 != *val1 {
		t.Errorf("Expected stack.Peek to return %v, was %v", *op2, *val1)
	}

	val2 := s.Pop()

	// should reference the exact same struct
	if *op2 != *val2 {
		t.Errorf("Expected stack.Pop to return %v, was %v", *op2, *val2)
	}

	// length should now be one
	if len := s.Len(); len != 1 {
		t.Errorf("Expected stack.Len to return 1, was %v", len)
	}

	val3 := s.Peek()

	// should now be op1
	if *op1 != *val3 {
		t.Errorf("Expected stack.Peek to return %v, was %v", *op1, *val3)
	}

	val4 := s.Pop()

	// should also be op1
	if *op1 != *val4 {
		t.Errorf("Expected stack.Peek to return %v, was %v", *op1, *val4)
	}

	if len := s.Len(); len != 0 {
		t.Errorf("Expected stack.Len to return 0, was %v", len)
	}

}
