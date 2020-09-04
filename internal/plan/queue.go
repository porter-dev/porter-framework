package plan

type (
	// OpQueue represents a queue of operations to be executed using
	// FIFO removal
	OpQueue struct {
		front  *node
		rear   *node
		length int
	}

	node struct {
		value *Operation
		prev  *node
		next  *node
	}
)

// NewOpQueue constructs an empty OpQueue
func NewOpQueue() *OpQueue {
	return &OpQueue{nil, nil, 0}
}

// IsEmpty returns true if an OpQueue is empty, false otherwise
func (q *OpQueue) IsEmpty() bool {
	return q.length == 0
}

// Len returns the number of items in the OpQueue
func (q *OpQueue) Len() int {
	return q.length
}

// Enqueue adds an operation to the rear of the queue
func (q *OpQueue) Enqueue(op *Operation) {
	// if no elements in queue, add as front and rear of queue
	if q.IsEmpty() {
		node := &node{op, nil, nil}
		q.front = node
		q.rear = node
		q.length++
		return
	}

	// else, get the current rear of queue and set to next
	next := q.rear
	node := &node{op, nil, next}

	// overwrite the current rear of queue and next.prev
	q.rear = node
	next.prev = node
	q.length++
	return
}

// Dequeue removes and returns the front operation from queue
func (q *OpQueue) Dequeue() (op *Operation) {
	// if no elements in queue, return nil
	if q.IsEmpty() {
		return nil
	}

	// if length is 1, make queue empty and return single node value
	if q.length == 1 {
		node := q.rear
		q.rear = nil
		q.front = nil
		q.length--
		return node.value
	}

	// else, get the current front of queue and rewrite
	node := q.front
	node.prev.next = nil
	q.front = node.prev
	q.length--
	return node.value
}

// Peek returns the element at the front
func (q *OpQueue) Peek() (op *Operation) {
	// if no elements in queue, return nil
	if q.IsEmpty() {
		return nil
	}

	res := q.front.value
	return res
}
