package json

import (
	"fmt"
	"strconv"

	v "github.com/porterdev/ego/internal/value"
)

// Inject takes in an array of bytes that contains Porter injection syntax {{}}
// and creates a Porter Value using those injected variables
func Inject(src string, v ...v.Value) (v.Value, error) {
	p := NewParser([]byte(src))
	p.injections = v
	p.currInj = 0

	return p.Parse()
}

// Parser type holds parser's internal state
type Parser struct {
	scanner Scanner

	// next token
	richTok RichToken

	// store injection variables
	injections v.Array
	currInj    int
}

// NewParser returns a new parser based on a set of input bytes
func NewParser(src []byte) (p Parser) {
	return Parser{
		scanner: *NewScanner(src),
	}
}

// Parse returns a Value representing the JSON object in Dynamic Syntax form
func (p *Parser) Parse() (v.Value, error) {
	err := p.next()

	if err != nil {
		return nil, err
	}

	var val v.Value

	val, err = p.parseValue()

	// advance to hopefully the end of parsing (EOF)
	p.next()

	if err == nil && p.richTok.tok != EOF {
		return nil, fmt.Errorf("Column %d, Line %d: Illegal token",
			p.scanner.srcPos.Column, p.scanner.srcPos.Line)
	}

	return val, err
}

// ----------------------------------------------------------------------------
// Parser helper methods
func (p *Parser) next() error {
	val, err := p.scanner.Scan()
	p.richTok = val

	return err
}

func (p *Parser) parseValue() (v.Value, error) {
	switch tok := p.richTok.tok; {
	case tok.IsLiteral():
		if tok == STRING {
			return v.String(p.richTok.lit), nil
		} else if tok == NUMBER {
			return p.parseNumber()
		} else if tok == TRUE {
			return v.Boolean(true), nil
		} else if tok == FALSE {
			return v.Boolean(false), nil
		} else if tok == NULL {
			return nil, nil
		}
	case tok.IsOperator():
		if tok == LBRACE {
			return p.parseObject()
		} else if tok == RBRACE {
			// this should have occurred in parseObject
			return nil, fmt.Errorf("Column %d, Line %d: Right brace } not preceded by left brace {",
				p.scanner.srcPos.Column, p.scanner.srcPos.Line)
		} else if tok == LBRACK {
			return p.parseArray()
		} else if tok == RBRACK {
			return nil, fmt.Errorf("Column %d, Line %d: Right bracket ] not preceded by left bracket [",
				p.scanner.srcPos.Column, p.scanner.srcPos.Line)
		} else if tok == COMMA {
			return nil, fmt.Errorf("Column %d, Line %d: Commas can only be present when separating array or object values",
				p.scanner.srcPos.Column, p.scanner.srcPos.Line)
		} else if tok == COLON {
			return nil, fmt.Errorf("Column %d, Line %d: Colons can only be used for string: value pairs inside objects",
				p.scanner.srcPos.Column, p.scanner.srcPos.Line)
		} else if tok == LINJECT {
			return p.parseInjection()
		} else if tok == RINJECT {
			return nil, fmt.Errorf("Column %d, Line %d: Right injection }} not preceded by left injection {{",
				p.scanner.srcPos.Column, p.scanner.srcPos.Line)
		}
	}

	return nil, nil
}

