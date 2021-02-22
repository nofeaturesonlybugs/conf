package conf_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/conf"
)

func TestConf_Spaces(t *testing.T) {
	chk := assert.New(t)
	s := `
	# Lines beginning with punctuation are comments.
; Also a comment!

# key=value pairs are assigned to the current section or the global section
# if no section has been defined yet.

	key with spaces = value with spaces
	key with spaces = other value with spaces!

	[ section can have spaces too ]
	section key 	=		 section value

	`
	type T struct {
		KeyWithSpaces   string   `conf:"key with spaces"`
		SliceWithSpaces []string `conf:"key with spaces"`

		// Tests [section]
		Section struct {
			SectionValue string `conf:"section key"`
		} `conf:"section can have spaces too"`
	}
	testIt := func(t *T) {
		// pointers created
		chk.NotNil(t)
		// global config
		chk.Equal("other value with spaces!", t.KeyWithSpaces)
		chk.Equal("value with spaces", t.SliceWithSpaces[0])
		chk.Equal("other value with spaces!", t.SliceWithSpaces[1])
		// section config
		chk.Equal("section value", t.Section.SectionValue)
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

func ExampleConf_Fill() {
	type T struct {
		First  string
		Second string

		Numbers struct {
			N []int
		}
	}
	s := `
First = I'm first!
Second = ...and I'm second.

[ Numbers ]
N = 42
N = 3
	`
	conf, err := conf.String(s)
	if err != nil {
		fmt.Println(err.Error())
	}
	var t T
	if err = conf.Fill(&t); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("First= %v\n", t.First)
	fmt.Printf("Second= %v\n", t.Second)
	fmt.Printf("N[0]= %v\n", t.Numbers.N[0])
	fmt.Printf("N[1]= %v\n", t.Numbers.N[1])
	// Output: First= I'm first!
	// Second= ...and I'm second.
	// N[0]= 42
	// N[1]= 3
}

func ExampleConf_Fill_withRepetition() {
	type T struct {
		Fruits []string

		Color []struct {
			Name string
			Rgb  string
		}
	}
	s := `
Fruits = apples
Fruits = oranges
Fruits = bananas

[ Color ]
Name = red
Rgb = ff0000

[ Color ]
Name = blue
Rgb = 0000ff

[ Color ]
Name = green
Rgb = 00ff00
		`
	conf, err := conf.String(s)
	if err != nil {
		fmt.Println(err.Error())
	}
	var t T
	if err = conf.Fill(&t); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v\n", t.Fruits)
	for _, color := range t.Color {
		fmt.Printf("%v %v\n", color.Name, color.Rgb)
	}

	// Output: [apples oranges bananas]
	// red ff0000
	// blue 0000ff
	// green 00ff00
}

func ExampleConf_FillByTag() {
	type T struct {
		First  string `conf:"first"`
		Second string `conf:"second"`

		Numbers struct {
			N []int `conf:"n"`
		} `conf:"numbers"`
	}
	s := `
first = I'm first!
second = ...and I'm second.

[ numbers ]
n = 42
n = 3
	`
	conf, err := conf.String(s)
	if err != nil {
		fmt.Println(err.Error())
	}
	var t T
	if err = conf.FillByTag("conf", &t); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("First= %v\n", t.First)
	fmt.Printf("Second= %v\n", t.Second)
	fmt.Printf("N[0]= %v\n", t.Numbers.N[0])
	fmt.Printf("N[1]= %v\n", t.Numbers.N[1])
	// Output: First= I'm first!
	// Second= ...and I'm second.
	// N[0]= 42
	// N[1]= 3
}

func ExampleConf_FillByTag_withRepetition() {
	type T struct {
		Fruits []string `conf:"fruits"`

		Color []struct {
			Name string `conf:"name"`
			Rgb  string `conf:"rgb"`
		} `conf:"color"`
	}
	s := `
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
		`
	conf, err := conf.String(s)
	if err != nil {
		fmt.Println(err.Error())
	}
	var t T
	if err = conf.FillByTag("conf", &t); err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Printf("%#v\n", conf) //TODO RM
	// fmt.Printf("%#v\n", t)    //TODO RM
	fmt.Printf("%v\n", t.Fruits)
	for _, color := range t.Color {
		fmt.Printf("%v %v\n", color.Name, color.Rgb)
	}

	// Output: [apples oranges bananas]
	// red ff0000
	// blue 0000ff
	// green 00ff00
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
