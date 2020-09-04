/**
 * Contains token definitions and keyword definitions for the
 * Porter configuration language.
 */

package translator

// Token is the set of lexical tokens of the Porter configuration language
type Token int

// RichToken denotes the position, literal value, and token constant of a
// token.
type RichToken struct {
	pos Position // the position of the token
	tok int      // the token
	lit string   // the literal value corresponding to a token
}

// create token enumeration
const (
	// literals
	LITERALBEG int = iota
	LHEREDOC       // <<yaml
	RHEREDOC       // yaml>>
	CODE           // code found within HEREDOCs
	LITERALEND

	OPERATORBEG
	LINJECT // {{
	RINJECT // }}
	OPERATOREND
)

// Tokens is an array of the string value of each token, indexed by the
// token enumeration.
var Tokens = [...]string{
	LHEREDOC: "LHEREDOC",
	RHEREDOC: "RHEREDOC",
	CODE:     "CODE",
	LINJECT:  "{{",
	RINJECT:  "}}",
}
