package json

// Position is the location of a token in a byte sequence.
type Position struct {
	// offset, starting at 0
	// needed to index the src bytes, so that we can retrieve slices
	// and such
	Offset int
	Line   int // line number, starting at 1
	Column int // column number, starting at 1
}
