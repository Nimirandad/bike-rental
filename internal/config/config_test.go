package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_WithDefaultValues(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("SQLITE_PATH")
	os.Unsetenv("HTTP_PORT")

	config := Load()

	assert.Equal(t, "data/bike_rental.db", config.SQLitePath)
	assert.Equal(t, "8080", config.Port)
}

func TestLoad_WithCustomValues(t *testing.T) {
	os.Setenv("SQLITE_PATH", "/custom/path/db.sqlite")
	os.Setenv("HTTP_PORT", "9090")
	defer func() {
		os.Unsetenv("SQLITE_PATH")
		os.Unsetenv("HTTP_PORT")
	}()

	config := Load()

	assert.Equal(t, "/custom/path/db.sqlite", config.SQLitePath)
	assert.Equal(t, "9090", config.Port)
}

func TestGetEnvDefault_WithValue(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnvDefault("TEST_VAR", "default")

	assert.Equal(t, "test_value", result)
}

func TestGetEnvDefault_WithoutValue(t *testing.T) {
	os.Unsetenv("NON_EXISTENT_VAR")

	result := getEnvDefault("NON_EXISTENT_VAR", "default_value")

	assert.Equal(t, "default_value", result)
}