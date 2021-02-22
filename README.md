[![Documentation](https://godoc.org/github.com/nofeaturesonlybugs/conf?status.svg)](http://godoc.org/github.com/nofeaturesonlybugs/conf)
[![Go Report Card](https://goreportcard.com/badge/github.com/nofeaturesonlybugs/conf)](https://goreportcard.com/report/github.com/nofeaturesonlybugs/conf)
[![Build Status](https://travis-ci.com/nofeaturesonlybugs/conf.svg?branch=master)](https://travis-ci.com/nofeaturesonlybugs/conf)
[![codecov](https://codecov.io/gh/nofeaturesonlybugs/conf/branch/master/graph/badge.svg)](https://codecov.io/gh/nofeaturesonlybugs/conf)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Package `conf` parses configuration data and populates a Go `struct` with the results.  

The package documentation contains the configuration EBNF but here are some syntax examples:

## Comments  
```ini
# Lines beginning with punctuation are comments.
; Any punctuation except [] can begin a comment.
^ Just know that syntax highlighters won't always know what to do!
```

## Keys and Values  
```
key = value
keys can have spaces = values run until the end of the line
keys.can.have.punctuation = so can values!

values can be quoted = 'This value is quoted!'
quotes preserve whitespace = `
This value
spans
4 lines!`

# Quotes are only recognized if they are matching runes at the beginning and end of a value.
not quoted = asdf '" fdsa

# Quotes can be used to add quote characters to the parsed value.
quotes picked up = `"Hello!" said Dave.`
```

## Invalid Keys  
```
# Invalid - whitespace must be followed by [a-z0-9]
key . = oops

# Invalid - multiple punctuation not allowed
key.. = oops again

# Invalid - punctuation must join [a-z0-9] to [a-z0-9] without whitespace
key . key_extra = nope
```

## Sections  
Create a section with square brackets:  
```
# These values are global
key1 = value1
key2 = value2

# This is a section named: color
[ color ]
name = red
rgb = #ff0000
```

## Lists Require No Special Syntax  
Create lists in your configuration by simply repeating `keys` or `sections`:  
```
fruits = apples
fruits = oranges
fruits = bananas

[ color ]
name = red
rgb = ff0000

[ color ]
name = blue
rgb = 0000ff

[ color ]
name = green
rgb = 00ff00
```

## Easily Populate a Conf Struct  
Create a `struct` matching the configuration.  

This `struct` is for the configuration just prior:  
```go
type T struct {
    Fruits []string `conf:"fruits"`

    Color []struct {
        Name string `conf:"name"`
        Rgb  string `conf:"rgb"`
    } `conf:"color"`
}
```

Consume the configuration either as a `string` or `file`.  
```go
conf, err := conf.String(s)
if err != nil {
    fmt.Println(err.Error())
}
```

We're using `FillByTag` to map the lowercase config values to our uppercase `struct` members:  
```go
var t T
if err = conf.FillByTag("conf", &t); err != nil {
    fmt.Println(err.Error())
}
fmt.Printf("%v\n", t.Fruits)
for _, color := range t.Color {
    fmt.Printf("%v %v\n", color.Name, color.Rgb)
}
```
