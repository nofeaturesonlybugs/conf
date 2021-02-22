package parser

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
