package parser

import (
	"io"
	"strings"

	"github.com/nofeaturesonlybugs/errors"
)

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
	return Parse(NewTokenizer(s), me.Runes)
}

// ParseReader parses the reader.
func (me Parser) ParseReader(r io.Reader) (Parsed, error) {
	s := &strings.Builder{}
	if _, err := io.Copy(s, r); err != nil {
		return nil, errors.Go(err)
	} else {
		return me.Parse(s.String())
	}
}
