package plan

import (
	"math"
	"strconv"

	v "github.com/porterdev/ego/internal/value"
)

// CreateOpQueue takes the current config (old) and generates a sequence of operations
// to create the desired config (new).
func CreateOpQueue(old v.Value, new v.Value, paths ...string) *OpQueue {
	q := NewOpQueue()
	genOperationRecursive(old, new, q, "", paths...)
	return q
}

func contains(paths []string, path string) bool {
	for _, _path := range paths {
		if _path == path {
			return true
		}
	}

	return false
}

func genOperationRecursive(old v.Value, new v.Value, q *OpQueue, prefix string, paths ...string) {
	if old == nil && new == nil {
		return
	} else if old == nil && new != nil {
		q.Enqueue(&Operation{
			Op:   CREATE,
			Path: prefix,
			Old:  old,
			New:  new})
		return
	} else if old != nil && new == nil {
		q.Enqueue(&Operation{
			Op:   DELETE,
			Path: prefix,
			Old:  old,
			New:  new})
		return
	}

	if contains(paths, prefix) {
		if !v.IsEqual(old, new) {
			q.Enqueue(&Operation{
				Op:   UPDATE,
				Path: prefix,
				Old:  old,
				New:  new})

			return
		}

		q.Enqueue(&Operation{
			Op:   READ,
			Path: prefix,
			Old:  old,
			New:  new})

		return
	}

	switch old.(type) {
	case v.Boolean:
		_, ok := new.(v.Boolean)
		enqueuePrimitiveOp(ok, old, new, q, prefix)
	case v.Float:
		_, ok := new.(v.Float)
		enqueuePrimitiveOp(ok, old, new, q, prefix)
	case v.Integer:
		_, ok := new.(v.Integer)
		enqueuePrimitiveOp(ok, old, new, q, prefix)
	case v.String:
		_, ok := new.(v.String)
		enqueuePrimitiveOp(ok, old, new, q, prefix)
	case v.Array:
		oldArr, _ := old.(v.Array)
		newArr, ok := new.(v.Array)

		if !ok {
			// If types are different, can treat this as a "primitive" operation that just
			// gets written as an UPDATE operation.
			enqueuePrimitiveOp(ok, old, new, q, prefix)
			return
		}

		oldLen := len(oldArr)
		newLen := len(newArr)

		// iterate through all shared indices, which get queued as UPDATE operations
		for i := 0.; i < math.Min(float64(oldLen), float64(newLen)); i++ {
			_i := int(i)
			genOperationRecursive(oldArr[_i], newArr[_i], q, prefix+"["+strconv.Itoa(_i)+"]", paths...)
		}

		if oldLen > newLen {
			// all old unshared indices become DELETE operations
			for i := newLen; i < oldLen; i++ {
				q.Enqueue(&Operation{
					Op:   DELETE,
					Path: prefix + "[" + strconv.Itoa(i) + "]",
					Old:  oldArr[i],
					New:  nil})
			}
		} else if oldLen < newLen {
			// all new unshared indices become CREATE operations
			for i := oldLen; i < newLen; i++ {
				q.Enqueue(&Operation{
					Op:   CREATE,
					Path: prefix + "[" + strconv.Itoa(i) + "]",
					Old:  nil,
					New:  newArr[i]})
			}
		}
	case v.Object:
		oldObj, _ := old.(v.Object)
		newObj, ok := new.(v.Object)

		if !ok {
			// If types are different, can treat this as a "primitive" operation that just
			// gets written as an UPDATE operation.
			enqueuePrimitiveOp(ok, old, new, q, prefix)
			return
		}

		// iterate through keys to discover which keys are shared, and recurse on each
		// shared key
		oldKeys := make(map[v.String]int)
		newKeys := make(map[v.String]int)

		// add each newObj key to newKeys
		for k := range newObj {
			newKeys[k] = 0
		}

		for k := range oldObj {
			// determine if the key exists in the new object
			if _, found := newKeys[k]; found {
				// remove from newKeys
				delete(newKeys, k)

				// recurse with both child objects
				genOperationRecursive(oldObj[k], newObj[k], q, prefix+"["+string(k)+"]", paths...)
			} else {
				// add to oldKeys
				oldKeys[k] = 0
			}
		}

		// all old unshared keys become DELETE operations
		for k := range oldKeys {
			q.Enqueue(&Operation{
				Op:   DELETE,
				Path: prefix + "[" + string(k) + "]",
				Old:  oldObj[k],
				New:  nil})
		}

		// all new unshared keys become CREATE operations
		for k := range newKeys {
			q.Enqueue(&Operation{
				Op:   CREATE,
				Path: prefix + "[" + string(k) + "]",
				Old:  nil,
				New:  newObj[k]})
		}
	}

	return
}

func enqueuePrimitiveOp(ok bool, old v.Value, new v.Value, q *OpQueue, prefix string) {
	if !ok || !v.IsEqual(old, new) {
		q.Enqueue(&Operation{
			Op:   UPDATE,
			Path: prefix,
			Old:  old,
			New:  new})
		return
	}

	q.Enqueue(&Operation{
		Op:   READ,
		Path: prefix,
		Old:  old,
		New:  new})
	return
}
