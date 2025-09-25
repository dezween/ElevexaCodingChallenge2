package kybertransit

import (
	"encoding/base64"
	"errors"

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
		return KeyPair{}, errors.New("kyber: failed to generate key pair: " + err.Error())
	}
	pub, err := pk.MarshalBinary()
	if err != nil {
		return KeyPair{}, errors.New("kyber: failed to marshal public key: " + err.Error())
	}
	priv, err := sk.MarshalBinary()
	if err != nil {
		return KeyPair{}, errors.New("kyber: failed to marshal private key: " + err.Error())
	}
	return KeyPair{PublicKey: pub, PrivateKey: priv}, nil
}

// Encrypt encrypts plaintext using the given Kyber public key.
// Returns base64-encoded ciphertext and encrypted data (encdata).
//
// NOTE: For demonstration, plaintext is XORed with the Kyber shared secret.
// This is NOT secure for production! Use authenticated encryption in real systems.
func Encrypt(pubKey []byte, plaintext []byte) (string, string, error) {
	scheme := kyber1024.Scheme()
	pk, err := scheme.UnmarshalBinaryPublicKey(pubKey)
	if err != nil || pk == nil {
		return "", "", errors.New("kyber: failed to unmarshal public key")
	}
	ct, ss, err := scheme.Encapsulate(pk)
	if err != nil {
		return "", "", errors.New("kyber: encapsulation failed: " + err.Error())
	}
	enc := make([]byte, len(plaintext))
	for i := range plaintext {
		enc[i] = plaintext[i] ^ ss[i%len(ss)]
	}
	return base64.StdEncoding.EncodeToString(ct), base64.StdEncoding.EncodeToString(enc), nil
}

// Decrypt decrypts base64-encoded ciphertext and encdata using the given Kyber private key.
// Returns the original plaintext, or error.
//
// NOTE: For demonstration, encdata is XORed with the Kyber shared secret.
// This is NOT secure for production! Use authenticated encryption in real systems.
func Decrypt(privKey []byte, b64ct string, b64enc string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(b64ct)
	if err != nil {
		return "", errors.New("kyber: invalid base64 ciphertext")
	}
	enc, err := base64.StdEncoding.DecodeString(b64enc)
	if err != nil {
		return "", errors.New("kyber: invalid base64 encdata")
	}
	scheme := kyber1024.Scheme()
	sk, err := scheme.UnmarshalBinaryPrivateKey(privKey)
	if err != nil || sk == nil {
		return "", errors.New("kyber: failed to unmarshal private key")
	}
	ss, err := scheme.Decapsulate(sk, ct)
	if err != nil {
		return "", errors.New("kyber: decapsulation failed: " + err.Error())
	}
	plaintext := make([]byte, len(enc))
	for i := range enc {
		plaintext[i] = enc[i] ^ ss[i%len(ss)]
	}
	return string(plaintext), nil
}
