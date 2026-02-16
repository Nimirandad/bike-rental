package types

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	err := WriteJSON(w, 200, data)

	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "key")
	assert.Contains(t, w.Body.String(), "value")
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	WriteError(w, 400, "test error")

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "test error")
}

func TestWriteValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	errors := map[string]string{
		"email":    "invalid email",
		"password": "too short",
	}

	WriteValidationErrors(w, errors)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Validation failed")
	assert.Contains(t, w.Body.String(), "invalid email")
}

func TestWriteSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]int{"count": 5}

	WriteSuccess(w, "success message", data)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success message")
	assert.Contains(t, w.Body.String(), "count")
}

func TestWritePaginatedSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := []string{"item1", "item2"}

	WritePaginatedSuccess(w, "test", data, 100, 2, 10)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "test")
	assert.Contains(t, w.Body.String(), "item1")
}

func TestWritePaginatedSuccess_ZeroTotal(t *testing.T) {
	w := httptest.NewRecorder()
	data := []string{}

	WritePaginatedSuccess(w, "empty", data, 0, 1, 10)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "empty")
}