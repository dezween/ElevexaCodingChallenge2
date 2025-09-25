package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dezween/ElevexaCodingChallenge2/internal/routes"
	"github.com/stretchr/testify/assert"
)

func TestServerRoutes(t *testing.T) {
	router := NewRouter()

	// Test that all main endpoints are registered and respond (even if with error)
	cases := []struct {
		method     string
		url        string
		body       string
		wantStatus int
	}{
		{"POST", routes.RouteCreateKey, "", http.StatusCreated},
		{"POST", routes.RouteCreateKey, "", http.StatusConflict}, // duplicate
		{"POST", routes.RouteEncrypt, `{"plaintext":"abc"}`, http.StatusOK},
		{"POST", "/transit/encrypt/unknown", `{"plaintext":"abc"}`, http.StatusNotFound},
		{"POST", routes.RouteDecrypt, `{"ciphertext":"bad","encdata":"bad"}`, http.StatusInternalServerError},
		{"POST", "/transit/decrypt/unknown", `{"ciphertext":"bad","encdata":"bad"}`, http.StatusNotFound},
	}

	for _, tc := range cases {
		url := strings.ReplaceAll(tc.url, "{name}", "testserver")
		req := httptest.NewRequest(tc.method, url, strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, tc.wantStatus, w.Code, "route %s %s", tc.method, url)
	}
}

func TestServerNotFound(t *testing.T) {
	router := NewRouter()
	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
