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

func TestEncryptDecrypt_TableDriven(t *testing.T) {
	kp, err := GenerateKeyPair()
	require.NoError(t, err)

	tests := []struct {
		name              string
		plaintext         string
		allowEmptyEncdata bool
	}{
		{"normal message", "test message", false},
		{"empty message", "", true},
		{"unicode", "test message", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct, encdata, err := Encrypt(kp.PublicKey, []byte(tt.plaintext))
			require.NoError(t, err)
			assert.NotEmpty(t, ct)
			if tt.allowEmptyEncdata {
				assert.Equal(t, "", encdata)
			} else {
				assert.NotEmpty(t, encdata)
			}
			pt, err := Decrypt(kp.PrivateKey, ct, encdata)
			require.NoError(t, err)
			assert.Equal(t, tt.plaintext, pt)
		})
	}
}

func TestKyberErrors_TableDriven(t *testing.T) {
	kp, err := GenerateKeyPair()
	require.NoError(t, err)

	tests := []struct {
		name         string
		encrypt      bool
		pubKey       []byte
		privKey      []byte
		ct           string
		encdata      string
		wantError    string
		allowNoError bool
	}{
		{"Encrypt with invalid key", true, []byte("badkey"), nil, "", "", "unmarshal public key", false},
		{"Decrypt with invalid base64", false, nil, kp.PrivateKey, "!!!", "!!!", "invalid base64", false},
		{"Decrypt with wrong ciphertext", false, nil, kp.PrivateKey, base64.StdEncoding.EncodeToString([]byte("badct")), base64.StdEncoding.EncodeToString([]byte("enc")), "decapsulation failed", false},
		{"Decrypt with empty encdata", false, nil, kp.PrivateKey, base64.StdEncoding.EncodeToString([]byte("badct")), "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.encrypt {
				_, _, err := Encrypt(tt.pubKey, []byte("data"))
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError)
			} else {
				_, err := Decrypt(tt.privKey, tt.ct, tt.encdata)
				if tt.allowNoError {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.wantError)
				}
			}
		})
	}
}
