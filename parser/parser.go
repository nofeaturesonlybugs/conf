package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/nofeaturesonlybugs/errors"
)

// State describes the parser state.
type State int

// Enums for parser state.
const (
	StateNone    State = 1 << iota
	StateComment State = 1 << iota
	StateSection State = 1 << iota
	StateKey     State = 1 << iota
	StateValue   State = 1 << iota
)

// String returns the State as a string.
func (s State) String() string {
	switch s {
	case StateNone:
		return "None"
	case StateComment:
		return "Comment"
	case StateSection:
		return "Section"
	case StateKey:
		return "Key"
	case StateValue:
		return "Value"
	}
	return fmt.Sprintf("Unknown %T= %v", s, int(s))
}

// Runes defines some important runes for a parser.
type Runes struct {
	// Any rune present in SectionOpen begins a configuration section line.
	SectionOpen []rune
	// Any rune present in SectionClose ends a configuration section line.
	SectionClose []rune
	// Any rune present in Assign acts as an assignment rune and delimits <key ASSIGN value>.
	Assign []rune
	// Any rune present in Quote acts as a quotation rune; quoted values must use the same rune
	// to start and end the quotation.
	Quote []rune
}

// Within tests if r is within possible.
func (me Runes) Within(r rune, possible []rune) bool {
	for _, v := range possible {
		if v == r {
			return true
		}
	}
	return false
}

// IsOpenSection returns true if the rune opens a section.
func (me Runes) IsOpenSection(r rune) bool {
	return me.Within(r, me.SectionOpen)
}

// IsCloseSection returns true if the rune closes a section.
func (me Runes) IsCloseSection(r rune) bool {
	return me.Within(r, me.SectionClose)
}

// IsAssign returns true if the rune is an assignment rune.
func (me Runes) IsAssign(r rune) bool {
	return me.Within(r, me.Assign)
}

// IsQuote returns true if the rune is a quotation rune.
func (me Runes) IsQuote(r rune) bool {
	return me.Within(r, me.Quote)
}

// Parser parses a string into a Configuration.
type Parser struct {
	Runes
}

// DefaultParser is a Parser with common settings.
var DefaultParser = Parser{
	Runes: Runes{
		Assign:       []rune{'='},
		Quote:        []rune{'\'', '"', '`'},
		SectionOpen:  []rune{'['},
		SectionClose: []rune{']'},
	},
}

// Parse parses a string.
func (me Parser) Parse(s string) (Parsed, error) {
	var err error
	//
	assign := func(s string, t Token) bool {
		return t == TokenPunct && me.IsAssign(rune(s[0]))
	}
	quote := func(s string, t Token) bool {
		return t == TokenPunct && me.IsQuote(rune(s[0]))
	}
	closeSection := func(s string, t Token) bool {
		return t == TokenPunct && me.IsCloseSection(rune(s[0]))
	}
	openSection := func(s string, t Token) bool {
		return t == TokenPunct && me.IsOpenSection(rune(s[0]))
	}
	//
	rv := make(Parsed)
	rv[""] = &SectionBlock{Last: make(Section)} // A default unnamed section.
	rv[""].Slice = []Section{rv[""].Last}
	current := rv[""].Last // Current block to put key=values into.
	//
	section, key, value, previous, quotation := "", "", "", "", ""
	//
	t, st := NewTokenizer(s), StateNone
	for err == nil && !t.EOF() {
		str, tok := t.Next()
		switch st {
		case StateNone:
			if tok == TokenAlphaNum {
				// Beginning of key = value
				st, key, value, previous, quotation = StateKey, str, "", "", ""
			} else if tok == TokenPunct {
				// Punctuation is either opening a section or beginning a comment.
				if openSection(str, tok) {
					st, section, previous = StateSection, "", ""
				} else {
					st = StateComment
				}
			}

		case StateComment:
			if tok == TokenNewline {
				st = StateNone
			}

		case StateSection:
			if tok == TokenAlphaNum {
				section, previous = section+previous+str, ""
			} else if tok == TokenWhiteSpace {
				if section != "" {
					previous = str
				}
				// Whitespace in section name has to be followed by section close or another alphanum.
				if peek, peekT := t.Peek(); peekT != TokenAlphaNum && !closeSection(peek, peekT) {
					err = errors.Errorf("Parsing section name; unexpected token= %v", peek)
				}
			} else if tok == TokenPunct {
				if closeSection(str, tok) {
					current = make(Section)
					if existing, ok := rv[section]; !ok {
						rv[section] = &SectionBlock{
							Last:  current,
							Slice: []Section{current},
						}
					} else {
						existing.Last = current
						existing.Slice = append(existing.Slice, current)
					}
					st = StateNone
				} else {
					previous = str
					// Punctuation in section name has to be followed by another alphanum.
					if peek, peekT := t.Peek(); peekT != TokenAlphaNum {
						err = errors.Errorf("Parsing section name; unexpected token= %v", peek)
					}
				}
			}

		case StateKey:
			if tok == TokenAlphaNum {
				key, previous = key+previous+str, ""
			} else if tok == TokenWhiteSpace {
				if key != "" {
					previous = str
				}
				// Whitespace in key has to be followed by assign or another alphanum.
				if peek, peekT := t.Peek(); peekT != TokenAlphaNum && !assign(peek, peekT) {
					err = errors.Errorf("Parsing key name; unexpected token= %v", peek)
				}
			} else if tok == TokenPunct {
				if assign(str, tok) {
					st, previous = StateValue, ""
				} else {
					previous = str
					// Punctuation in key has to be followed by another alphanum.
					if peek, peekT := t.Peek(); peekT != TokenAlphaNum {
						err = errors.Errorf("Parsing key name; unexpected token= %v", peek)
					}
				}
			}

		case StateValue:
			if tok == TokenPunct && quote(str, tok) {
				if value == "" && quotation == "" {
					quotation = str
				} else if str != quotation {
					value, previous = value+previous+str, ""
				} else {
					quotation = ""
					st = StateNone
				}
			} else if tok == TokenPunct || tok == TokenAlphaNum {
				value, previous = value+previous+str, ""
			} else if tok == TokenWhiteSpace {
				if quotation != "" {
					value, previous = value+previous+str, ""
				} else if value != "" {
					previous = str
				}
			} else if tok == TokenNewline {
				if quotation != "" {
					value, previous = value+previous+str, ""
				} else {
					st = StateNone
				}
			}
			if st == StateNone { // Intentionally not attached to previous if..else block
				// Value was completed.
				if _, ok := current[key]; !ok {
					current[key] = &Value{}
				}
				current[key].Last = value
				current[key].Slice = append(current[key].Slice, value)
			}
		}
	}
	if st != StateNone {
		return rv, errors.Errorf("Unexpected EOF while parsing %v", st.String())
	}
	return rv, err
}

// ParseReader parses the reader.
func (me Parser) ParseReader(r io.Reader) (Parsed, error) {
	s := &strings.Builder{}
	if _, err := io.Copy(s, r); err != nil {
		return nil, errors.Go(err)
	}
	return me.Parse(s.String())
}
