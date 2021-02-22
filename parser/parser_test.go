package parser_test

import (
	"fmt"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/conf/parser"
	"github.com/nofeaturesonlybugs/errors"
)

func TestParser(t *testing.T) {
	chk := assert.New(t)
	//
	str := `
# Ignore this comment.
key = value
key1 = value1
`
	parsed, err := parser.DefaultParser.Parse(str)
	chk.NoError(err)
	chk.NotNil(parsed)
	//
	chk.NotNil(parsed[""])
	chk.Nil(parsed["hello"])
}

func TestParserNoWhiteSpace(t *testing.T) {
	chk := assert.New(t)
	//
	s := `
[service]
name=examplesvc
label=Example Service
description=This is an example service.
`
	//
	parsed, err := parser.DefaultParser.Parse(s)
	chk.NoError(err)
	chk.NotNil(parsed["service"])
	chk.Equal("examplesvc", parsed["service"].Last["name"].Last)
	chk.Equal("Example Service", parsed["service"].Last["label"].Last)
	chk.Equal("This is an example service.", parsed["service"].Last["description"].Last)
}

func TestParserEmptyValue(t *testing.T) {
	chk := assert.New(t)
	//
	s := `
; a comment
empty =
; new comment
`
	//
	parsed, err := parser.DefaultParser.Parse(s)
	chk.NoError(err)
	chk.NotNil(parsed[""])
	chk.Equal("", parsed[""].Last["empty"].Last)
}

