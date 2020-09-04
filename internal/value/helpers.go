package value

import (
	"fmt"
	"strconv"
)

// IsEqual compares two values, and returns true if the values are
// equal. Configuration equality is defined as:
//
// - Key equality: each key for each map[String]Value contains the same literal
// - Type equality: each type is the exact same at the most restrictive type
//   definition
// - Literal equality: each literal contains the same value
//
// This function returns false on any value that is not considered a "Porter Configuration"
// type -- see types.go in porter package for explicit Porter types.
func IsEqual(v1, v2 Value) bool {
	if v1 == nil && v2 == nil {
		return true
	} else if v1 == nil || v2 == nil {
		return false
	}

	switch v1.(type) {
	case Boolean:
		_, ok := v2.(Boolean)
		return ok && v1 == v2
	case Float:
		_, ok := v2.(Float)
		return ok && v1 == v2
	case Integer:
		_, ok := v2.(Integer)
		return ok && v1 == v2
	case String:
		_, ok := v2.(String)
		return ok && v1 == v2
	case Array:
		v1Arr := v1.(Array)
		v2Arr, ok := v2.(Array)

		if !ok {
			return false
		}

		// check that arrays are the same length
		if len(v1Arr) != len(v2Arr) {
			return false
		}

		// iterate through each element
		for i := range v1Arr {
			if !IsEqual(v1Arr[i], v2Arr[i]) {
				return false
			}
		}

		return true
	case Object:
		v1Obj := v1.(Object)
		v2Obj, ok := v2.(Object)

		if !ok {
			return false
		}

		if len(v1Obj) != len(v2Obj) {
			return false
		}

		for k, v := range v1Obj {
			w, ok := v2Obj[k]

			if !ok || !IsEqual(v, w) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

// Get retrieves a Value at a certain path within a configuration object
func Get(v Value, path string) (Value, error) {
	if path == "" {
		return v, nil
	}

	offs := 0
	curr := rune(path[0])

	// if here, indexing array or object
	if curr == '[' {
		for curr != ']' {
			offs++
			curr = rune(path[offs])
		}

		// make sure we're at an object or an array
		switch v.(type) {
		case Array:
			vArr := v.(Array)

			// parse path to integer
			i, err := strconv.ParseInt(path[1:offs], 10, 64)

			if err != nil {
				return nil, err
			}

			if len(path) == offs+1 {
				return vArr[i], nil
			}

			return Get(vArr[i], path[offs+1:])
		case Object:
			vObj := v.(Object)

			if len(path) == offs+1 {
				return vObj[String(path[1:offs])], nil
			}

			return Get(vObj[String(path[1:offs])], path[offs+1:])
		default:
			return nil, fmt.Errorf("Not an object, cannot index using brackets []")
		}
	}

	// if leading ., indexing on object -- can just remove period
	if curr == '.' {
		if len(path) > 1 {
			return Get(v, path[1:])
		}

		return nil, fmt.Errorf("Path cannot end in period (.)")
	}

	// iterate until period
	for i := 0; i < len(path); i++ {
		curr = rune(path[i])

		if curr == '.' || curr == '[' {
			// make sure we're at an object
			vObj, ok := v.(Object)

			if !ok {
				return nil, fmt.Errorf("Not an object: cannot index (cannot use [] or .)")
			}

			if len(path) > i && curr == '.' {
				return Get(vObj[String(path[offs:i])], path[i+1:])
			} else if len(path) > i && curr == '[' {
				return Get(vObj[String(path[offs:i])], path[i:])
			}

			return nil, fmt.Errorf("Cannot end in period (.)")
		}
	}

	// if here, must be indexing an object
	vObj, ok := v.(Object)

	if !ok {
		return nil, fmt.Errorf("Not an object: cannot index on a field")
	}

	return vObj[String(path[0:len(path)])], nil
}
