package json

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Scanner defines a lexical scanner to extract tokens from a buffer.
//
type Scanner struct {
	buf *bytes.Buffer // Source buffer for advancing and scanning
	src []byte        // Source buffer for immutable access

	ch rune // current character

	err func(pos Position, msg string)

	srcPos      Position
	prevPos     Position
	lastCharLen int
}

const eof rune = rune(0)

// NewScanner creates and initializes a new instance of Scanner using src as
// its source content.
//
func NewScanner(src []byte) *Scanner {
	b := bytes.NewBuffer(src)

	s := &Scanner{
		buf: b,
		src: src,
	}

	// line number starts at 1
	s.srcPos.Line = 1

	s.ch = ' '

	return s
}

// next reads the next rune, returns eof if error or io.EOF is reached
//
func (s *Scanner) next() (rune, error) {
	r, size, err := s.buf.ReadRune()

	if err != nil {
		s.srcPos.Column++
		s.srcPos.Offset += s.lastCharLen
		s.lastCharLen = size
		s.ch = eof
		return eof, nil
	}

	s.prevPos = s.srcPos

	s.srcPos.Column++
	s.srcPos.Offset += s.lastCharLen
	s.lastCharLen = size

	s.ch = r

	switch {
	case r == utf8.RuneError && size == 1: // utf-8 replacement character
		return eof, fmt.Errorf("Column %d, Line %d: Improper UTF-8 encoding", s.srcPos.Column,
			s.srcPos.Line)
	case r == '\n': // new line
		s.srcPos.Line++
		s.srcPos.Column = 0
	case r == '\x00': // null character
		return eof, fmt.Errorf("Column %d, Line %d: Unexpected null character (0x00)", s.srcPos.Column,
			s.srcPos.Line)
	}

	return r, nil
}

// unread unreads the previous read Rune and updates the source position
func (s *Scanner) unread() {
	if err := s.buf.UnreadRune(); err != nil {
		panic(err) // this is user fault, we should catch it
	}
	s.srcPos = s.prevPos // put back last position
}

// peek reads the next rune without advancing the scanner
func (s *Scanner) peek() rune {
	peek, _, err := s.buf.ReadRune()
	if err != nil {
		return eof
	}

	s.buf.UnreadRune()
	return peek
}

// Scan gets the next token, skipping white space, and returns the token position,
// the token, and its literal string if applicable.
//
// If the token is illegal, Scan will log a syntax error and increment the error
// count, but will not panic.
//
// If the returned token is a literal, the literal string has the corresponding
// value. Otherwise, the literal string is the raw text value of the token.
func (s *Scanner) Scan() (RichToken, error) {
	s.skipWhitespace()
	var lit string
	var tok Token

	if s.ch == eof {
		return RichToken{
			pos: s.srcPos,
			tok: EOF,
			lit: "EOF",
		}, nil
	}

	switch ch := s.ch; {

	// if letter, determine if it is a literal name token (true, false, null)
	case isLetter(ch):
		richTok, err := s.scanIdentifier()

		if err != nil {
			return RichToken{}, err
		}

		return richTok, nil
	// if decimal, ensure that this is a number
	case isDecimal(ch) || ch == '-' && isDecimal(s.peek()):
		richTok, err := s.scanNumber()

		if err != nil {
			return RichToken{}, err
		}

		return richTok, nil
	case ch == '{' && s.peek() == '{':
		// consume next {
		s.next()
		lit = "{{"
		tok = LINJECT
	case ch == '}' && s.peek() == '}':
		// consume next }
		s.next()
		lit = "}}"
		tok = RINJECT
	default:
		lit = string(ch)

		switch ch {
		case eof:
			tok = EOF
		case '"':
			richTok, err := s.scanString()

			if err != nil {
				return RichToken{}, err
			}

			return richTok, nil
		case ',':
			tok = COMMA
		case ':':
			tok = COLON
		case '[':
			tok = LBRACK
		case ']':
			tok = RBRACK
		case '{':
			tok = LBRACE
		case '}':
			tok = RBRACE
		default:
			tok = ILLEGAL
			return RichToken{}, fmt.Errorf("Column %d, Line %d: Illegal token",
				s.srcPos.Column, s.srcPos.Line)
		}
	}

	pos := s.srcPos

	_, err := s.next()

	if err != nil {
		return RichToken{}, err
	}

	return RichToken{
		pos: pos,
		tok: tok,
		lit: lit,
	}, nil
}

