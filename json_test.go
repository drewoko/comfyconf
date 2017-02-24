package comfyconf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson_Init_WithCustomReader(t *testing.T) {

	jp := NewJSONWithCustomReader(func(j *JSON) ([]byte, error) {
		return []byte(`
			{
			   "c0deum":"NSobolew",
			   "Dunkon":{
			      "megweg":true,
			      "Bushwacker":1,
			      "mofa":{
			         "ews":"AG DobeR",
			         "ichursin":[
			            "Villian.zip"
			         ]
			      }
			   }
			}
		`), nil
	})

	err := jp.Init()

	assert.Nil(t, err)
}

func TestDefaultJsonReader(t *testing.T) {

	j := &JSON{
		path: "testdata/testJsonConfiguration.json",
	}

	b, err := DefaultJSONReader(j)
	b2, err2 := ioutil.ReadFile(j.path)

	assert.Equal(t, err2, err)
	assert.Equal(t, b2, b)
}

func TestJson_Init(t *testing.T) {
	jp := NewJSON("/path")

	assert.Equal(t, "/path", jp.path)
	assert.False(t, jp.flags)
}

func TestNewJsonWithFlagFile(t *testing.T) {
	origArgs := os.Args

	os.Args = []string{"--config=/testpath"}
	jp := NewJSONWithCustomFileMiddleware("c", "config", "/path", NewFlags())

	assert.Equal(t, "/testpath", jp.path)

	os.Args = origArgs
}

func TestJson_ParseInt_Short(t *testing.T) {

	jp := NewJSON("testdata/testJsonConfiguration.json")
	err := jp.Init()

	assert.Nil(t, err)

	v, isOk := jp.ParseInt("Bushwacker", "Long.Bushwacker")

	assert.True(t, isOk)
	assert.Equal(t, 1, v)

	jp = NewJSON("testdata/testJsonConfiguration.json")
	err = jp.Init()

	assert.Nil(t, err)

	v, isOk = jp.ParseInt("Random2", "Random")

	assert.False(t, isOk)
	assert.Equal(t, 0, v)
}

func TestJson_ParseString_Short_Success(t *testing.T) {

	jp := NewJSON("testdata/testJsonConfiguration.json")
	err := jp.Init()

	assert.Nil(t, err)

	v, isOk := jp.ParseString("ews", "mofa")

	assert.True(t, isOk)
	assert.Equal(t, "AG DobeR", v)

	jp = NewJSON("testdata/testJsonConfiguration.json")
	err = jp.Init()

	assert.Nil(t, err)

	v, isOk = jp.ParseString("Random2", "Random")

	assert.False(t, isOk)
	assert.Equal(t, "", v)
}

func TestJson_ParseBool_Success(t *testing.T) {

	jp := NewJSON("testdata/testJsonConfiguration.json")
	err := jp.Init()

	assert.Nil(t, err)

	v, isOk := jp.ParseBool("test", "Dunkon.megweg")

	assert.True(t, isOk)
	assert.Equal(t, true, v)

	jp = NewJSON("testdata/testJsonConfiguration.json")
	err = jp.Init()

	assert.Nil(t, err)

	v, isOk = jp.ParseBool("Random2", "random")

	assert.False(t, isOk)
	assert.Equal(t, false, v)
}

func TestJson_ParseExistence(t *testing.T) {

	jp := NewJSON("testdata/testJsonConfiguration.json")
	err := jp.Init()

	assert.Nil(t, err)

	v, isOk := jp.ParseExistence("c0deum", "c0deum")

	assert.True(t, isOk)
	assert.Equal(t, true, v)
}

func TestJson_ParseSlice(t *testing.T) {

	jp := NewJSON("testdata/testJsonConfiguration.json")
	err := jp.Init()

	assert.Nil(t, err)

	v, isOk := jp.ParseSlice("ichursin", "Dunkon.mofa.ichursin")

	assert.True(t, isOk)
	assert.Equal(t, "Villian.zip", v[0])
}
