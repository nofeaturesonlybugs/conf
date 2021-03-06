package parser

import (
	"fmt"
	"unicode"
)

// Token describes a token type.
type Token int

// Enums for token types returned from the Tokenizer.
const (
	TokenNone       Token = 1 << iota
	TokenAlphaNum   Token = 1 << iota
	TokenNewline    Token = 1 << iota
	TokenPunct      Token = 1 << iota
	TokenWhiteSpace Token = 1 << iota
)

// String returns the Token as a string.
func (t Token) String() string {
	switch t {
	case TokenNone:
		return "None"
	case TokenAlphaNum:
		return "AlphaNum"
	case TokenNewline:
		return "Newline"
	case TokenPunct:
		return "Punctuation"
	case TokenWhiteSpace:
		return "Whitespace"
	}
	return fmt.Sprintf("Unknown %T= %v", t, int(t))
}

// Tokenizer returns tokens.
type Tokenizer interface {
	// EOF returns true when the tokenizer has no more tokens to return.
	EOF() bool
	// Memory records the current Tokenizer position and a call to Rewid() will reset the Tokenizer to
	// this position.  Use this to make consecutive calls to Peek() or Next() for look-ahead past a single token.
	Memory()
	// Peek returns the next token and its type without advancing the internal counters; an empty string signals
	// the end of the Tokenizer's input.
	Peek() (string, Token)
	// Next returns the token and its type that would be returned by Peek() but then advances internal
	// counters to the next possible token.
	Next() (string, Token)
	// Rewind resets the Tokenizer to the position recorded by calling Memory() or to the beginning
	// of the string if Memory was never called.
	Rewind()
}

// tokenizer reads a string into tokens.
type tokenizer struct {
	// String to tokenize.
	s string
	// Length of s.
	max int
	// Current index into s.
	n int
	// The number of runes peeked at; if 0 then nothing yet peeked.
	peek int
	// The peeked token's type.
	typ Token
	// A stack of n, peek, and typ to allow for forward scanning but rewinding.
	rewindN    int
	rewindPeek int
	rewindTyp  Token
}

// NewTokenizer creates a new Tokenizer type.
func NewTokenizer(s string) Tokenizer {
	rv := &tokenizer{s, len(s), 0, 0, TokenNone, 0, 0, TokenNone}
	return rv
}

// EOF returns true when the tokenizer has no more tokens to return.
func (me *tokenizer) EOF() bool {
	return me.n >= me.max
}

// Memory records the current Tokenizer position and a call to Rewid() will reset the Tokenizer to
// this position.  Use this to make consecutive calls to Peek() or Next() for look-ahead past a single token.
func (me *tokenizer) Memory() {
	me.rewindN, me.rewindPeek, me.rewindTyp = me.n, me.peek, me.typ
}

// Peek returns the next token and its type without advancing the internal counters; an empty string signals
// the end of the Tokenizer's input.
func (me *tokenizer) Peek() (string, Token) {
	if me.n >= me.max {
		return "", TokenNone
	} else if me.peek != 0 {
		return me.s[me.n : me.n+me.peek], me.typ
	}
	//
	next := func() (rune, bool) {
		me.peek++
		if me.n+me.peek >= me.max {
			return ' ', false
		}
		return rune(me.s[me.n+me.peek]), true
	}
	//
	r, ok := rune(me.s[me.n]), true
	if r == ' ' || r == '\t' {
		// Consumes spaces and tabs.
		for ; ok && (r == ' ' || r == '\t'); r, ok = next() {
		}
		me.typ = TokenWhiteSpace
	} else if r == '\r' || r == '\n' {
		// Consumes newlines.
		for ; ok && (r == '\r' || r == '\n'); r, ok = next() {
		}
		me.typ = TokenNewline
	} else if unicode.IsDigit(r) || unicode.IsLetter(r) {
		for ; ok && (unicode.IsDigit(r) || unicode.IsLetter(r)); r, ok = next() {
		}
		me.typ = TokenAlphaNum
	} else {
		me.peek = 1
		me.typ = TokenPunct
	}
	//
	return me.s[me.n : me.n+me.peek], me.typ
}

// Next returns the token and its type that would be returned by Peek() but then advances internal
// counters to the next possible token.
func (me *tokenizer) Next() (string, Token) {
	defer func() {
		me.n = me.n + me.peek
		me.peek = 0
	}()
	return me.Peek()
}

// Rewind resets the Tokenizer to the position recorded by calling Memory() or to the beginning
// of the string if Memory was never called.
func (me *tokenizer) Rewind() {
	me.n, me.peek, me.typ = me.rewindN, me.rewindPeek, me.rewindTyp
}
