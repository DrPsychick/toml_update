package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoEnv(t *testing.T) {
	os.Clearenv()

	got := getEnv("PFX")

	assert.Empty(t, got, "Environment variables must be empty")
}

func TestEmptyEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "")     //nolint:errcheck
	os.Setenv("PFX_Test2", "foo") //nolint:errcheck

	got := getEnv("PFX")
	assert.Len(t, got, 0, "Environment variables must be skipped")

	got = getEnv("FOO")
	assert.Len(t, got, 0, "Environment variables must be skipped")
}

func TestSimpleEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "foo.bar=bar") //nolint:errcheck
	os.Setenv("PFX_Test2", "bar=")       //nolint:errcheck

	got := getEnv("PFX")

	require.Len(t, got, 2, "2 environment variables must be found")
	assert.Equal(t, "bar", got["foo.bar"])
	assert.Equal(t, "", got["bar"])
}

func TestListEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_TEST", "foo=[\"bar\", \"foo\"]") //nolint:errcheck

	got := getEnv("PFX")

	require.Len(t, got, 1, "Environment variables must be skipped")
	assert.Equal(t, "[\"bar\", \"foo\"]", got["foo"])
}

func TestIntEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "foo=\"1\"") //nolint:errcheck

	got := getEnv("PFX")

	require.Len(t, got, 1)
	assert.Equal(t, "1", got["foo"])
}

func TestParseInt(t *testing.T) {
	got := parseValue("1")
	assert.Equal(t, int64(1), got)

	got = parseValue("2.5")
	assert.NotEqual(t, int64(2), got)
}

func TestParseFloat(t *testing.T) {
	got := parseValue("1.5")
	assert.Equal(t, float64(1.5), got)
}

func TestPrseBool(t *testing.T) {
	got := parseValue("true")
	assert.Equal(t, true, got)
}

func TestParseList(t *testing.T) {
	got := parseValue("[\"foo\"]")
	require.Len(t, got, 1)
	assert.Equal(t, []string{"foo"}, got)

	got = parseValue("[\"bar\", \"foo\", \"\"]")
	require.Len(t, got, 3)
	assert.Equal(t, []string{"bar", "foo", ""}, got)

	got = parseValue("[]")
	require.IsType(t, []string{}, got)
	assert.Len(t, got, 0)
}

func TestParseString(t *testing.T) {
	got := parseValue("foo bar")
	assert.Equal(t, "foo bar", got)
}

func TestConfNoEnv(t *testing.T) {
	os.Clearenv()
	config := []byte("\n[section]\n  key = \"value\"\n")
	_, err := updateConfigFromEnv(config, "PFX")

	assert.Error(t, err)
}

func TestConfEmptySection(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "foo.bar=#no config") //nolint:errcheck

	config := []byte("\n[section]\n  key = \"value\"\n")
	got, err := updateConfigFromEnv(config, "PFX")

	assert.NoError(t, err)
	assert.Contains(t, string(got), "[foo]")
	assert.Contains(t, string(got), "# bar = \"no config\"")
}

func TestUpdateConf(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test", "section.key=value") //nolint:errcheck
	config := []byte("\n[section]\n  key = \"value\"\n")

	got, err := updateConfigFromEnv(config, "PFX")

	assert.NoError(t, err)
	assert.Equal(t, string(config), string(got))
}

func TestComplexConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("PFX_Test0", "section.comment=\"#value\"")     //nolint:errcheck
	os.Setenv("PFX_Test1", "section.integer=1")              //nolint:errcheck
	os.Setenv("PFX_Test2", "section.list=[\"foo\",\"bar\"]") //nolint:errcheck
	os.Setenv("PFX_Test3", "section.float=\"1.337\"")        //nolint:errcheck
	os.Setenv("PFX_Test3", "section.bool=\"false\"")         //nolint:errcheck

	config := []byte("\n[section]\n  bool = false\n  # comment = \"value\"\n  float = 1.337\n  integer = 1\n  list = [\"foo\", \"bar\"]\n")

	got, err := updateConfigFromEnv(config, "PFX")

	assert.NoError(t, err)
	assert.Equal(t, string(config), string(got))
}
