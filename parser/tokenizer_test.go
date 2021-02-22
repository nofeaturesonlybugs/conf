package parser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/conf/parser"
)

func TestTokenizerNext(t *testing.T) {
	chk := assert.New(t)

	str := "asdf 1234 ;\t\t\r\nfinale;"
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)

	tok, typ := "", parser.TokenNone
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "asdf")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "1234")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, ";")
	chk.Equal(typ, parser.TokenPunct)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "\t\t")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "\r\n")
	chk.Equal(typ, parser.TokenNewline)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "finale")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, ";")
	chk.Equal(typ, parser.TokenPunct)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "")
	chk.Equal(typ, parser.TokenNone)
	//
	chk.Equal(true, tokenizer.Eof())
}

func TestTokenizerNextHeredoc(t *testing.T) {
	chk := assert.New(t)
	//
	str := `
[main]
# This is a comment.
value = true

[other]

`
	//
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	tok, typ := "", parser.TokenNone
	//
	_, typ = tokenizer.Next() // Ignore the specific newline values as they may change depending on dev environment.
	chk.Equal(typ, parser.TokenNewline)
	//
	// [main]
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "[")
	chk.Equal(typ, parser.TokenPunct)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "main")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "]")
	chk.Equal(typ, parser.TokenPunct)
	//
	_, typ = tokenizer.Next() // Ignoring the newline once again.
	chk.Equal(typ, parser.TokenNewline)
	//
	// # This is a comment.
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "#")
	chk.Equal(typ, parser.TokenPunct)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "This")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "is")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "a")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "comment")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, ".")
	chk.Equal(typ, parser.TokenPunct)
	//
	_, typ = tokenizer.Next() // Ignoring specific newline chars once again.
	chk.Equal(typ, parser.TokenNewline)
	//
	// value = true
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "value")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "=")
	chk.Equal(typ, parser.TokenPunct)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, " ")
	chk.Equal(typ, parser.TokenWhiteSpace)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "true")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	_, typ = tokenizer.Next() // Newlines...
	chk.Equal(typ, parser.TokenNewline)
	//
	// [other]
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "[")
	chk.Equal(typ, parser.TokenPunct)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "other")
	chk.Equal(typ, parser.TokenAlphaNum)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "]")
	chk.Equal(typ, parser.TokenPunct)
	//
	// Trailing newline(s)
	_, typ = tokenizer.Next()
	chk.Equal(typ, parser.TokenNewline)
	//
	// End
	tok, typ = tokenizer.Next()
	chk.Equal(tok, "")
	chk.Equal(typ, parser.TokenNone)
	//
	chk.Equal(true, tokenizer.Eof())
}

func TestTokenizerNextUnclosedQuote(t *testing.T) {
	chk := assert.New(t)
	//
	str := "'value"
	//
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	tok, typ := tokenizer.Next()
	chk.Equal("'", tok)
	chk.Equal(parser.TokenPunct, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("value", tok)
	chk.Equal(parser.TokenAlphaNum, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("", tok)
	chk.Equal(parser.TokenNone, typ)
	//
	chk.Equal(true, tokenizer.Eof())
}

func TestTokenizerMemoryRewind(t *testing.T) {
	chk := assert.New(t)
	//
	str := "hello = world\r\n[main"
	//
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	tok, typ := tokenizer.Next()
	chk.Equal("hello", tok)
	chk.Equal(parser.TokenAlphaNum, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(" ", tok)
	chk.Equal(parser.TokenWhiteSpace, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("=", tok)
	chk.Equal(parser.TokenPunct, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal(" ", tok)
	chk.Equal(parser.TokenWhiteSpace, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("world", tok)
	chk.Equal(parser.TokenAlphaNum, typ)
	//
	_, typ = tokenizer.Next()
	chk.Equal(parser.TokenNewline, typ)
	//
	tokenizer.Memory()
	//
	tok, typ = tokenizer.Next()
	chk.Equal("[", tok)
	chk.Equal(parser.TokenPunct, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("main", tok)
	chk.Equal(parser.TokenAlphaNum, typ)
	//
	chk.Equal(true, tokenizer.Eof())
	//
	tokenizer.Rewind()
	//
	tok, typ = tokenizer.Next()
	chk.Equal("[", tok)
	chk.Equal(parser.TokenPunct, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("main", tok)
	chk.Equal(parser.TokenAlphaNum, typ)
	//
	chk.Equal(true, tokenizer.Eof())
}

func TestTokenStrings(t *testing.T) {
	chk := assert.New(t)
	chk.Equal("None", parser.TokenNone.String())
	chk.Equal("AlphaNum", parser.TokenAlphaNum.String())
	chk.Equal("Newline", parser.TokenNewline.String())
	chk.Equal("Punctuation", parser.TokenPunct.String())
	chk.Equal("Whitespace", parser.TokenWhiteSpace.String())
	chk.Equal(true, strings.HasPrefix(parser.Token(-10).String(), "Unknown"))
}
