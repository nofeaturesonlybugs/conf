package parser

import (
	"io"
	"strings"

	"github.com/nofeaturesonlybugs/errors"
)

// Parser parses a string into a Configuration.
type Parser struct {
	options *Options
}

// NewParser creates a new Parser type.
func NewParser(options ...interface{}) *Parser {
	rv := &Parser{}
	rv.options = NewOptions(options...)
	return rv
}

// Parse parses a string.
func (me *Parser) Parse(s string) (Parsed, error) {
	if me == nil {
		return nil, errors.NilReceiver()
	}
	tokenizer := NewTokenizer(s)
	return Parse(tokenizer, me.options)
}

// ParseReader parses the reader.
func (me *Parser) ParseReader(r io.Reader) (Parsed, error) {
	s := &strings.Builder{}
	if _, err := io.Copy(s, r); err != nil {
		return nil, errors.Go(err)
	} else {
		return me.Parse(s.String())
	}
}
