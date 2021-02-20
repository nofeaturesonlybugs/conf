package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/conf/parser"
)

func TestParse(t *testing.T) {
	chk := assert.New(t)
	//
	str := `
# Ignore this comment.
key = value
key1 = value1
`
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	options := parser.NewOptions()
	chk.NotNil(options)
	//
	parsed, err := parser.Parse(tokenizer, options)
	chk.NoError(err)
	chk.NotNil(parsed)
	//
	chk.NotNil(parsed[""])
	chk.Nil(parsed["hello"])
}

func TestParse_NoWhiteSpace(t *testing.T) {
	chk := assert.New(t)
	//
	s := `
[service]
name=examplesvc
label=Example Service
description=This is an example service.
`
	//
	parser := parser.NewParser()
	parsed, err := parser.Parse(s)
	chk.NoError(err)
	chk.NotNil(parsed["service"])
	chk.Equal("examplesvc", parsed["service"].Last["name"].Last)
	chk.Equal("Example Service", parsed["service"].Last["label"].Last)
	chk.Equal("This is an example service.", parsed["service"].Last["description"].Last)
}

func TestParse_EmptyValue(t *testing.T) {
	chk := assert.New(t)
	//
	s := `
; a comment
empty = 
; new comment
`
	//
	parser := parser.NewParser()
	parsed, err := parser.Parse(s)
	chk.NoError(err)
	chk.NotNil(parsed[""])
	chk.Equal("", parsed[""].Last["empty"].Last)
}

func TestParseSections(t *testing.T) {
	chk := assert.New(t)
	//
	str := `
# Ignore this comment.
key = value
key1 = value1

# Options section.
[options]
foo = bar
baz = faz

# Qwerty section.
[qwerty]
name = Fred
age = 42
`
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	options := parser.NewOptions()
	chk.NotNil(options)
	//
	parsed, err := parser.Parse(tokenizer, options)
	chk.NoError(err)
	chk.NotNil(parsed)
	//
	chk.NotNil(parsed[""])
	chk.NotNil(parsed["options"])
	chk.NotNil(parsed["qwerty"])
	chk.Nil(parsed["hello"])
	// Checking values.
	chk.Equal("bar", parsed["options"].Last["foo"].Last)
	chk.Equal("faz", parsed["options"].Last["baz"].Last)
	chk.Equal("Fred", parsed["qwerty"].Last["name"].Last)
	chk.Equal("42", parsed["qwerty"].Last["age"].Last)
	// Checking slices.
	chk.Equal(1, len(parsed["options"].Last["foo"].Slice))
	chk.Equal(1, len(parsed["options"].Last["baz"].Slice))
	chk.Equal(1, len(parsed["qwerty"].Last["name"].Slice))
	chk.Equal(1, len(parsed["qwerty"].Last["age"].Slice))
}

func TestParseSectionsSlices(t *testing.T) {
	chk := assert.New(t)
	//
	str := `
# Ignore this comment.
key = value
key1 = value1

# Options section.
[options]
foo = bar
baz = faz

# Repeated Options section.
[options]
foo = boo
baz = zaz

# Qwerty section.
[qwerty]
name = Fred
age = 42

# Repeated Qwerty section.
[qwerty]
name = Wilma
age = 38
`
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	options := parser.NewOptions()
	chk.NotNil(options)
	//
	parsed, err := parser.Parse(tokenizer, options)
	chk.NoError(err)
	chk.NotNil(parsed)
	//
	chk.NotNil(parsed[""])
	chk.NotNil(parsed["options"])
	chk.NotNil(parsed["qwerty"])
	chk.Nil(parsed["hello"])
	// Checking values.
	chk.Equal("boo", parsed["options"].Last["foo"].Last)
	chk.Equal("zaz", parsed["options"].Last["baz"].Last)
	chk.Equal("Wilma", parsed["qwerty"].Last["name"].Last)
	chk.Equal("38", parsed["qwerty"].Last["age"].Last)
	// slice=0
	chk.Equal("bar", parsed["options"].Slice[0]["foo"].Last)
	chk.Equal("faz", parsed["options"].Slice[0]["baz"].Last)
	chk.Equal("Fred", parsed["qwerty"].Slice[0]["name"].Last)
	chk.Equal("42", parsed["qwerty"].Slice[0]["age"].Last)
	// slice=1
	chk.Equal("boo", parsed["options"].Slice[1]["foo"].Last)
	chk.Equal("zaz", parsed["options"].Slice[1]["baz"].Last)
	chk.Equal("Wilma", parsed["qwerty"].Slice[1]["name"].Last)
	chk.Equal("38", parsed["qwerty"].Slice[1]["age"].Last)
	// Checking slices.
	chk.Equal(2, len(parsed["options"].Slice))
	chk.Equal(2, len(parsed["qwerty"].Slice))
	//
	chk.Equal(1, len(parsed["options"].Last["foo"].Slice))
	chk.Equal(1, len(parsed["options"].Last["baz"].Slice))
	chk.Equal(1, len(parsed["qwerty"].Last["name"].Slice))
	chk.Equal(1, len(parsed["qwerty"].Last["age"].Slice))
}

