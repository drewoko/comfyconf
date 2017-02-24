package comfyconf

import (
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

func TestOption_Put_String_Success(t *testing.T) {

	testStr := "Hi, there"

	var testVarStr string

	opt := Option{
		defaultValue: "Test Default Val",
		variable:     &testVarStr,
	}

	opt.Put(testStr)

	assert.Equal(t, "*string", reflect.TypeOf(opt.variable).String())

	testStr2 := true

	var testVarStr2 string

	opt = Option{
		defaultValue: "Test Default Val",
		variable:     &testVarStr2,
	}

	opt.Put(testStr2)

	assert.Equal(t, "*string", reflect.TypeOf(opt.variable).String())
}

func TestOption_Put_Bool(t *testing.T) {

	testBool := true

	var testVarStr bool

	opt := Option{
		defaultValue: "Test Default Val",
		variable:     &testVarStr,
	}

	opt.Put(testBool)

	assert.Equal(t, "*bool", reflect.TypeOf(opt.variable).String())

	testInt := false

	var testVarStr2 int

	opt = Option{
		defaultValue: "Test Default Val",
		variable:     &testVarStr2,
	}

	opt.Put(testInt)

	assert.Equal(t, "*int", reflect.TypeOf(opt.variable).String())
}

func TestOption_Put_Int_Success(t *testing.T) {

	testBool := 1

	var testVarStr int

	opt := Option{
		defaultValue: "Test Default Val",
		variable:     &testVarStr,
	}

	opt.Put(testBool)

	assert.Equal(t, "*int", reflect.TypeOf(opt.variable).String())
}
