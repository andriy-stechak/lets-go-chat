package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "_health", nil)
	assert.Nil(t, err, "%v", err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	expected := `{"status":"ok"}`
	assert.Equal(t, expected, rr.Body.String(), "handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
}