func TestParserSections(t *testing.T) {
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
	parsed, err := parser.DefaultParser.Parse(str)
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

func TestParserSectionsSlices(t *testing.T) {
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
	parsed, err := parser.DefaultParser.Parse(str)
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

func TestParserSectionsSlicesValueSlices(t *testing.T) {
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
	parsed, err := parser.DefaultParser.Parse(str)
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

func TestParserComment(t *testing.T) {
	chk := assert.New(t)
	parsed, err := parser.DefaultParser.Parse(`
		# Hello, World!
	`)
	chk.NoError(err)
	chk.NotNil(parsed)
}

func TestParserSection(t *testing.T) {
	chk := assert.New(t)
	{ // Expected
		parsed, err := parser.DefaultParser.Parse(`
		[ hello ]
	`)
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.NotNil(parsed["hello"])
	}
	{ // Expected with spaces
		parsed, err := parser.DefaultParser.Parse(`
		[ hello world ]
	`)
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.NotNil(parsed["hello world"])
	}
	{ // Expected with punctuation
		parsed, err := parser.DefaultParser.Parse(`
		[ hello.world ]
	`)
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.NotNil(parsed["hello.world"])
	}
	{ // Error with unexpected token
		parsed, err := parser.DefaultParser.Parse(`
		[ hello . ]
	`)
		chk.Error(err)
		chk.NotNil(parsed)
	}
	{ // Error with punctuation not followed by alpha num
		parsed, err := parser.DefaultParser.Parse(`
		[ hello.. ]
	`)
		chk.Error(err)
		chk.NotNil(parsed)
	}
}

func TestParserSectionMultiple(t *testing.T) {
	chk := assert.New(t)
	{ // Expected
		parsed, err := parser.DefaultParser.Parse(`
		[ hello ]
		message = "Hi there!"
		message = "Ok..."

		[ hello ]
		message = "Fred says, 'Good day!'"
		message = "This message
spans
multiple lines!"
	`)
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.NotNil(parsed["hello"])
		chk.Equal(2, len(parsed["hello"].Slice))
		a, b := parsed["hello"].Slice[0], parsed["hello"].Slice[1]
		chk.Equal(parsed["hello"].Last, b)
		//
		chk.Equal(2, len(a["message"].Slice))
		chk.Equal("Hi there!", a["message"].Slice[0])
		chk.Equal("Ok...", a["message"].Slice[1])
		//
		chk.Equal(2, len(b["message"].Slice))
		chk.Equal("Fred says, 'Good day!'", b["message"].Slice[0])
		chk.Equal("This message\nspans\nmultiple lines!", b["message"].Slice[1])
	}
}

func TestParserReader(t *testing.T) {
	chk := assert.New(t)
	{ // Expected
		parsed, err := parser.DefaultParser.ParseReader(strings.NewReader(`
		[ hello ]
		message = "Hi there!"
		message = "Ok..."

		[ hello ]
		message = "Fred says, 'Good day!'"
		message = "This message
spans
multiple lines!"
	`))
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.NotNil(parsed["hello"])
		chk.Equal(2, len(parsed["hello"].Slice))
		a, b := parsed["hello"].Slice[0], parsed["hello"].Slice[1]
		chk.Equal(parsed["hello"].Last, b)
		//
		chk.Equal(2, len(a["message"].Slice))
		chk.Equal("Hi there!", a["message"].Slice[0])
		chk.Equal("Ok...", a["message"].Slice[1])
		//
		chk.Equal(2, len(b["message"].Slice))
		chk.Equal("Fred says, 'Good day!'", b["message"].Slice[0])
		chk.Equal("This message\nspans\nmultiple lines!", b["message"].Slice[1])
	}
}

func TestParserKeyValue(t *testing.T) {
	chk := assert.New(t)
	//
	{
		parsed, err := parser.DefaultParser.Parse(`
		key = value1
		key = value2
		other = something else   
	`)
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.Equal(1, len(parsed[""].Slice))
		section := parsed[""].Last
		chk.NotNil(section)
		//
		chk.NotNil(section["key"].Slice)
		chk.Equal(2, len(section["key"].Slice))
		chk.Equal("value2", section["key"].Last)
		chk.Equal("value1", section["key"].Slice[0])
		chk.Equal("value2", section["key"].Slice[1])
		//
		chk.NotNil(section["other"].Slice)
		chk.Equal(1, len(section["other"].Slice))
		chk.Equal("something else", section["other"].Last)
		chk.Equal("something else", section["other"].Slice[0])
	}
	{ // key with punctuation
		parsed, err := parser.DefaultParser.Parse(`
		key.1 = value1
		key.2 = value2
	`)
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.Equal(1, len(parsed[""].Slice))
		section := parsed[""].Last
		chk.NotNil(section)
		//
		chk.NotNil(section["key.1"].Slice)
		chk.Equal(1, len(section["key.1"].Slice))
		chk.Equal("value1", section["key.1"].Last)
		chk.Equal("value1", section["key.1"].Slice[0])
		chk.Equal(1, len(section["key.2"].Slice))
		chk.Equal("value2", section["key.2"].Last)
		chk.Equal("value2", section["key.2"].Slice[0])
	}
	{ // key with whitespace not followed by alphanum or assign
		parsed, err := parser.DefaultParser.Parse(`
		key   . = value1
	`)
		chk.Error(err)
		chk.NotNil(parsed)
	}
	{ // key with bad punctuation
		parsed, err := parser.DefaultParser.Parse(`
		key. = value1
	`)
		chk.Error(err)
		chk.NotNil(parsed)
	}
}

func TestParserKeyValueQuotes(t *testing.T) {
	chk := assert.New(t)
	//
	{
		parsed, err := parser.DefaultParser.Parse(`
		key = "value1"
		key = "  value2  "
		other = "something 'else   '  "
	`)
		chk.NoError(err)
		chk.NotNil(parsed)
		chk.Equal(1, len(parsed[""].Slice))
		section := parsed[""].Last
		chk.NotNil(section)
		//
		chk.NotNil(section["key"].Slice)
		chk.Equal(2, len(section["key"].Slice))
		chk.Equal("  value2  ", section["key"].Last)
		chk.Equal("value1", section["key"].Slice[0])
		chk.Equal("  value2  ", section["key"].Slice[1])
		//
		chk.NotNil(section["other"].Slice)
		chk.Equal(1, len(section["other"].Slice))
		chk.Equal("something 'else   '  ", section["other"].Last)
		chk.Equal("something 'else   '  ", section["other"].Slice[0])
	}
}

func ExampleParser_globalSection() {
	s := `
hello = world
foo = bar
`
	//
	if parsed, err := parser.DefaultParser.Parse(s); err != nil {
		fmt.Println(err)
	} else if _, ok := parsed[""]; ok {
		fmt.Printf("The global section.\n")
		fmt.Printf("hello = %v\n", parsed[""].Last["hello"].Last)
		fmt.Printf("foo = %v", parsed[""].Last["foo"].Last)
	} else {
		fmt.Println("Global section not found.")
	}
	// Output: The global section.
	// hello = world
	// foo = bar
}

func ExampleParser_string() {
	s := `
[main]
hello = world
foo = bar
`
	//
	if parsed, err := parser.DefaultParser.Parse(s); err != nil {
		fmt.Println(err)
	} else if _, ok := parsed["main"]; ok {
		fmt.Printf("Section `main` exists.\n")
		fmt.Printf("hello = %v\n", parsed["main"].Last["hello"].Last)
		fmt.Printf("foo = %v\n", parsed["main"].Last["foo"].Last)
	} else {
		fmt.Println("Section `main` not found.")
	}
	// Output: Section `main` exists.
	// hello = world
	// foo = bar
}

func ExampleParser_reader() {
	s := `
[main]
hello = world
foo = bar
`
	reader := strings.NewReader(s)
	//
	if parsed, err := parser.DefaultParser.ParseReader(reader); err != nil {
		fmt.Println(err)
	} else if _, ok := parsed["main"]; ok {
		fmt.Printf("Section `main` exists.\n")
		fmt.Printf("hello = %v\n", parsed["main"].Last["hello"].Last)
		fmt.Printf("foo = %v", parsed["main"].Last["foo"].Last)
	} else {
		fmt.Println("Section `main` not found.")
	}
	// Output: Section `main` exists.
	// hello = world
	// foo = bar
}

func ExampleParser_options() {
	s := `
(main)
message ~ /Hello World!/
`
	//
	parser := &parser.Parser{
		Runes: parser.Runes{
			Assign:       []rune{'~'},
			Quote:        []rune{'/'},
			SectionOpen:  []rune{'('},
			SectionClose: []rune{')'},
		},
	}
	if parsed, err := parser.Parse(s); err != nil {
		fmt.Println(err)
	} else if _, ok := parsed["main"]; ok {
		fmt.Printf("Section `main` exists.\n")
		fmt.Printf("message = %v\n", parsed["main"].Last["message"].Last)
	} else {
		fmt.Println("Section `main` not found.")
	}
	// Output: Section `main` exists.
	// message = Hello World!
}

func TestParserReaderError(t *testing.T) {
	chk := assert.New(t)
	//
	{
		r := iotest.ErrReader(errors.Errorf("kaboom"))
		chk.NotNil(r)
		parsed, err := parser.DefaultParser.ParseReader(r)
		chk.Nil(parsed)
		chk.Error(err)
	}
}

func TestStateStrings(t *testing.T) {
	chk := assert.New(t)
	chk.Equal("None", parser.StateNone.String())
	chk.Equal("Comment", parser.StateComment.String())
	chk.Equal("Section", parser.StateSection.String())
	chk.Equal("Key", parser.StateKey.String())
	chk.Equal("Value", parser.StateValue.String())
	chk.Equal(true, strings.HasPrefix(parser.State(-10).String(), "Unknown"))
}
