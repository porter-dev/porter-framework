package plan

import v "github.com/porterdev/ego/internal/value"

type (
	// OpType represents an enumerated operation type
	OpType int

	// Operation represents an executable operation that takes a Value from an
	// old state to a new state
	Operation struct {
		Op   OpType
		Path string
		Old  v.Value
		New  v.Value
	}
)

// The operation enumeration types: CREATE, READ, UPDATE, DELETE
const (
	CREATE OpType = iota
	READ
	UPDATE
	DELETE
)
