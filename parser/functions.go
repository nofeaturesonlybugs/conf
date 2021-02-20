package parser

import (
	"github.com/nofeaturesonlybugs/errors"
)

// IsAssign is a function that returns true if the given rune is used to assign a value to a key.
type IsAssign func(rune) bool

// IsQuote is a function that returns true if the given rune is used to quote strings.
type IsQuote func(rune) bool

// Value is the value in a key=value configuration section; values can be singular or slices.
type Value struct {
	Last  string
	Slice []string
}

// Section is the key=value store of a configuration section.
type Section map[string]*Value

// Map returns section as a map[string][]string.
func (me Section) Map() map[string][]string {
	rv := make(map[string][]string)
	for k, v := range me {
		rv[k] = append([]string{}, v.Slice...)
	}
	return rv
}

// SectionBlock contains a slice of all sections that had the same name as well as the last section
// that had the name.
type SectionBlock struct {
	Last  Section
	Slice []Section
}

// Map returns the section block as a []map[string][]string.
func (me *SectionBlock) Map() []map[string][]string {
	rv := []map[string][]string{}
	for _, v := range me.Slice {
		rv = append(rv, v.Map())
	}
	return rv
}

// Parsed is a named map of name=SectionBlock.
type Parsed map[string]*SectionBlock

// Map returns parsed as a map[string][]map[string][]string
func (me Parsed) Map() map[string][]map[string][]string {
	rv := make(map[string][]map[string][]string)
	for k, v := range me {
		rv[k] = v.Map()
	}
	return rv
}

// Parse parses a configuration string.
func Parse(t *Tokenizer, options *Options) (Parsed, error) {
	if t == nil {
		return nil, errors.NilArgument("t").Type(t)
	} else if options == nil {
		return nil, errors.NilArgument("options").Type(options)
	}
	rv := make(Parsed)
	rv[""] = &SectionBlock{Last: make(Section)} // A default unnamed section.
	rv[""].Slice = []Section{rv[""].Last}
	current := rv[""].Last // Current block to put key=values into.
	for !t.Eof() {
		ScanWhileType(t, TokenNewline|TokenWhiteSpace) // Scan whitespace and newlines.
		if tok, typ := t.Peek(); typ == TokenPunct && rune(tok[0]) != options.sectionOpen {
			// This is a comment.
			if _, err := ScanComment(t); err != nil {
				return nil, errors.Go(err)
			}
		} else if typ == TokenPunct && rune(tok[0]) == options.sectionOpen {
			// This should be a section.
			t.Next()                          // Leading punctuation in section name.
			ScanWhileType(t, TokenWhiteSpace) // Ignore whitespace.
			if tok, typ = t.Next(); typ != TokenAlphaNum {
				return nil, errors.Errorf("Parse expects section name; got [%v], %v", typ, tok)
			} else if _, ok := rv[tok]; !ok {
				rv[tok] = &SectionBlock{} // First time a section with this name is found; create a block.
			}
			//
			rv[tok].Last = make(Section)                        // Create new section.
			rv[tok].Slice = append(rv[tok].Slice, rv[tok].Last) // Append new section to slice.
			current = rv[tok].Last
			//
			ScanWhileType(t, TokenWhiteSpace) // Ignore whitespace.
			if tok, typ = t.Next(); typ != TokenPunct && rune(tok[0]) != options.sectionClose {
				return nil, errors.Errorf("Parse expects section closing punctuation; got [%v], %v", typ, tok)
			}
		} else if typ == TokenAlphaNum {
			// Should be a key=value pair.
			key, value, err := tok, "", error(nil)
			//
			if key, err = ScanKey(t, options.isAssign); err != nil {
				return nil, errors.Go(err)
			}
			//
			// ScanKey scans leading and trailing whitespace around the key so next token should be
			// an assignment rune.
			if tok, typ = t.Next(); !options.isAssign(rune(tok[0])) {
				return nil, errors.Errorf("Parse expects assignment rune; got [%v] %v", typ, tok)
			}
			//
			ScanWhileType(t, TokenWhiteSpace) // Ignore whitespace.
			//
			// Only scan a value if next token is not a newline; if the next value is a newline
			// do nothing and append empty value.
			if _, typ2 := t.Peek(); typ2 != TokenNewline {
				if value, err = ScanValue(t, options.isQuote); err != nil {
					return nil, errors.Go(err)
				}
			}
			if _, ok := current[key]; !ok {
				current[key] = &Value{}
			}
			current[key].Last = value                              // Record last value in this section.
			current[key].Slice = append(current[key].Slice, value) // Also add it to the slice.
		}
	}
	return rv, nil
}

