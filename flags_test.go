package comfyconf

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags_NewFlags(t *testing.T) {

	f := NewFlags()

	assert.Equal(t, "*comfyconf.Flags", reflect.TypeOf(f).String())
	assert.NotNil(t, f.args)
	assert.NotNil(t, f.parsed)
}

func TestFlags_NewFlagsWithCustomParser(t *testing.T) {

	f := NewFlagsWithCustomParser(func(arg string) (string, string) {
		return "", ""
	})

	assert.Equal(t, "*comfyconf.Flags", reflect.TypeOf(f).String())
	assert.NotNil(t, f.args)
	assert.NotNil(t, f.parsed)
}

func TestFlags_Init(t *testing.T) {

	f := prepareFlags([]string{"-t=Panzer"}, "=")
	assert.Nil(t, f.Init())

	assert.Equal(t, f.parsed["t"], "Panzer")
}

func TestFlags_get_Success(t *testing.T) {
	f := prepareFlags([]string{"--test=ArtifexHomicida", "-t2=drewoko"}, "=")
	assert.Nil(t, f.Init())

	r1, r2 := f.get("t", "test")

	assert.True(t, r2)
	assert.Equal(t, r1, "ArtifexHomicida")

	r1, r2 = f.get("t2", "test2")

	assert.True(t, r2)
	assert.Equal(t, r1, "drewoko")

	f = prepareFlags([]string{"--test Askar"}, " ")
	assert.Nil(t, f.Init())

	r1, r2 = f.get("t2", "test2")

	assert.False(t, r2)
	assert.Equal(t, r1, "")
}

func TestFlags_ParseInt_Success(t *testing.T) {
	f := prepareFlags([]string{"--test 1"}, " ")
	assert.Nil(t, f.Init())

	r1, r2 := f.ParseInt("t", "test")

	assert.True(t, r2)
	assert.Equal(t, r1, 1)

	f = prepareFlags([]string{"-t=The_Guy"}, "=")
	assert.Nil(t, f.Init())

	r1, r2 = f.ParseInt("t", "test")

	assert.False(t, r2)
	assert.Equal(t, r1, 0)
}

func TestFlags_ParseString_Success(t *testing.T) {
	f := prepareFlags([]string{"--test=Daimon"}, "=")
	assert.Nil(t, f.Init())

	r1, r2 := f.ParseString("t", "test")

	assert.True(t, r2)
	assert.Equal(t, r1, "Daimon")

	f = prepareFlags([]string{"-t2=judgegc"}, "=")
	assert.Nil(t, f.Init())

	r1, r2 = f.ParseString("test", "t")

	assert.Equal(t, r2, false)
	assert.Equal(t, r1, "")
}

func TestFlags_ParseBool_Success(t *testing.T) {
	f := prepareFlags([]string{"-t=true"}, "=")
	assert.Nil(t, f.Init())

	r1, r2 := f.ParseBool("t", "test")

	assert.True(t, r2)
	assert.True(t, r1)

	f = prepareFlags([]string{"-t2=false"}, "=")
	assert.Nil(t, f.Init())

	r1, r2 = f.ParseBool("t", "test")

	assert.False(t, r2)
	assert.False(t, r1)
}

func TestFlags_ParseExistence(t *testing.T) {
	f := prepareFlags([]string{"-t"}, "=")
	assert.Nil(t, f.Init())

	r11, r12 := f.ParseExistence("t", "test")

	assert.True(t, r11)
	assert.True(t, r12)

	r21, r22 := f.ParseExistence("t2", "test2")

	assert.False(t, r21)
	assert.True(t, r22)
}

func TestFlags_ParseSlice(t *testing.T) {
	f := prepareFlags([]string{"-t[]=false", "-t[]=2", "-t[1]=3", "--test1"}, "=")
	assert.Nil(t, f.Init())

	r1, isOk1 := f.ParseSlice("t", "test")

	assert.Len(t, r1, 3)
	assert.True(t, isOk1)
	assert.Equal(t, "2", r1[1])
}

func TestFlags_Parser_FullFormat(t *testing.T) {

	str1, str2 := DefaultFlagsParser("--test=Koddi", "=")

	assert.Equal(t, str1, "test")
	assert.Equal(t, str2, "Koddi")
}

func TestFlags_Parser_ShortFormat(t *testing.T) {

	str1, str2 := DefaultFlagsParser("-t=Koddi", "=")

	assert.Equal(t, str1, "t")
	assert.Equal(t, str2, "Koddi")
}

func TestFlags_Parser_ShortFormat_LongString(t *testing.T) {

	str1, str2 := DefaultFlagsParser("-t drewoko loves Koddi", " ")

	assert.Equal(t, str1, "t")
	assert.Equal(t, str2, "drewoko loves Koddi")
}

func TestFlags_Parser_NotFound(t *testing.T) {

	str1, str2 := DefaultFlagsParser("-test=Don.Gray", "=")

	assert.NotEqual(t, str1, "t")
	assert.NotEqual(t, str2, "pugaman")
}

func prepareFlags(args []string, assignment string) *Flags {
	return &Flags{
		args,
		func(arg string) (string, string) {
			return DefaultFlagsParser(arg, assignment)
		},
		make(map[string]string),
		make(map[string][]interface{}),
	}
}
