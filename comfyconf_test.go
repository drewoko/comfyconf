package comfyconf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func prepareConf(args []string, assignment string) *Conf {
	return New(prepareFlags(args, assignment))
}

func TestConf_AddMiddleware(t *testing.T) {
	conf := prepareConf([]string{}, " ")
	conf.AddMiddleware(&Env{})

	assert.Len(t, conf.middleware, 2)
}

func TestConf_Parse(t *testing.T) {
	conf := prepareConf([]string{}, " ")
	err := conf.Parse()

	assert.Nil(t, err)
}

func TestConf_TestDefaultValue_WithoutParseCall_Full(t *testing.T) {
	conf := prepareConf([]string{"-t=Nixel"}, "=")

	t1 := conf.String("test", "t", "Ультвэ", "Basic description")

	assert.Equal(t, *t1, "Ультвэ")
}

func TestConf_ParseString_Assignment1_Full(t *testing.T) {

	conf := prepareConf([]string{"-t=mofa"}, "=")

	t1 := conf.String("test", "t", "ews", "Basic description")
	t2 := conf.String("test2", "t2", "c0deum", "Basic description 2")

	assert.Nil(t, conf.Parse())

	assert.Equal(t, *t1, "mofa")
	assert.Equal(t, *t2, "c0deum")
}

func TestConf_ParseString_Assignment2_Full(t *testing.T) {

	conf := prepareConf([]string{"-t Jare", "--test2 Hi, My tagName is FRAG"}, " ")

	t1 := conf.String("test", "t", "Hi, There", "Basic description")
	t2 := conf.String("test2", "t2", "Hi, There 2", "Basic description 2")

	assert.Nil(t, conf.Parse())

	assert.Equal(t, *t1, "Jare")
	assert.Equal(t, *t2, "Hi, My tagName is FRAG")
}

func TestConf_ParseBool_Full(t *testing.T) {

	conf := prepareConf([]string{"-t=true", "--test2=false"}, "=")

	t1 := conf.Bool("test", "t", false, "Basic description")
	t2 := conf.Bool("test2", "t2", true, "Basic description 2")
	t3 := conf.Bool("test3", "t3", true, "Basic description 3")

	assert.Nil(t, conf.Parse())

	assert.True(t, *t1)
	assert.False(t, *t2)
	assert.True(t, *t3)
}

func TestConf_ParseExistence_Full(t *testing.T) {

	conf := prepareConf([]string{"-t", "--test", "-t2"}, "=")

	t1 := conf.Exist("test", "t", "Basic description")
	t2 := conf.Exist("test2", "t2", "Basic description 2")
	t3 := conf.Exist("test3", "t3", "Basic description 3")

	assert.Nil(t, conf.Parse())

	assert.True(t, *t1)
	assert.True(t, *t2)
	assert.False(t, *t3)
}

func TestConf_Slice(t *testing.T) {

	conf := prepareConf([]string{"-t[]=1", "-t[]=2", "-t2"}, "=")

	t1 := conf.Slice("test", "t", append(make([]interface{}, 1), "3"), "Basic description")

	assert.Nil(t, conf.Parse())

	assert.Len(t, *t1, 2)
}

func TestConf_ToStruct(t *testing.T) {

	var testStruct struct {
		TestSlice     []interface{} `comfyname:"test"`
		TestExistence bool          `comfyname:"t2"`
		TestPrtString *string       `comfyname:"test3"`
		TestStruct    struct {
			Test string `comfyname:"test4"`
		}
	}

	conf := prepareConf([]string{"-t[]=1", "-t[]=2", "-t2", "--test3=Pepsioner"}, "=")

	conf.Slice("test", "t", append(make([]interface{}, 1), "3"), "Basic description")
	conf.Exist("test2", "t2", "Basic description 2")
	conf.String("test3", "t3", "c0deum", "Basic description 3")
	conf.String("test4", "t4", "JAre", "Basic description 4")
	conf.String("test5", "t5", "drewoko", "Basic description 4")

	assert.Nil(t, conf.Parse())
	conf.ToStruct(&testStruct)

	assert.Len(t, testStruct.TestSlice, 2)
	assert.True(t, testStruct.TestExistence)
	assert.Equal(t, "Pepsioner", *testStruct.TestPrtString)
	assert.Equal(t, "JAre", testStruct.TestStruct.Test)
}

func TestConf_PrintHelp(t *testing.T) {

	conf := prepareConf([]string{"--Ayangar=FireGM", "-0sk0L0k=Kergan"}, "=")
	conf.String("test", "t", "Ayangar", "Basic description")

	conf.PrintHelp(func(options map[OptionKey]*Option) {
		opt := options[OptionKey{"test", "t"}]
		assert.Equal(t, "Basic description", opt.GetDescription())
		assert.Equal(t, "Ayangar", opt.GetDefaultValue())
		assert.Equal(t, stringType, opt.GetOptionType())
	})
}

func TestDefaultHelpPrinter(t *testing.T) {
	conf := prepareConf([]string{"-t[]=1", "-t[]=2", "-t2", "--test3=Pepsioner"}, "=")

	conf.Slice("test", "t", append(make([]interface{}, 1), "3"), "Basic description")
	conf.Exist("test2", "t2", "Basic description 2")
	conf.String("test3", "t3", "c0deum", "Basic description 3")
	conf.String("test4", "t4", "JAre", "Basic description 4")
	conf.String("test5", "t5", "drewoko", "Basic description 4")

	conf.PrintHelp(DefaultHelpPrinter)
}