// ScanComment reads a comment from the Tokenizer and returns it.
func ScanComment(t *Tokenizer) (string, error) {
	if t == nil {
		return "", errors.NilArgument("t").Type(t)
	} else if tok, typ := t.Peek(); typ == TokenNone {
		return "", nil
	} else if typ != TokenPunct {
		return "", errors.Errorf("Scan comment expects punct; got [%v] %v", typ, tok)
	} else if rv, _ := t.Next(); false {
	} else if rest, err := ScanUntilType(t, TokenNewline); err != nil {
		return rv, errors.Go(err)
	} else {
		return rv + rest, nil
	}
	return "", nil
}

// ScanUntilType reads from the Tokenizer until the specified Token type would be the next token.
func ScanUntilType(t *Tokenizer, until Token) (string, error) {
	if t == nil {
		return "", errors.NilArgument("t").Type(t)
	}
	rv := ""
	for {
		if tok, typ := t.Peek(); typ == TokenNone {
			goto done
		} else if typ == until {
			goto done
		} else {
			tok, _ = t.Next()
			rv = rv + tok
		}

	}
done:
	return rv, nil
}

// ScanWhileType reads from the Tokenizer while the tokens match the given type(s).
func ScanWhileType(t *Tokenizer, types Token) (string, error) {
	if t == nil {
		return "", errors.NilArgument("t").Type(t)
	}
	rv := ""
	for {
		if tok, typ := t.Peek(); typ == TokenNone {
			goto done
		} else if typ&types == 0 {
			goto done
		} else {
			tok, _ = t.Next()
			rv = rv + tok
		}
	}
done:
	return rv, nil
}

// ScanKey reads from the Tokenizer a key, which is everything up until the first instance of an assignment rune
// with leading and trailing whitespace removed.
//
// Example keys:
//	hello =										-> hello
//	my message = 								-> 'my message'
//	  with whitespace      =					-> 'with whitespace'
//	punctuated.key =							-> 'punctuated.key'
func ScanKey(t *Tokenizer, isAssign IsAssign) (string, error) {
	if t == nil {
		return "", errors.NilArgument("t").Type(t)
	} else if isAssign == nil {
		return "", errors.NilArgument("isAssign").Type(isAssign)
	} else if tok, typ := t.Peek(); typ == TokenNone {
		return "", nil
	} else if typ != TokenAlphaNum {
		return "", errors.Errorf("Scan key expects first token to be alphanumeric; got [%v] %v", tok, typ)
	}
	//
	key := ""
	t.Memory()
	for {
		if tok2, typ2 := t.Peek(); typ2 == TokenPunct && isAssign(rune(tok2[0])) {
			// Assignment rune.
			if key == "" {
				t.Rewind()
				return "", errors.Errorf("Scan key hit assignment rune before a key.")
			} else {
				goto done
			}
		}
		tok, typ := t.Next()
		if typ == TokenAlphaNum {
			key = key + tok
		} else if typ == TokenPunct {
			if isAssign(rune(tok[0])) {
				goto done
			} else {
				key = key + tok
			}
		} else if typ == TokenWhiteSpace {
			// If next token is an assignment rune then this whitespace is trailing whitespace; otherwise
			// it is part of the key.
			if tok2, typ2 := t.Peek(); typ2 != TokenPunct || !isAssign(rune(tok2[0])) {
				key = key + tok
			} else {
				// We must `goto done` here or the next iteration will consume the token via t.Next()
				goto done
			}
		} else if typ == TokenNone || typ == TokenNewline {
			t.Rewind()
			return "", errors.Errorf("Scan key expects assignment rune but found end-of-line or end-of-input @ %v", key)
		}
	}
done:
	return key, nil
}

// ScanValue reads from the Tokenizer a value, either quoted or unquoted.
func ScanValue(t *Tokenizer, isQuote IsQuote) (string, error) {
	if t == nil {
		return "", errors.NilArgument("t").Type(t)
	} else if isQuote == nil {
		return "", errors.NilArgument("isQuote").Type(isQuote)
	} else if tok, typ := t.Peek(); typ == TokenNone {
		return "", nil
	} else if typ == TokenWhiteSpace || typ == TokenNewline {
		return "", errors.Errorf("Scan expected value but got whitespace or newline.")
	} else if typ == TokenPunct && isQuote(rune(tok[0])) {
		// A quoted string.
		quote, value := tok, ""
		t.Next() // Consume the starting quote value.
		for tok, typ = t.Next(); typ != TokenNone && tok != quote; tok, typ = t.Next() {
			value = value + tok
		}
		if typ == TokenNone {
			return "", errors.Errorf("Scan expected quoted value but value is missing end-quote.")
		}
		return value, nil
	} else {
		// Scan until end of line or end of tokenizer; however if last
		// token is whitespace then discard as trailing whitespace is only
		// kept with quoted values.
		rv := ""
		for {
			tok, typ := t.Next()
			if _, peek := t.Peek(); peek == TokenNone || peek == TokenNewline {
				if typ != TokenWhiteSpace {
					rv = rv + tok
				}
				goto done
			}
			rv = rv + tok
		}
	done:
		return rv, nil
	}
}
