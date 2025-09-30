package kybertransit

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/cloudflare/circl/kem/kyber/kyber1024"
)

// KeyPair holds a Kyber public and private key in binary form.
// Use GenerateKeyPair to create a new key pair.
type KeyPair struct {
	PublicKey  []byte // Serialized Kyber public key
	PrivateKey []byte // Serialized Kyber private key
}

// GenerateKeyPair generates a new Kyber-1024 key pair using the CIRCL library.
// Returns a KeyPair with serialized keys, or error on failure.
func GenerateKeyPair() (KeyPair, error) {
	scheme := kyber1024.Scheme()
	pk, sk, err := scheme.GenerateKeyPair()
	if err != nil {
		return KeyPair{}, fmt.Errorf("kyber: failed to generate key pair: %w", err)
	}
	pub, err := pk.MarshalBinary()
	if err != nil {
		return KeyPair{}, fmt.Errorf("kyber: failed to marshal public key: %w", err)
	}
	priv, err := sk.MarshalBinary()
	if err != nil {
		return KeyPair{}, fmt.Errorf("kyber: failed to marshal private key: %w", err)
	}
	return KeyPair{PublicKey: pub, PrivateKey: priv}, nil
}

// Encrypt encrypts plaintext using the given Kyber public key.
// Returns base64-encoded ciphertext and encrypted data (encdata).
//
// SECURITY WARNING: For demonstration, plaintext is XORed with the Kyber shared secret.
// This is NOT secure for production! Use authenticated encryption (e.g., KEM-DEM with AEAD) in real systems.
func Encrypt(pubKey []byte, plaintext []byte) (string, string, error) {
	scheme := kyber1024.Scheme()
	pk, err := scheme.UnmarshalBinaryPublicKey(pubKey)
	if err != nil || pk == nil {
		return "", "", fmt.Errorf("kyber: failed to unmarshal public key: %w", err)
	}
	ct, ss, err := scheme.Encapsulate(pk)
	if err != nil {
		return "", "", fmt.Errorf("kyber: encapsulation failed: %w", err)
	}
	if len(ss) == 0 {
		return "", "", errors.New("kyber: shared secret is empty")
	}
	enc := make([]byte, len(plaintext))
	for i := range plaintext {
		enc[i] = plaintext[i] ^ ss[i%len(ss)]
	}
	// Allow empty plaintext (empty encdata is valid)
	return base64.StdEncoding.EncodeToString(ct), base64.StdEncoding.EncodeToString(enc), nil
}

// Decrypt decrypts base64-encoded ciphertext and encdata using the given Kyber private key.
// Returns the original plaintext, or error.
//
// SECURITY WARNING: For demonstration, encdata is XORed with the Kyber shared secret.
// This is NOT secure for production! Use authenticated encryption (e.g., KEM-DEM with AEAD) in real systems.
func Decrypt(privKey []byte, b64ct string, b64enc string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(b64ct)
	if err != nil {
		return "", fmt.Errorf("kyber: invalid base64 ciphertext: %w", err)
	}
	if b64enc == "" {
		// Allow empty encdata (valid for empty plaintext)
		return "", nil
	}
	enc, err := base64.StdEncoding.DecodeString(b64enc)
	if err != nil {
		return "", fmt.Errorf("kyber: invalid base64 encdata: %w", err)
	}
	scheme := kyber1024.Scheme()
	sk, err := scheme.UnmarshalBinaryPrivateKey(privKey)
	if err != nil || sk == nil {
		return "", fmt.Errorf("kyber: failed to unmarshal private key: %w", err)
	}
	ss, err := scheme.Decapsulate(sk, ct)
	if err != nil {
		return "", fmt.Errorf("kyber: decapsulation failed: %w", err)
	}
	if len(ss) == 0 {
		return "", errors.New("kyber: shared secret is empty")
	}
	plaintext := make([]byte, len(enc))
	for i := range enc {
		plaintext[i] = enc[i] ^ ss[i%len(ss)]
	}
	return string(plaintext), nil
}
