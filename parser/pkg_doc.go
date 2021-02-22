// Package parser is the configurable parser for the conf package.
//
// The primary type exported by this package is Parser.  You can use the default parser:
//	var DefaultParser = Parser{
//		Runes: Runes{
//			Assign:       []rune{'='},
//			Quote:        []rune{'\'', '"', '`'},
//			SectionOpen:  []rune{'['},
//			SectionClose: []rune{']'},
//		},
//	}
//
// or create your own:
//	myParser := parser.Parser{
//		Runes: Runes{
//			Assign:       []rune{'='},
//			Quote:        []rune{'`'}, // Backtick quotes only.
//			SectionOpen:  []rune{'{'}, // Curly bracket sections.
//			SectionClose: []rune{'}'},
//		},
//	}
//
// The given runes should be mutually exclusive sets when creating a Parser; the behavior for disregarding
// this rule is undefined.
//
// The Parsed Type
//
// When parsing succeeds a type Parsed is returned.  It is a map[string]*SectionBlock.  Semantically it is
// a map["section-name"]*SectionBlock.
//
// Global Section
//
// When parsing the input the default starting section is an unnammed section represented by an empty string.
// Once a section name is parsed all configuration goes into that section and any following sections.  No more
// configuration can be placed in the global section.
//
// See example Global Section under Parser.
//
// Repeated Sections
//
// A SectionBlock contains two members:
//	Last Section
//	Slice []Section
//
// The Slice member will contain all sections with the same name in the order they were encountered.  The Last
// member contains the last setion encountered.
//
// The example configuration:
//	[domain]
//	listen = 0.0.0.0
//
//	[domain]
//	listen = example.com
//
// Creates a SectionBlock where:
//	Slice[0] is the section where listen = 0.0.0.0
//	Slice[1] is the section where listen = example.com
//	Last is the same section as Slice[1]
//
// Section and its Values
//
// The Section and Value types repeat some of concepts encountered already.  A Section is a map[string]Value or
// semantically a map["key"]Value.  A Value contains two members:
//	Last string
//	Slice []string
//
// The Slice member will contain all key=value lines with the same key-name in the order they were encountered.
// The Last member contains the last key=value encountered for a specific key-name.
//
// The example configuration:
//	[domains]
//	listen = 0.0.0.0
//	listen = example.com
//
// Creates a Section where:
//	Section["listen"].Slice[0] = 0.0.0.0
//	Section["listen"].Slice[1] = example.com
//	Section["listen"].Last is the same as Section["listen"].Slice[1]
//
// The End Result
//
// The end result is a convenient configuration syntax that allows repeated sections and repeated key=values
// without adding any special syntax.
package parser
