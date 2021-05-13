package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoEnv(t *testing.T) {
	os.Clearenv()
	envs := getEnv("PFX")
	assert.Empty(t, envs, "Environment variables must be empty")
}

func TestEmptyEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "")
	os.Setenv("PFX_Test2", "foo")
	envs := getEnv("PFX")
	assert.Len(t, envs, 0, "Environment variables must be skipped")

	envs = getEnv("FOO")
	assert.Len(t, envs, 0, "Environment variables must be skipped")
}

func TestSimpleEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "foo.bar=bar")
	os.Setenv("PFX_Test2", "bar=")

	envs := getEnv("PFX")
	assert.Len(t, envs, 2, "Environment variables must be skipped")
	assert.Equal(t, "bar", envs["foo.bar"])
	assert.Equal(t, "", envs["bar"])
}

func TestListEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_TEST", "foo=[\"bar\", \"foo\"]")
	envs := getEnv("PFX")

	assert.Len(t, envs, 1, "Environment variables must be skipped")
	assert.Equal(t, "[\"bar\", \"foo\"]", envs["foo"])
}

func TestIntEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "foo=\"1\"")
	envs := getEnv("PFX")

	assert.Len(t, envs, 1)
	assert.Equal(t, "1", envs["foo"])
}

func TestParseInt(t *testing.T) {
	val := parseValue("1")
	assert.Equal(t, int64(1), val)

	val = parseValue("2.5")
	assert.NotEqual(t, int64(2), val)
}

func TestParseFloat(t *testing.T) {
	val := parseValue("1.5")
	assert.Equal(t, float64(1.5), val)
}

func TestParseList(t *testing.T) {
	val := parseValue("[\"foo\"]")
	assert.Len(t, val, 1)
	assert.Equal(t, []string{"foo"}, val)

	val = parseValue("[\"bar\", \"foo\", \"\"]")
	assert.Len(t, val, 3)
	assert.Equal(t, []string{"bar", "foo", ""}, val)

	val = parseValue("[]")
	assert.IsType(t, []string{}, val)
	assert.Len(t, val, 0)
}

func TestParseString(t *testing.T) {
	val := parseValue("foo bar")
	assert.Equal(t, "foo bar", val)
}

func TestConfNoEnv(t *testing.T) {
	os.Clearenv()
	config := []byte(fmt.Sprintf("\n[section]\n  key = \"value\"\n"))
	_, err := updateConfigFromEnv(config, "PFX")

	assert.NotNil(t, err)
}

func TestUpdateConf(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "section.key=value")
	config := []byte("\n[section]\n  key = \"value\"\n")

	conf, err := updateConfigFromEnv(config, "PFX")

	assert.Nil(t, err)
	assert.Equal(t, string(config), string(conf))
}

func TestComplexConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test0", "section.comment=#value")
	os.Setenv("PFX_Test1", "section.integer=1")
	os.Setenv("PFX_Test2", "section.list=[\"foo\",\"bar\"]")
	os.Setenv("PFX_Test3", "section.float=1.337")

	config := []byte("\n[section]\n  # comment = \"value\"\n  float = 1.337\n  integer = 1\n  list = [\"foo\", \"bar\"]\n")

	conf, err := updateConfigFromEnv(config, "PFX")

	assert.Nil(t, err)
	assert.Equal(t, string(config), string(conf))
}
