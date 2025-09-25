package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/dezween/ElevexaCodingChallenge2/internal/kybertransit"
	"github.com/gorilla/mux"
)

var (
	keyStore      = make(map[string]kybertransit.KeyPair)
	keyStoreMutex sync.RWMutex
)

// writeJSON writes a JSON response with the given status code.
// Not exported: only for internal use in handlers.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// CreateKeyHandler handles POST /transit/keys/{name}.
// Generates a new Kyber key pair and stores it in memory.
// Returns 201 on success, 409 if key exists, 500 on internal error.
func CreateKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	keyStoreMutex.Lock()
	defer keyStoreMutex.Unlock()
	if _, exists := keyStore[name]; exists {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Key already exists"})
		return
	}
	kp, err := kybertransit.GenerateKeyPair()
	if err != nil {
		log.Printf("[ERROR] failed to generate key pair: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	keyStore[name] = kp
	writeJSON(w, http.StatusCreated, map[string]string{"message": "Key created"})
}

// EncryptHandler handles POST /transit/encrypt/{name}.
// Encrypts plaintext using the Kyber public key.
// Returns 200 and ciphertext+encdata on success, 404 if key not found, 400/500 on error.
func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	keyStoreMutex.RLock()
	key, exists := keyStore[name]
	keyStoreMutex.RUnlock()
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Key not found"})
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	var req struct {
		Plaintext string `json:"plaintext"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}
	ct, encdata, err := kybertransit.Encrypt(key.PublicKey, []byte(req.Plaintext))
	if err != nil {
		log.Printf("[ERROR] encrypt failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"ciphertext": ct,
		"encdata":    encdata,
	})
}

// DecryptHandler handles POST /transit/decrypt/{name}.
// Decrypts ciphertext+encdata using the Kyber private key.
// Returns 200 and plaintext on success, 404 if key not found, 400/500 on error.
func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	keyStoreMutex.RLock()
	key, exists := keyStore[name]
	keyStoreMutex.RUnlock()
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Key not found"})
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	var req struct {
		Ciphertext string `json:"ciphertext"`
		Encdata    string `json:"encdata"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}
	if req.Ciphertext == "" || req.Encdata == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing ciphertext or encdata"})
		return
	}
	plaintext, err := kybertransit.Decrypt(key.PrivateKey, req.Ciphertext, req.Encdata)
	if err != nil {
		log.Printf("[ERROR] decrypt failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"plaintext": plaintext,
	})
}
