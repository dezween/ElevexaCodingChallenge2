package handlers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dezween/ElevexaCodingChallenge2/internal/routes"
	"github.com/dezween/ElevexaCodingChallenge2/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testKey1   = "testkey"
	testKey2   = "testkey2"
	testKey3   = "testkey3"
	edgeKey    = "edgekey"
	unknownKey = "unknown"
)

func TestCreateKeyHandler_Success(t *testing.T) {
	r := server.NewRouter()
	url, _ := r.Get(routes.RouteNameCreateKey).URL("name", testKey1)
	req := httptest.NewRequest("POST", url.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateKeyHandler_Duplicate(t *testing.T) {
	r := server.NewRouter()
	url, _ := r.Get(routes.RouteNameCreateKey).URL("name", testKey1)
	// first create
	req := httptest.NewRequest("POST", url.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// duplicate
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestEncryptHandler_Success(t *testing.T) {
	r := server.NewRouter()
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", testKey2)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Encrypt
	encReq := map[string]string{"plaintext": "hello quantum world"}
	encBody, _ := json.Marshal(encReq)
	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", testKey2)
	req = httptest.NewRequest("POST", encURL.String(), bytes.NewReader(encBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var encResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &encResp)
	require.NoError(t, err)
	assert.NotEmpty(t, encResp["ciphertext"])
	assert.NotEmpty(t, encResp["encdata"])
}

func TestDecryptHandler_Success(t *testing.T) {
	r := server.NewRouter()
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", testKey2)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Encrypt
	encReq := map[string]string{"plaintext": "hello quantum world"}
	encBody, _ := json.Marshal(encReq)
	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", testKey2)
	req = httptest.NewRequest("POST", encURL.String(), bytes.NewReader(encBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var encResp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &encResp)
	ct := encResp["ciphertext"]
	encdata := encResp["encdata"]

	// Decrypt
	decReq := map[string]string{"ciphertext": ct, "encdata": encdata}
	decBody, _ := json.Marshal(decReq)
	decURL, _ := r.Get(routes.RouteNameDecrypt).URL("name", testKey2)
	req = httptest.NewRequest("POST", decURL.String(), bytes.NewReader(decBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var decResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &decResp)
	require.NoError(t, err)
	assert.Equal(t, "hello quantum world", decResp["plaintext"])
}

func TestEncryptHandler_UnknownKey(t *testing.T) {
	r := server.NewRouter()
	encReq := map[string]string{"plaintext": "data"}
	encBody, _ := json.Marshal(encReq)
	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", unknownKey)
	req := httptest.NewRequest("POST", encURL.String(), bytes.NewReader(encBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDecryptHandler_UnknownKey(t *testing.T) {
	r := server.NewRouter()
	decReq := map[string]string{"ciphertext": "abc", "encdata": "abc"}
	decBody, _ := json.Marshal(decReq)
	decURL, _ := r.Get(routes.RouteNameDecrypt).URL("name", unknownKey)
	req := httptest.NewRequest("POST", decURL.String(), bytes.NewReader(decBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDecryptHandler_InvalidCiphertext(t *testing.T) {
	r := server.NewRouter()
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", testKey3)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Encrypt
	encReq := map[string]string{"plaintext": "data"}
	encBody, _ := json.Marshal(encReq)
	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", testKey3)
	req = httptest.NewRequest("POST", encURL.String(), bytes.NewReader(encBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var encResp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &encResp)

	badct := base64.StdEncoding.EncodeToString([]byte("bad"))
	decReq := map[string]string{"ciphertext": badct, "encdata": encResp["encdata"]}
	decBody, _ := json.Marshal(decReq)
	decURL, _ := r.Get(routes.RouteNameDecrypt).URL("name", testKey3)
	req = httptest.NewRequest("POST", decURL.String(), bytes.NewReader(decBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestEncryptHandler_EmptyPlaintext(t *testing.T) {
	r := server.NewRouter()
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", edgeKey)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	encReq := map[string]string{"plaintext": ""}
	encBody, _ := json.Marshal(encReq)
	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", edgeKey)
	req = httptest.NewRequest("POST", encURL.String(), bytes.NewReader(encBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var encResp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &encResp)
	require.NoError(t, err)
	ct := encResp["ciphertext"]
	assert.NotEmpty(t, ct)
}

func TestDecryptHandler_EmptyEncdata(t *testing.T) {
	r := server.NewRouter()
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", edgeKey)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	encReq := map[string]string{"plaintext": ""}
	encBody, _ := json.Marshal(encReq)
	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", edgeKey)
	req = httptest.NewRequest("POST", encURL.String(), bytes.NewReader(encBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var encResp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &encResp)
	ct := encResp["ciphertext"]

	// Decrypt with empty encdata
	decReq := map[string]string{"ciphertext": ct, "encdata": ""}
	decBody, _ := json.Marshal(decReq)
	decURL, _ := r.Get(routes.RouteNameDecrypt).URL("name", edgeKey)
	req = httptest.NewRequest("POST", decURL.String(), bytes.NewReader(decBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestEncryptHandler_InvalidJSON(t *testing.T) {
	r := server.NewRouter()
	createURL, _ := r.Get(routes.RouteNameCreateKey).URL("name", edgeKey)
	req := httptest.NewRequest("POST", createURL.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	encURL, _ := r.Get(routes.RouteNameEncrypt).URL("name", edgeKey)
	req = httptest.NewRequest("POST", encURL.String(), bytes.NewReader([]byte("notjson")))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