func (p *Parser) parseObject() (v.Value, error) {
	var obj v.Object = map[v.String]v.Value{}

	// LBRACE means we've entered an object
	if p.richTok.tok == LBRACE {
		// loop through strings until RBRACE
		prevTok := p.richTok.tok

		for {
			err := p.next()

			if err != nil {
				return nil, err
			}

			richTok := p.richTok

			// if an injection, parse and rewrite richTok
			if richTok.tok == LINJECT {
				inj, err := p.parseInjection()

				if err != nil {
					return nil, err
				}

				// verify value is a string
				str, ok := inj.(v.String)

				if !ok {
					return nil, fmt.Errorf("Column %d, Line %d: Key must be a string for a key injection",
						p.scanner.srcPos.Column, p.scanner.srcPos.Line)
				}

				richTok = RichToken{
					tok: STRING,
					pos: p.scanner.srcPos,
					lit: string(str),
				}
			}

			if richTok.tok == STRING || richTok.tok == LINJECT {
				err = p.next()

				if err != nil {
					return nil, err
				}

				if p.richTok.tok == COLON {
					err = p.next()

					if err != nil {
						return nil, err
					}

					val, _err := p.parseValue()

					if _err != nil {
						return nil, _err
					}

					obj[v.String(richTok.lit)] = val
				} else {
					return nil, fmt.Errorf("Column %d, Line %d: Must use colon : to define a string: value pair",
						p.scanner.srcPos.Column, p.scanner.srcPos.Line)
				}
			} else if (richTok.tok == RBRACE && prevTok == COMMA) ||
				(richTok.tok == COMMA && prevTok == COMMA) {
				return nil, fmt.Errorf("Column %d, Line %d: Commas must be followed by a string: value pair",
					p.scanner.srcPos.Column, p.scanner.srcPos.Line)
			} else if richTok.tok == RBRACE {
				break
			} else if richTok.tok != COMMA {
				return nil, fmt.Errorf("Column %d, Line %d: Left brace { not followed by comma-separated string: value pairs or right brace }",
					p.scanner.srcPos.Column, p.scanner.srcPos.Line)
			}

			prevTok = p.richTok.tok
		}
	}

	return obj, nil
}

func (p *Parser) parseNumber() (v.Value, error) {
	var err error

	// attempt conversion to int first
	_int, err := strconv.ParseInt(p.richTok.lit, 10, 64)

	if err != nil {
		_float, _err := strconv.ParseFloat(p.richTok.lit, 64)

		if _err != nil {
			return nil, fmt.Errorf("Column %d, Line %d: Unable to parse number into int or float",
				p.scanner.srcPos.Column, p.scanner.srcPos.Line)
		}

		// if number can be an integer, convert
		if _float == float64(int(_float)) {
			return v.Integer(_float), nil
		}

		return v.Float(_float), nil
	}

	return v.Integer(_int), nil
}

func (p *Parser) parseArray() (v.Value, error) {
	var arr v.Array = []v.Value{}

	// LBRACK means we've entered an array
	if p.richTok.tok == LBRACK {
		prevTok := LBRACK
		// loop through values until RBRACK
		for {
			err := p.next()

			if err != nil {
				return nil, err
			}

			richTok := p.richTok

			if richTok.tok == EOF {
				return nil, fmt.Errorf("Column %d, Line %d: No closing bracket in array",
					p.scanner.srcPos.Column, p.scanner.srcPos.Line)
			} else if richTok.tok == RBRACK && prevTok == COMMA {
				return nil, fmt.Errorf("Column %d, Line %d: Commas must be followed by a value",
					p.scanner.srcPos.Column, p.scanner.srcPos.Line)
			} else if richTok.tok == RBRACK {
				break
			} else if richTok.tok != COMMA && prevTok != COMMA && prevTok != LBRACK {
				return nil, fmt.Errorf("Column %d, Line %d: Values must be separated by a comma",
					p.scanner.srcPos.Column, p.scanner.srcPos.Line)
			} else if richTok.tok == COMMA && (prevTok == LBRACK || prevTok == COMMA) {
				return nil, fmt.Errorf("Column %d, Line %d: Must have a value between array elements",
					p.scanner.srcPos.Column, p.scanner.srcPos.Line)
			} else if richTok.tok != COMMA {
				obj, err := p.parseValue()

				if err != nil {
					return nil, err
				}

				arr = append(arr, obj)
			}

			prevTok = richTok.tok
		}
	}

	return arr, nil
}

func (p *Parser) parseInjection() (v.Value, error) {
	// iterate until hitting right injection
	for p.richTok.tok != RINJECT {
		if p.richTok.tok == EOF {
			return nil, fmt.Errorf("Column %d, Line %d: No closing injection }} found",
				p.scanner.srcPos.Column, p.scanner.srcPos.Line)
		}

		p.next()
	}

	res := p.injections[p.currInj]

	p.currInj++

	return res, nil
}
