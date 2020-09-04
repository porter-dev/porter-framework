package translator

import (
	"bytes"
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
func (s *Scanner) next() rune {
	r, size, err := s.buf.ReadRune()

	if err != nil {
		s.srcPos.Column++
		s.srcPos.Offset += s.lastCharLen
		s.lastCharLen = size
		s.ch = eof
		return eof
	}

	s.prevPos = s.srcPos

	s.srcPos.Column++
	s.srcPos.Offset += s.lastCharLen
	s.lastCharLen = size

	s.ch = r

	switch {
	case r == utf8.RuneError && size == 1: // utf-8 replacement character
		s.err(s.srcPos, "Improper UTF-8 encoding")
	case r == '\n': // new line
		s.srcPos.Line++
		s.srcPos.Column = 0
	case r == '\x00': // null character
		s.err(s.srcPos, "unexpected null character (0x00)")
		return eof
	}

	return r
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

// Scan gets the next set of tokens that match Porter's HEREDOC syntax,
// skipping white space, and returns the token position, the token, and
// its literal string if applicable.
//
// At the moment, this skips all tokens that are not HEREDOC tokens.
//
// If the token is illegal, Scan will log a syntax error and increment the error
// count, but will not panic.
//
// If the returned token is a literal, the literal string has the corresponding
// value. Otherwise, the literal string is the raw text value of the token.
func (s *Scanner) Scan() []RichToken {
	s.skipWhitespace()

	var tokens []RichToken

	if s.ch == '<' {
		if s.peek() == '<' {
			tokens = s.scanHeredoc()
		}
	}

	s.next()

	return tokens
}

// scanHeredoc scans a HEREDOC statement. This scan looks for code blocks,
// denoted by {{ and }}.
func (s *Scanner) scanHeredoc() []RichToken {
	startOffs := s.srcPos.Offset
	startCol := s.srcPos.Column
	var ch rune
	var name string

	// consume second <
	s.next()

	// parse first line until whitespace, store heredoc name
	nameOffset := s.srcPos.Offset
	for {
		ch := s.next()

		if ch == ' ' || ch == '\t' || ch == '\n' {
			name = string(s.src[nameOffset+s.lastCharLen : s.srcPos.Offset])
			break
		}

		if ch == eof {
			s.err(s.srcPos, "heredoc not terminated")
			return []RichToken{}
		}
	}

	var tokens = []RichToken{
		// add the opening LHEREDOC to the final tokens
		{
			pos: Position{
				Offset: startOffs,
				Column: startCol,
				Line:   s.srcPos.Line,
			},
			tok: LHEREDOC,
			lit: string(s.src[startOffs:s.srcPos.Offset]),
		},
	}

	ch = s.next()
	termNameCol := s.srcPos.Column
	termNameOffs := s.srcPos.Offset

	// when this loop starts we are on a new line/character
	for {
		if ch == eof {
			s.err(s.srcPos, "heredoc not terminated")
		}

		if ch == ' ' || ch == '\t' || ch == '\n' {
			termNameOffs = s.srcPos.Offset + 1
			termNameCol = s.srcPos.Column + 1
		}

		// look for HEREDOC terminus
		if ch == '>' && s.peek() == '>' {
			testName := string(s.src[termNameOffs:s.srcPos.Offset])

			if testName == name {
				s.next()

				tokens = append(tokens, RichToken{
					pos: Position{
						Offset: termNameOffs,
						Column: termNameCol,
						Line:   s.srcPos.Line,
					},
					tok: RHEREDOC,
					lit: string(s.src[termNameOffs : s.srcPos.Offset+1]),
				})

				break
			}
		}

		// we have hit a code block, get code literal
		if ch == '{' && s.peek() == '{' {
			tokens = append(tokens, RichToken{
				pos: s.srcPos,
				tok: LINJECT,
				lit: string(s.src[s.srcPos.Offset : s.srcPos.Offset+2]),
			})

			s.next() // consume second {

			var codeToken []RichToken = s.scanCode()

			tokens = append(tokens, codeToken...)
		}

		ch = s.next()
	}

	return tokens
}

// Scans past injected code, looking for end token }}
// Returns tokens [CODE, RINJECT]
func (s *Scanner) scanCode() []RichToken {
	ch := s.next()

	startOffs := s.srcPos.Offset
	startCol := s.srcPos.Column

	var tokens []RichToken

	for {
		if ch == '}' && s.peek() == '}' {
			codeToken := RichToken{
				pos: Position{
					Offset: startOffs,
					Column: startCol,
					Line:   s.srcPos.Line,
				},
				tok: CODE,
				lit: string(s.src[startOffs:s.srcPos.Offset]),
			}

			rInjectToken := RichToken{
				pos: s.srcPos,
				tok: RINJECT,
				lit: string(s.src[s.srcPos.Offset : s.srcPos.Offset+2]),
			}

			s.next()

			tokens = []RichToken{
				codeToken,
				rInjectToken,
			}

			break
		}

		ch = s.next()
	}

	return tokens
}

// skipWhitespace advances the scanner to the next non-space and non-tab statement
func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\r' {
		s.next()
	}
}