func TestParseSectionsSlicesValueSlices(t *testing.T) {
	chk := assert.New(t)
	//
	str := `
# Ignore this comment.
key = value
key1 = value1

# Options section.
[options]
foo = bar
foo = bar2
baz = faz
baz = faz2

# Repeated Options section.
[options]
foo = boo
foo = boo2
baz = zaz
baz = zaz2

`
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	options := parser.NewOptions()
	chk.NotNil(options)
	//
	parsed, err := parser.Parse(tokenizer, options)
	chk.NoError(err)
	chk.NotNil(parsed)
	//
	chk.NotNil(parsed[""])
	chk.NotNil(parsed["options"])
	chk.Nil(parsed["hello"])
	// Checking values.
	chk.Equal("boo2", parsed["options"].Last["foo"].Last)
	chk.Equal("zaz2", parsed["options"].Last["baz"].Last)
	// slice=0
	chk.Equal("bar2", parsed["options"].Slice[0]["foo"].Last)
	chk.Equal("faz2", parsed["options"].Slice[0]["baz"].Last)
	// slice=1
	chk.Equal("boo2", parsed["options"].Slice[1]["foo"].Last)
	chk.Equal("zaz2", parsed["options"].Slice[1]["baz"].Last)
	// Checking slices.
	chk.Equal(2, len(parsed["options"].Slice))
	//
	chk.Equal(2, len(parsed["options"].Last["foo"].Slice))
	chk.Equal(2, len(parsed["options"].Last["baz"].Slice))
}
func TestScanComment(t *testing.T) {
	chk := assert.New(t)
	//
	str := "# This is a comment."
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	scanned, err := parser.ScanComment(tokenizer)
	chk.NoError(err)
	chk.Equal("# This is a comment.", scanned)
	//
	//
	str = `# This is a comment.
# This is the second comment.
`
	tokenizer = parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	scanned, err = parser.ScanComment(tokenizer)
	chk.NoError(err)
	chk.Equal("# This is a comment.", scanned)
	//
	_, typ := tokenizer.Next()
	chk.Equal(parser.TokenNewline, typ)
	//
	scanned, err = parser.ScanComment(tokenizer)
	chk.NoError(err)
	chk.Equal("# This is the second comment.", scanned)

}

func TestScanUntilType(t *testing.T) {
	chk := assert.New(t)
	//
	str := " This is a comment."
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	scanned, err := parser.ScanUntilType(tokenizer, parser.TokenPunct)
	chk.NoError(err)
	chk.Equal(" This is a comment", scanned)
	//
	//
	str = " This is a comment.\r\nkey = value"
	tokenizer = parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	scanned, err = parser.ScanUntilType(tokenizer, parser.TokenNewline)
	chk.NoError(err)
	chk.Equal(" This is a comment.", scanned)
}

func TestScanKey(t *testing.T) {
	chk := assert.New(t)
	//
	isAssign := func(r rune) bool {
		return r == '='
	}
	//
	str := "key="
	tokenizer := parser.NewTokenizer(str)
	scanned, err := parser.ScanKey(tokenizer, isAssign)
	chk.Equal("key", scanned)
	chk.NoError(err)
	//
	str = "key = "
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanKey(tokenizer, isAssign)
	chk.Equal("key", scanned)
	chk.NoError(err)
	//
	str = "key       = "
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanKey(tokenizer, isAssign)
	chk.Equal("key", scanned)
	chk.NoError(err)
	//
	str = "my key       = "
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanKey(tokenizer, isAssign)
	chk.Equal("my key", scanned)
	chk.NoError(err)
	//
	str = "my-key       = "
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanKey(tokenizer, isAssign)
	chk.Equal("my-key", scanned)
	chk.NoError(err)
	//
	// End of line without assignment rune.
	str = "key       \r\n"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanKey(tokenizer, isAssign)
	chk.Equal("", scanned)
	chk.Error(err)
	//
	// End of input without assignment rune.
	str = "key       "
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanKey(tokenizer, isAssign)
	chk.Equal("", scanned)
	chk.Error(err)
	//
}

