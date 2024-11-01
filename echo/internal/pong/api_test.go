package pong

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPing(t *testing.T) {
	// Create a new Echo instance
	e := PongEchoServer()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	// Serve the HTTP request
	e.ServeHTTP(rec, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"ping":"pong"}`, rec.Body.String())
}
