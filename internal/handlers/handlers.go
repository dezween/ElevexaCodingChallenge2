package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/dezween/ElevexaCodingChallenge2/internal/kybertransit"
	"github.com/gorilla/mux"
)

// KeyStoreManager manages Kyber key pairs in a thread-safe in-memory store.
type KeyStoreManager struct {
	store map[string]kybertransit.KeyPair
	mu    sync.RWMutex
}

// NewKeyStoreManager creates a new in-memory key store manager.
func NewKeyStoreManager() *KeyStoreManager {
	return &KeyStoreManager{store: make(map[string]kybertransit.KeyPair)}
}

// CreateKey creates a new Kyber key pair with the given name.
func (m *KeyStoreManager) CreateKey(name string) (kybertransit.KeyPair, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.store[name]; exists {
		return kybertransit.KeyPair{}, true, nil
	}
	kp, err := kybertransit.GenerateKeyPair()
	if err != nil {
		return kybertransit.KeyPair{}, false, err
	}
	m.store[name] = kp
	return kp, false, nil
}

// GetKey returns the Kyber key pair by name.
func (m *KeyStoreManager) GetKey(name string) (kybertransit.KeyPair, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	kp, exists := m.store[name]
	return kp, exists
}

// Reset clears all keys (for test isolation).
func (m *KeyStoreManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store = make(map[string]kybertransit.KeyPair)
}

// keyStore is a volatile in-memory key-value store for Kyber key pairs.
// All keys are lost on server restart. Not for production use.
// TODO: Use persistent storage for production.
var keyStoreManager = NewKeyStoreManager()

// writeJSON writes a JSON response with the given status code.
// Logs encoding errors.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[ERROR] failed to encode JSON response: %v", err)
	}
}

// CreateKeyHandler handles POST /transit/keys/{name}.
// Generates a new Kyber key pair and stores it in memory.
// Returns 201 on success, 409 if key exists, 500 on internal error.
func CreateKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	kp, exists, err := keyStoreManager.CreateKey(name)
	if err != nil {
		log.Printf("[ERROR] failed to generate key pair: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	if exists {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Key already exists"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{
		"message":    "Key created",
		"public_key": base64.StdEncoding.EncodeToString(kp.PublicKey),
	})
}

// EncryptHandler handles POST /transit/encrypt/{name}.
// Encrypts plaintext using the Kyber public key.
// Returns 200 and ciphertext+encdata on success, 404 if key not found, 400/500 on error.
func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	key, exists := keyStoreManager.GetKey(name)
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
	if req.Plaintext == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing plaintext"})
		return
	}
	ct, encdata, err := kybertransit.Encrypt(key.PublicKey, []byte(req.Plaintext))
	if err != nil {
		log.Printf("[ERROR] encrypt failed: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Encryption failed: invalid input or internal error"})
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
	key, exists := keyStoreManager.GetKey(name)
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
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Decryption failed: invalid ciphertext, encdata, or internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"plaintext": plaintext,
	})
}

// HealthHandler returns 200 OK for health checks.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("[ERROR] failed to write health check response: %v", err)
	}
}

// ResetKeyStore clears the in-memory store. Intended for tests to ensure isolation.
func ResetKeyStore() {
	keyStoreManager.Reset()
}