func TestScanValue(t *testing.T) {
	chk := assert.New(t)
	//
	isQuote := func(r rune) bool {
		return r == '\'' || r == '"' || r == '`'
	}
	//
	str := " value"
	tokenizer := parser.NewTokenizer(str)
	scanned, err := parser.ScanValue(tokenizer, isQuote)
	chk.Equal("", scanned)
	chk.Error(err)
	//
	str = "'value"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("", scanned)
	chk.Error(err)
	//
	str = "\r\n'value"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("", scanned)
	chk.Error(err)
	//
	str = "value"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("value", scanned)
	chk.NoError(err)
	//
	str = "'value'"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("value", scanned)
	chk.NoError(err)
	//
	str = "...value"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("...value", scanned)
	chk.NoError(err)
	//
	str = "value\n"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("value", scanned)
	chk.NoError(err)
	//
	str = "value   \n"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("value", scanned)
	chk.NoError(err)
	//
	str = "value   with\t\tspaces\n"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("value   with\t\tspaces", scanned)
	chk.NoError(err)
	//
	str = "value   with\t\tspaces and punctuation!!!\n"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.Equal("value   with\t\tspaces and punctuation!!!", scanned)
	chk.NoError(err)
}

func TestScanValue_email(t *testing.T) {
	chk := assert.New(t)
	//
	isQuote := func(r rune) bool {
		return r == '\'' || r == '"' || r == '`'
	}
	//
	str := "foo@example.com"
	tokenizer := parser.NewTokenizer(str)
	scanned, err := parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("foo@example.com", scanned)
	//
	str = "foo..bar@example.com"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("foo..bar@example.com", scanned)
	//
	str = "'\"foo..bar\"@example.com'"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("\"foo..bar\"@example.com", scanned)
	//
	str = "Foo Bar <foo.bar@example.com>"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("Foo Bar <foo.bar@example.com>", scanned)
}

func TestScanValue_delimited(t *testing.T) {
	chk := assert.New(t)
	//
	isQuote := func(r rune) bool {
		return r == '\'' || r == '"' || r == '`'
	}
	//
	str := "@ foo @ bar"
	tokenizer := parser.NewTokenizer(str)
	scanned, err := parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("@ foo @ bar", scanned)
	//
	str = "@ foo @ bar \t\t"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("@ foo @ bar", scanned)
	//
	str = "@ \"foo\" @ \"bar\" @"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("@ \"foo\" @ \"bar\" @", scanned)
	//
	str = "@ \"foo\" @ \"bar\" @\t\t\r\n"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("@ \"foo\" @ \"bar\" @", scanned)
}

func TestScanWhileType(t *testing.T) {
	chk := assert.New(t)
	//
	str := "hello\n\r\t\t\r\nworld!"
	tokenizer := parser.NewTokenizer(str)
	chk.NotNil(tokenizer)
	//
	tok, typ := tokenizer.Next()
	chk.Equal("hello", tok)
	chk.Equal(parser.TokenAlphaNum, typ)
	//
	scanned, err := parser.ScanWhileType(tokenizer, parser.TokenNewline|parser.TokenWhiteSpace)
	chk.Equal("\n\r\t\t\r\n", scanned)
	chk.NoError(err)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("world", tok)
	chk.Equal(parser.TokenAlphaNum, typ)
	//
	tok, typ = tokenizer.Next()
	chk.Equal("!", tok)
	chk.Equal(parser.TokenPunct, typ)
	//
	chk.Equal(true, tokenizer.Eof())
}

func TestScanValue_cron(t *testing.T) {
	chk := assert.New(t)
	//
	isQuote := func(r rune) bool {
		return r == '\'' || r == '"' || r == '`'
	}
	//
	str := "6 0 * * *"
	tokenizer := parser.NewTokenizer(str)
	scanned, err := parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("6 0 * * *", scanned)
	//
	str = "6 0 * * * "
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("6 0 * * *", scanned)
	//
	str = "* 6 0 * * *"
	tokenizer = parser.NewTokenizer(str)
	scanned, err = parser.ScanValue(tokenizer, isQuote)
	chk.NoError(err)
	chk.Equal("* 6 0 * * *", scanned)
}
