package parser_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/conf/parser"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	chk := assert.New(t)
	//
	{
		s := `
			[section]
			a = A
			b = B

			[section]
			a = AA
			b = BB
		`
		parsed, err := parser.DefaultParser.Parse(s)
		chk.NoError(err)
		chk.NotNil(parsed)
		mapped := parsed.Map()
		chk.NotNil(mapped)
		sectionSlice := mapped["section"]
		chk.NotNil(sectionSlice)
		chk.Equal(2, len(sectionSlice))
		//
		first := sectionSlice[0]
		chk.NotNil(first)
		chk.Equal(1, len(first["a"]))
		chk.Equal("A", first["a"][0])
		chk.Equal("B", first["b"][0])
		chk.Equal(1, len(first["b"]))
		//
		second := sectionSlice[1]
		chk.NotNil(second)
		chk.Equal("AA", second["a"][0])
		chk.Equal("BB", second["b"][0])
	}
}
