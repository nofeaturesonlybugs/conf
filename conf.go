package conf

import (
	"os"

	"github.com/nofeaturesonlybugs/conf/parser"
	"github.com/nofeaturesonlybugs/errors"
	"github.com/nofeaturesonlybugs/set"
)

// Conf is a parsed configuration.
type Conf struct {
	parsed parser.Parsed
}

// File returns a Conf type by reading and parsing the given file.
func File(file string) (*Conf, error) {
	handle, err := os.Open(file)
	if err != nil {
		return nil, errors.Go(err)
	}
	defer handle.Close()
	//
	parser := parser.NewParser()
	parsed, err := parser.ParseReader(handle)
	if err != nil {
		return nil, errors.Go(err)
	}
	//
	return &Conf{parsed}, nil
}

// String returns a Conf type by parsing the given string of configuration data.
func String(s string) (*Conf, error) {
	parser := parser.NewParser()
	parsed, err := parser.Parse(s)
	if err != nil {
		return nil, errors.Go(err)
	}
	//
	return &Conf{parsed}, nil
}

// fill populates target either ByTag or ByFieldName as determined by tag == "".
func (me *Conf) fill(target interface{}, tag string) error {
	if me == nil {
		return errors.NilReceiver()
	}
	//
	value := set.V(target)
	//
	var fields []set.Field
	if tag == "" {
		fields = value.Fields()
	} else {
		fields = value.FieldsByTag(tag)
	}
	//
	m := me.parsed.Map()
	globalSection, getter := set.MapGetter(m[""][len(m[""])-1]), set.MapGetter(m)
	scalars := map[string]set.Getter{}
	for _, field := range fields {
		if field.Value.IsScalar {
			if tag == "" {
				scalars[field.Field.Name] = globalSection
			} else {
				scalars[field.TagValue] = globalSection
			}
		}
	}
	//
	fn := set.GetterFunc(func(name string) interface{} {
		if scalarGetter, ok := scalars[name]; ok {
			return scalarGetter.Get(name)
		}
		return getter.Get(name)
	})
	if tag == "" {
		return value.Fill(fn)
	}
	return value.FillByTag(tag, fn)
}

// Fill places the data from the configuration into the given target which should be
// a pointer to struct.
func (me *Conf) Fill(target interface{}) error {
	return me.fill(target, "")
}

// FillByTag places the data from the configuration into the given target which should be
// a pointer to struct.
func (me *Conf) FillByTag(tag string, target interface{}) error {
	return me.fill(target, tag)
}
