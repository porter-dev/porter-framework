package plan

type (
	// OpStack represents a stack of operations to be executed using
	// LIFO removal
	OpStack struct {
		top    *stackNode
		length int
	}

	stackNode struct {
		value *Operation
		prev  *stackNode
	}
)

// NewOpStack constructs an empty OpStack
func NewOpStack() *OpStack {
	return &OpStack{nil, 0}
}

// IsEmpty returns true if a OpStack is empty, false otherwise
func (s *OpStack) IsEmpty() bool {
	return s.length == 0
}

// Len returns the number of items in the OpStack
func (s *OpStack) Len() int {
	return s.length
}

// Peek views the top item on the OpStack. If the OpStack
// is empty, Peek returns nil.
func (s *OpStack) Peek() *Operation {
	if s.IsEmpty() {
		return nil
	}

	return s.top.value
}

// Pop removes and returns the top item of the OpStack. If the OpStack
// is empty, Pop returns nil.
func (s *OpStack) Pop() *Operation {
	if s.IsEmpty() {
		return nil
	}

	n := s.top
	s.top = n.prev
	s.length--
	return n.value
}

// Push adds an operation to the top of the OpStack.
func (s *OpStack) Push(op *Operation) {
	n := &stackNode{op, s.top}
	s.top = n
	s.length++
}
