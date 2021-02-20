package parser

import (
	"unicode"
)

// Token describes a token type.
type Token int

const (
	TokenNone       Token = 1 << iota
	TokenAlphaNum   Token = 1 << iota
	TokenNewline    Token = 1 << iota
	TokenPunct      Token = 1 << iota
	TokenWhiteSpace Token = 1 << iota
)

// Tokenizer reads a string into tokens.
type Tokenizer struct {
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
func NewTokenizer(s string) *Tokenizer {
	rv := &Tokenizer{s, len(s), 0, 0, TokenNone, 0, 0, TokenNone}
	return rv
}

// Eof returns true when the tokenizer has no more tokens to return.
func (me *Tokenizer) Eof() bool {
	if me == nil {
		return true
	}
	return me.n >= me.max
}

// Memory records the current Tokenizer position and a call to Rewid() will reset the Tokenizer to
// this position.  Use this to make consecutive calls to Peek() or Next() for look-ahead past a single token.
func (me *Tokenizer) Memory() {
	if me != nil {
		me.rewindN, me.rewindPeek, me.rewindTyp = me.n, me.peek, me.typ
	}
}

// Peek returns the next token and its type without advancing the internal counters; an empty string signals
// the end of the Tokenizer's input.
func (me *Tokenizer) Peek() (string, Token) {
	if me == nil || me.n >= me.max {
		return "", TokenNone
	} else if me.peek != 0 {
		return me.s[me.n : me.n+me.peek], me.typ
	}
	//
	next := func() (rune, bool) {
		me.peek++
		if me.n+me.peek >= me.max {
			return ' ', false
		} else {
			return rune(me.s[me.n+me.peek]), true
		}
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
func (me *Tokenizer) Next() (string, Token) {
	if me != nil {
		defer func() {
			me.n = me.n + me.peek
			me.peek = 0
		}()
	}
	return me.Peek()
}

// Rewind resets the Tokenizer to the postion recorded by calling Memory() or to the beginning
// of the string if Memory was never called.
func (me *Tokenizer) Rewind() {
	if me != nil {
		me.n, me.peek, me.typ = me.rewindN, me.rewindPeek, me.rewindTyp
	}
}
