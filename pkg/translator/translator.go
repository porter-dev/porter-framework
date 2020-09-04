package translator

// Translator implements a translator from .gop to .go files.
type Translator struct {
	src []byte // Source buffer
	res []byte // resulting translation

	scanner Scanner // scanner instance

	currBlock []RichToken // the current token block
}

// NewTranslator creates a new Translator based on a source byte array
func NewTranslator(src []byte) Translator {
	trans := Translator{
		src:     src,
		scanner: *NewScanner(src),
	}

	return trans
}

// TranslateToString takes in a HEREDOC that isn't a supported kind
// and translates it to dynamic string using the injected variables
func (t *Translator) TranslateToString() []byte {
	var prevPos int = 0

	// iterate through token blocks, translating HEREDOCs into Go
	for t.scanner.ch != eof {
		t.currBlock = t.scanner.Scan()

		// iterate through tokens, extract code blocks
		for _, richTok := range t.currBlock {
			switch richTok.tok {
			case LHEREDOC:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset]...)
				t.res = append(t.res, []byte{'`'}...)
			case LINJECT:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset]...)
				t.res = append(t.res, []byte{'`', ' ', '+', ' '}...)
			case CODE:
				t.res = append(t.res, []byte(richTok.lit)...)
			case RINJECT:
				t.res = append(t.res, []byte{' ', '+', ' ', '`'}...)
			case RHEREDOC:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset]...)
				t.res = append(t.res, []byte{'`'}...)
			}

			prevPos = richTok.pos.Offset + len(richTok.lit)
		}
	}

	t.res = append(t.res, t.src[prevPos:t.scanner.srcPos.Offset]...)

	return t.res
}

// TranslateToJSON takes in a JSON HEREDOC and translates it to a function that
// generates a Porter configuration using the json package.
func (t *Translator) TranslateToJSON() []byte {
	var prevPos int = 0
	var injections []string = make([]string, 4)

	// iterate through token blocks, translating HEREDOCs into Go
	for t.scanner.ch != eof {
		t.currBlock = t.scanner.Scan()

		// iterate through tokens, extract code blocks
		for _, richTok := range t.currBlock {
			switch richTok.tok {
			case LHEREDOC:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset]...)
				t.res = append(t.res, []byte("json.Inject(`")...)
			case LINJECT:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset+len(richTok.lit)]...)
			case CODE:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset+len(richTok.lit)]...)
				injections = append(injections, richTok.lit)
			case RINJECT:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset+len(richTok.lit)]...)
			case RHEREDOC:
				t.res = append(t.res, t.src[prevPos:richTok.pos.Offset]...)

				// add injected variable literals to end of argument
				terminus := "`"

				for _, c := range injections {
					if len(c) > 0 {
						terminus += ","
						terminus += c
					}
				}

				terminus += ")"

				t.res = append(t.res, []byte(terminus)...)

			}

			prevPos = richTok.pos.Offset + len(richTok.lit)
		}
	}

	t.res = append(t.res, t.src[prevPos:t.scanner.srcPos.Offset]...)

	return t.res
}
