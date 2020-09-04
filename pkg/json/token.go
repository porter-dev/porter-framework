/**
 * Contains token definitions for Porter's JSON parser.
 */

package json

// Token is the set of lexical tokens of the JSON text syntax.
type Token int

// RichToken denotes the position, literal value, and token constant of a
// token.
type RichToken struct {
	pos Position // the position of the token
	tok Token    // the token
	lit string   // the literal value corresponding to a token
}

// create token enumeration
const (
	SPECIALBEG Token = iota
	EOF
	ILLEGAL
	SPECIALEND

	LITERALBEG
	NUMBER // 123
	STRING // "foo"
	TRUE   // true
	FALSE  // false
	NULL   // null
	LITERALEND

	OPERATORBEG
	LBRACE  // {
	LBRACK  // [
	RBRACE  // }
	RBRACK  // ]
	COMMA   // ,
	COLON   // :
	LINJECT // {{
	RINJECT // }}
	OPERATOREND
)

// Tokens is an array of the string value of each token, indexed by the
// token enumeration.
var Tokens = [...]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	NUMBER:  "NUMBER",
	STRING:  "STRING",
	TRUE:    "TRUE",
	FALSE:   "FALSE",
	NULL:    "NULL",
	LBRACE:  "{",
	LBRACK:  "[",
	RBRACE:  "}",
	RBRACK:  "]",
	COMMA:   ",",
	COLON:   ":",
	LINJECT: "{{",
	RINJECT: "}}",
}

// IsLiteral determines if a given token tok is a literal; true if yes,
// false otherwise
//
func (tok Token) IsLiteral() bool {
	return LITERALBEG < tok && tok < LITERALEND
}

// IsOperator determines if a given token tok is an operator; true if yes,
// false otherwise
//
func (tok Token) IsOperator() bool {
	return OPERATORBEG < tok && tok < OPERATOREND
}
