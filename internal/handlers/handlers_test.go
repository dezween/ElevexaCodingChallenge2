package handlers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dezween/ElevexaCodingChallenge2/internal/handlers"
	"github.com/dezween/ElevexaCodingChallenge2/internal/routes"
	"github.com/dezween/ElevexaCodingChallenge2/internal/server"
	"github.com/stretchr/testify/assert"
)

const (
	testKey1   = "testkey"
	testKey2   = "testkey2"
	testKey3   = "testkey3"
	edgeKey    = "edgekey"
	unknownKey = "unknown"
)

func TestCreateKeyHandler_TableDriven(t *testing.T) {
	// Ensure test isolation
	handlers.ResetKeyStore()
	r := server.NewRouter()
	tests := []struct {
		name       string
		keyName    string
		firstReq   bool
		wantStatus int
		wantField  string
		wantValue  string
	}{
		{"create new key", testKey1, true, http.StatusCreated, "message", "Key created"},
		{"duplicate key", testKey1, false, http.StatusConflict, "error", "Key already exists"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, _ := r.Get(routes.RouteNameCreateKey).URL("name", tt.keyName)
			if tt.firstReq {
				req := httptest.NewRequest("POST", url.String(), nil)
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				assert.Equal(t, tt.wantStatus, w.Code)
				var resp map[string]string
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Equal(t, tt.wantValue, resp[tt.wantField])
			} else {
				// duplicate
				req := httptest.NewRequest("POST", url.String(), nil)
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				assert.Equal(t, tt.wantStatus, w.Code)
				var resp map[string]string
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Equal(t, tt.wantValue, resp[tt.wantField])
			}
		})
	}
}

func TestEncryptHandler_TableDriven(t *testing.T) {
	handlers.ResetKeyStore()
	r := server.NewRouter()
	// First, we create a key
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", testKey2)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	tests := []struct {
		name       string
		keyName    string
		body       interface{}
		wantStatus int
		wantField  string
		wantValue  string
	}{
		{"success", testKey2, map[string]string{"plaintext": "hello quantum world"}, http.StatusOK, "ciphertext", ""},
		{"unknown key", unknownKey, map[string]string{"plaintext": "data"}, http.StatusNotFound, "error", "Key not found"},
		{"invalid JSON", testKey2, "notjson", http.StatusBadRequest, "error", "Invalid JSON"},
		{"missing plaintext", testKey2, map[string]string{}, http.StatusBadRequest, "error", "Missing plaintext"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", tt.keyName)
			var bodyBytes []byte
			if s, ok := tt.body.(string); ok {
				bodyBytes = []byte(s)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}
			req := httptest.NewRequest("POST", encURL.String(), bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			var resp map[string]string
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			if tt.wantValue != "" {
				assert.Equal(t, tt.wantValue, resp[tt.wantField])
			} else {
				assert.NotEmpty(t, resp[tt.wantField])
			}
		})
	}
}

func TestDecryptHandler_TableDriven(t *testing.T) {
	handlers.ResetKeyStore()
	r := server.NewRouter()
	// First, we create a key and encrypt the data.
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", testKey3)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	encReq := map[string]string{"plaintext": "data"}
	encBody, _ := json.Marshal(encReq)
	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", testKey3)
	req = httptest.NewRequest("POST", encURL.String(), bytes.NewReader(encBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var encResp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &encResp)
	ct := encResp["ciphertext"]
	encdata := encResp["encdata"]

	tests := []struct {
		name       string
		keyName    string
		body       map[string]string
		wantStatus int
		wantField  string
		wantValue  string
	}{
		{"success", testKey3, map[string]string{"ciphertext": ct, "encdata": encdata}, http.StatusOK, "plaintext", "data"},
		{"unknown key", unknownKey, map[string]string{"ciphertext": ct, "encdata": encdata}, http.StatusNotFound, "error", "Key not found"},
		{"invalid ciphertext", testKey3, map[string]string{"ciphertext": base64.StdEncoding.EncodeToString([]byte("bad")), "encdata": encdata}, http.StatusBadRequest, "error", "Decryption failed: invalid ciphertext, encdata, or internal error"},
		{"missing ciphertext", testKey3, map[string]string{"encdata": encdata}, http.StatusBadRequest, "error", "Missing ciphertext or encdata"},
		{"missing encdata", testKey3, map[string]string{"ciphertext": ct}, http.StatusBadRequest, "error", "Missing ciphertext or encdata"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decURL, _ := r.Get(routes.RouteNameDecrypt).URL("name", tt.keyName)
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", decURL.String(), bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			var resp map[string]string
			_ = json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tt.wantValue, resp[tt.wantField])
		})
	}
}
