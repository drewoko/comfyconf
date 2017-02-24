package comfyconf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnv(t *testing.T) {
	env := NewEnv()

	assert.NotNil(t, env)
	assert.Equal(t, "ENV_", env.prefix)
	assert.NotNil(t, env.parsed)
	assert.NotNil(t, env.parsedSlice)
}

func TestNewEnvWithPrefix(t *testing.T) {
	env := NewEnvWithPrefix("TEST_")

	assert.NotNil(t, env)
	assert.Equal(t, "TEST_", env.prefix)
	assert.NotNil(t, env.parsed)
	assert.NotNil(t, env.parsedSlice)
}

func TestEnv_Init(t *testing.T) {

	env := NewEnvWithPrefix("TEST_")
	assert.NotNil(t, env)

	_ = os.Setenv("TEST_Test1", "1")
	_ = os.Setenv("TEST_Test2", "String")
	_ = os.Setenv("TEST_Test3", "false")
	_ = os.Setenv("TEST_Test4[0]", "test1")
	_ = os.Setenv("TEST_Test4[1]", "test2")
	_ = os.Setenv("TEST_Test4[]", "test3")

	assert.NoError(t, env.Init())
}
