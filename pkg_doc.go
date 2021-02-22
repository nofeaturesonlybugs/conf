// Package conf can parse configuration text and populate structs with the results
//
// Basic Syntax
//
// Example:
//
//	# Lines beginning with punctuation are comments.
//	; Also a comment!
//
//	# key=value pairs are assigned to the current section or the global section
//	# if no section has been defined yet.
//
//	key1 = value1
//	key with spaces = value with spaces!
//
// 	quote1 = 'This value is "quoted"'
// 	quote2 = "  whitespace
// 	is preserved in
// 	quoted values!   "
// 	quote3 = `backticks quote too!`
//
//	# Create a section with [ section-name ]
//	[ section ]
//	# This key value pair goes into a section named: section
//	key = value
//
//	[ sections can also have spaces ]
//	keys.can.have.punctuation.too = neat!
//
// Repetition
//
// Both keys and sections can be turned into lists by repeating them:
//	fruits = apples
//	fruits = oranges
//	fruits = bananas
//
//	[ color ]
//	name = red
//	rgb = ff0000
//
//	[ color ]
//	name = blue
//	rgb = 0000ff
//
//	[ color ]
//	name = green
//	rgb = 00ff00
//
// Fill
//
// Use Conf.Fill() and Conf.FillByTag() to populate parsed configuration into your structures.  Examples are provided below.
//
// Configuration EBNF
//
// Here lies the EBNF for configuration syntax:
//	conf
//		: line*
//		;
//
//	line
//		: comment
//		| section
//		| key_value
//		;
//
//	section
//		: '[' key ']'
//		;
//
//	key_value
//		: key '=' value
//		;
//
//	key
//		: [a-z0-9]+ key_extend*
//		;
//
//	key_extend
//		: punct [a-z0-9]+
//		: ws+ [a-z0-9]+
//		;
//
//	value
//		: ['] ~[']* [']
//		: ["] ~["]* ["]
//		: [`] ~[`]* [`]
//		: ~[rn]*
//		;
//
//	comment
//		: punct ~[rn]*
//		;
//
//	punct
//		: [!@#$%&*()_+=\|;:'"``,.<>/?~-^]
//		;
//
//	ws
//		: [ t]
//		;
package conf
