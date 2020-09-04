package plan

import (
	"strconv"

	v "github.com/porterdev/ego/internal/value"
)

// IsEqual compares two operations to determine if they are equal.
func (op1 *Operation) IsEqual(op2 *Operation) bool {
	// check enumeration is equal
	res := op1.Op == op2.Op

	// check old and new are equal
	res = res && v.IsEqual(op1.Old, op2.Old) && v.IsEqual(op1.New, op2.New)

	// check paths are equal
	res = res && op1.Path == op2.Path

	return res
}

// ToString converts an operation to a string based on the operation type and the
// path. For example, an UPDATE operation at the path [example] becomes
// "2:[example]"
func (op1 *Operation) ToString() string {
	return strconv.Itoa(int(op1.Op)) + ":" + op1.Path
}