// scanIdentifier scans an identifier and returns the literal name token
// and the accompanying literal string
func (s *Scanner) scanIdentifier() (RichToken, error) {
	offs := s.srcPos.Offset
	startCol := s.srcPos.Column
	ch := s.ch
	var err error

	for isLetter(ch) {
		ch, err = s.next()

		if err != nil {
			return RichToken{}, err
		}
	}

	lit := string(s.src[offs:s.srcPos.Offset])
	var tok Token

	switch lit {
	case "true":
		tok = TRUE
	case "false":
		tok = FALSE
	case "null":
		tok = NULL
	default:
		return RichToken{}, fmt.Errorf("Column %d, Line %d: Not a valid name token: must be true, false, or null. Strings must be enclosed in quotes",
			s.srcPos.Column, s.srcPos.Line)
	}

	pos := Position{
		Offset: offs,
		Column: startCol,
		Line:   s.srcPos.Line,
	}

	return RichToken{
		pos: pos,
		tok: tok,
		lit: lit,
	}, nil
}

// scanString scans a string and returns the literal value
// without "
func (s *Scanner) scanString() (RichToken, error) {
	// opening demarcation " already consumed
	ch, err := s.next()

	if err != nil {
		return RichToken{}, nil
	}

	startPos := s.srcPos

	for {
		if ch == eof {
			return RichToken{}, fmt.Errorf("Column %d, Line %d: String not terminated",
				s.srcPos.Column, s.srcPos.Line)
		}

		if ch == '"' {
			break
		}

		ch, err = s.next()

		if err != nil {
			return RichToken{}, nil
		}
	}

	lit := string(s.src[startPos.Offset:s.srcPos.Offset])

	// consume final "
	s.next()

	return RichToken{
		pos: startPos,
		tok: STRING,
		lit: lit,
	}, nil
}

// scanNumber scans a number and returns the token NUMBER
// and the accompanying value as a literal string
func (s *Scanner) scanNumber() (RichToken, error) {
	offs := s.srcPos.Offset
	startCol := s.srcPos.Column
	tok := NUMBER
	ch := s.ch
	var err error

	// if ch is -, move on to next digit
	if ch == '-' {
		ch, err = s.next()
	}

	if ch == '0' {
		if isDecimal(s.peek()) {
			return RichToken{}, fmt.Errorf("Column %d, Line %d: 0 cannot be followed by digit without . or exponent",
				s.srcPos.Column, s.srcPos.Line)
		}

		ch, err = s.next()
	}

	// can only have one 0, or arbitrary number of decimals
	if ch != '0' {
		for isDecimal(ch) {
			ch, err = s.next()
		}
	}

	if err != nil {
		return RichToken{}, err
	}

	// can only have one ., following by arbitrary number of decimals
	if ch == '.' {
		ch, err = s.next()

		if err != nil {
			return RichToken{}, err
		} else if !isDecimal(ch) {
			// must have at least one decimal
			return RichToken{}, fmt.Errorf("Column %d, Line %d: . must be followed by at least one digit",
				s.srcPos.Column, s.srcPos.Line)
		}

		for isDecimal(ch) {
			ch, err = s.next()

			if err != nil {
				return RichToken{}, err
			}
		}
	}

	// can have e or E
	if ch == 'e' || ch == 'E' {
		ch, err = s.next()

		if err != nil {
			return RichToken{}, err
		}

		// exponent has optional sign
		if ch == '+' || ch == '-' {
			ch, err = s.next()

			if err != nil {
				return RichToken{}, err
			}
		}

		// must have at least one decimal
		if !isDecimal(ch) {
			return RichToken{}, fmt.Errorf("Column %d, Line %d: e or E must be followed by at least one digit",
				s.srcPos.Column, s.srcPos.Line)
		}

		// exponent followed by arbitrary number of decimals
		for isDecimal(ch) {
			ch, err = s.next()

			if err != nil {
				return RichToken{}, err
			}
		}
	}

	return RichToken{
		pos: Position{
			Offset: offs,
			Column: startCol,
			Line:   s.srcPos.Line,
		},
		tok: tok,
		lit: string(s.src[offs:s.srcPos.Offset]),
	}, nil
}

// isLetter returns true if the given rune is a letter
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
}

// isDecimal returns true if the given rune is a decimal number
func isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// skipWhitespace advances the scanner to the next non-space and non-tab statement
func (s *Scanner) skipWhitespace() {
	for s.ch == '\u0009' || s.ch == '\u000A' || s.ch == '\u000D' || s.ch == '\u0020' {
		s.next()
	}
}
