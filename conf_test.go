package conf_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	// "github.com/nofeaturesonlybugs/assert" // TODO RM
	"github.com/nofeaturesonlybugs/conf"
)

func TestConf_FillByTag(t *testing.T) {
	chk := assert.New(t)
	//
	s := `
	global-value = I'm a global!
global-value = I'm also global!

global with spaces = yay!

[section]
section-value = I'm a section value!

[schedule]
archive = 0 6 * * *
archive = 0 */2 * * *
`
	type T struct {
		Global           string `conf:"global-value"`
		GlobalWithSpaces string `conf:"global with spaces"`

		// Tests [section]
		Section struct {
			SectionValue string `conf:"section-value"`
		} `conf:"section"`
		SectionPtr *struct {
			SectionValue string `conf:"section-value"`
		} `conf:"section"`
		SectionSlicePtr []*struct {
			SectionValue string `conf:"section-value"`
		} `conf:"section"`
		PtrSectionSlicePtr *[]*struct {
			SectionValue string `conf:"section-value"`
		} `conf:"section"`

		// Tests [schedule]
		Schedule struct {
			Archive []string `conf:"archive"`
		} `conf:"schedule"`
	}
	testIt := func(t *T) {
		// pointers created
		chk.NotNil(t)
		chk.NotNil(t.SectionPtr)
		chk.NotNil(t.SectionSlicePtr)
		chk.NotNil(t.PtrSectionSlicePtr)
		// global config
		chk.Equal("I'm also global!", t.Global)
		chk.Equal("yay!", t.GlobalWithSpaces)
		// section config
		chk.Equal("I'm a section value!", t.Section.SectionValue)
		chk.Equal("I'm a section value!", t.SectionPtr.SectionValue)
		// slice and ptr-to-slice sections
		chk.Equal(1, len(t.SectionSlicePtr))
		chk.Equal("I'm a section value!", t.SectionSlicePtr[0].SectionValue)
		chk.Equal(1, len(*t.PtrSectionSlicePtr))
		chk.Equal("I'm a section value!", (*t.PtrSectionSlicePtr)[0].SectionValue)
		//
		chk.Equal(2, len(t.Schedule.Archive))
		chk.Equal("0 6 * * *", t.Schedule.Archive[0])
		chk.Equal("0 */2 * * *", t.Schedule.Archive[1])
	}

	{
		tmpfile, err := ioutil.TempFile("", "gotest")
		chk.NoError(err)
		defer os.Remove(tmpfile.Name())
		_, err = tmpfile.Write([]byte(s))
		chk.NoError(err)
		err = tmpfile.Close()
		chk.NoError(err)
		conf, err := conf.File(tmpfile.Name())
		chk.NoError(err)
		chk.NotNil(conf)
		var t *T
		err = conf.FillByTag("conf", &t)
		chk.NoError(err)
		testIt(t)
	}
	{
		conf, err := conf.String(s)
		chk.NoError(err)
		chk.NotNil(conf)
		var t *T
		err = conf.FillByTag("conf", &t)
		chk.NoError(err)
		testIt(t)
	}
}

func TestConf_Fill(t *testing.T) {
	chk := assert.New(t)
	//
	s := `
	Global = I'm a global!
Global = I'm also global!

[Section]
Value = I'm a section value!

[Schedule]
Archive = 0 6 * * *
Archive = 0 */2 * * *
`
	type T struct {
		Global string

		// Tests [section]
		Section struct {
			Value string
		}

		// Tests [schedule]
		Schedule struct {
			Archive []string
		}
	}
	testIt := func(t *T) {
		// pointers created
		chk.NotNil(t)
		// global config
		chk.Equal("I'm also global!", t.Global)
		// section config
		chk.Equal("I'm a section value!", t.Section.Value)
		//
		chk.Equal(2, len(t.Schedule.Archive))
		chk.Equal("0 6 * * *", t.Schedule.Archive[0])
		chk.Equal("0 */2 * * *", t.Schedule.Archive[1])
	}

	{
		tmpfile, err := ioutil.TempFile("", "gotest")
		chk.NoError(err)
		defer os.Remove(tmpfile.Name())
		_, err = tmpfile.Write([]byte(s))
		chk.NoError(err)
		err = tmpfile.Close()
		chk.NoError(err)
		conf, err := conf.File(tmpfile.Name())
		chk.NoError(err)
		chk.NotNil(conf)
		var t *T
		err = conf.Fill(&t)
		chk.NoError(err)
		testIt(t)
	}
	{
		conf, err := conf.String(s)
		chk.NoError(err)
		chk.NotNil(conf)
		var t *T
		err = conf.Fill(&t)
		chk.NoError(err)
		testIt(t)
	}
}

func TestConf_parseError(t *testing.T) {
	chk := assert.New(t)
	//
	s := `
	Global  I'm a global!
Global = I'm also global!
`
	{
		tmpfile, err := ioutil.TempFile("", "gotest")
		chk.NoError(err)
		defer os.Remove(tmpfile.Name())
		_, err = tmpfile.Write([]byte(s))
		chk.NoError(err)
		err = tmpfile.Close()
		chk.NoError(err)
		conf, err := conf.File(tmpfile.Name())
		chk.Error(err)
		chk.Nil(conf)
	}
	{
		conf, err := conf.String(s)
		chk.Error(err)
		chk.Nil(conf)
	}
}

func TestConf_Fill_nilReceiver(t *testing.T) {
	chk := assert.New(t)
	//
	var conf *conf.Conf
	err := conf.Fill(nil)
	chk.Error(err)
}

func TestConf_File_doesNotExist(t *testing.T) {
	chk := assert.New(t)
	//
	_, err := conf.File("asldjflaksdjflaksjflasjdf")
	chk.Error(err)
}
