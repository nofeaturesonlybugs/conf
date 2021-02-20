package parser

// OptAssign defines the delimiter(s) for assigning a value to a key.
type OptAssign []rune

// OptQuote defines the rune(s) for a quoted value.
type OptQuote []rune

// OptSectionRunes defines the open-close pair that defines a section.
type OptSectionRunes [2]rune

// Options gathers and summarizes the various options.
type Options struct {
	sectionOpen  rune
	sectionClose rune
	isAssign     IsAssign
	isQuote      IsQuote
}

// NewOptions creates a new Options type.
func NewOptions(options ...interface{}) *Options {
	//
	// Default values.
	assign := OptAssign([]rune{'='})
	quote := OptQuote([]rune{'\'', '"', '`'})
	sectionOpen, sectionClose := '[', ']'
	//
	// Override defaults with whatever was passed.
	for _, option := range options {
		switch value := option.(type) {
		case OptAssign:
			assign = value
		case OptQuote:
			quote = value
		case OptSectionRunes:
			// TODO Test length?
			sectionOpen, sectionClose = value[0], value[1]
			// default: // TODO
			// TODO ERROR
		}
	}
	// O(1) lookup for assignment runes.
	assignMap := make(map[rune]struct{})
	for _, r := range assign {
		assignMap[r] = struct{}{}
	}
	isAssign := func(r rune) bool {
		_, ok := assignMap[r]
		return ok
	}
	// O(1) lookup for quote runes.
	quoteMap := make(map[rune]struct{})
	for _, r := range quote {
		quoteMap[r] = struct{}{}
	}
	isQuote := func(r rune) bool {
		_, ok := quoteMap[r]
		return ok
	}
	//
	rv := &Options{
		isAssign:     isAssign,
		isQuote:      isQuote,
		sectionOpen:  sectionOpen,
		sectionClose: sectionClose,
	}
	return rv
}
