package parser_test

import (
	"fmt"
	"strings"

	"github.com/nofeaturesonlybugs/conf/parser"
)

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
		fmt.Printf("foo = %v", parsed["main"].Last["foo"].Last)
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
