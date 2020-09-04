package plan

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

func TestQueueEmpty(t *testing.T) {
	// create an empty queue
	queue := NewOpQueue()

	// check that it's empty and that read operations return <nil>
	if isEmpty := queue.IsEmpty(); !isEmpty {
		t.Errorf("Expected queue.IsEmpty to be true, was %v", isEmpty)
	}

	if peek := queue.Peek(); peek != nil {
		t.Errorf("Expected queue.Peek to return <nil>, was %v", peek)
	}

	if deq := queue.Dequeue(); deq != nil {
		t.Errorf("Expected queue.Dequeue to return <nil>, was %v", deq)
	}
}

// initQueue is a helper method to create a simple OpQueue.
func initQueue(q *OpQueue) (*Operation, *Operation) {
	op1 := &Operation{
		Op:   CREATE,
		Path: "",
		Old:  v.String("hello"),
		New:  v.Value("there"),
	}

	q.Enqueue(op1)

	op2 := &Operation{
		Op:   READ,
		Path: ".hello",
		Old:  v.String("beep"),
		New:  v.Value("beep"),
	}

	q.Enqueue(op2)

	return op1, op2
}

func TestQueue(t *testing.T) {
	q := NewOpQueue()

	op1, op2 := initQueue(q)

	val1 := q.Peek()

	// length should be 2
	if len := q.Len(); len != 2 {
		t.Errorf("Expected queue.Len to return 2, was %v", len)
	}

	// front should be op1
	if *op1 != *val1 {
		t.Errorf("Expected queue.Peek to return %v, was %v", *op1, *val1)
	}

	val2 := q.Dequeue()

	// length should be 1
	if len := q.Len(); len != 1 {
		t.Errorf("Expected queue.Len to return 1, was %v", len)
	}

	// front should be op1
	if *op1 != *val2 {
		t.Errorf("Expected queue.Peek to return %v, was %v", *op1, *val2)
	}

	val3 := q.Dequeue()

	// length should be 0
	if len := q.Len(); len != 0 {
		t.Errorf("Expected queue.Len to return 0, was %v", len)
	}

	// front should now be op2
	if *op2 != *val3 {
		t.Errorf("Expected queue.Peek to return %v, was %v", *op2, *val3)
	}
}
