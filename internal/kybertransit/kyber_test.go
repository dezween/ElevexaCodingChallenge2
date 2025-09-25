package kybertransit

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	require.NoError(t, err)
	assert.NotEmpty(t, kp.PublicKey)
	assert.NotEmpty(t, kp.PrivateKey)
}

func TestEncryptDecrypt(t *testing.T) {
	kp, err := GenerateKeyPair()
	require.NoError(t, err)
	plaintext := "test message"
	ct, encdata, err := Encrypt(kp.PublicKey, []byte(plaintext))
	require.NoError(t, err)
	assert.NotEmpty(t, ct)
	assert.NotEmpty(t, encdata)
	pt, err := Decrypt(kp.PrivateKey, ct, encdata)
	require.NoError(t, err)
	assert.Equal(t, plaintext, pt)
}

func TestEncryptInvalidKey(t *testing.T) {
	_, _, err := Encrypt([]byte("badkey"), []byte("data"))
	assert.Error(t, err)
}

func TestDecryptInvalidBase64(t *testing.T) {
	kp, err := GenerateKeyPair()
	require.NoError(t, err)
	_, err = Decrypt(kp.PrivateKey, "!!!", "!!!")
	assert.Error(t, err)
}

func TestDecryptWrongCiphertext(t *testing.T) {
	kp, err := GenerateKeyPair()
	require.NoError(t, err)
	badct := base64.StdEncoding.EncodeToString([]byte("badct"))
	encdata := base64.StdEncoding.EncodeToString([]byte("enc"))
	_, err = Decrypt(kp.PrivateKey, badct, encdata)
	assert.Error(t, err)
}
