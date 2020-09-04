package value

type (
	// Value can be any Porter configuration value:
	// Object, Array, Integer, Float, String, True, False, Null
	Value interface{}

	// Object is a JSON object, consisting of String/Value pairs.
	Object map[String]Value

	// Array is a comma-separated list of Values
	Array []Value

	// Integer is a Go integer, subset of JSON numbers
	Integer int

	// Float is a Go float, subset of JSON numbers
	Float float64

	// String is a Go string -- equivalent to Go strings
	String string

	// Boolean is a Go bool, or JSON literal name tokens true or false
	Boolean bool
)
