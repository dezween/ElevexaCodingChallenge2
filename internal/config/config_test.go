package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Default(t *testing.T) {
	prev, had := os.LookupEnv("KYBER_SERVER_PORT")
	if had {
		defer os.Setenv("KYBER_SERVER_PORT", prev)
	} else {
		defer os.Unsetenv("KYBER_SERVER_PORT")
	}
	os.Unsetenv("KYBER_SERVER_PORT")

	cfg := LoadConfig()
	assert.Equal(t, ":8080", cfg.Port, "expected default port :8080")
}

func TestLoadConfig_NoColon(t *testing.T) {
	prev, had := os.LookupEnv("KYBER_SERVER_PORT")
	if had {
		defer os.Setenv("KYBER_SERVER_PORT", prev)
	} else {
		defer os.Unsetenv("KYBER_SERVER_PORT")
	}
	os.Setenv("KYBER_SERVER_PORT", "9090")

	cfg := LoadConfig()
	assert.Equal(t, ":9090", cfg.Port, "expected port :9090 for KYBER_SERVER_PORT=9090")
}

func TestLoadConfig_WithColon(t *testing.T) {
	prev, had := os.LookupEnv("KYBER_SERVER_PORT")
	if had {
		defer os.Setenv("KYBER_SERVER_PORT", prev)
	} else {
		defer os.Unsetenv("KYBER_SERVER_PORT")
	}
	os.Setenv("KYBER_SERVER_PORT", ":9091")

	cfg := LoadConfig()
	assert.Equal(t, ":9091", cfg.Port, "expected port :9091 for KYBER_SERVER_PORT=:9091")
}
