package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		setEnv   bool
		envValue string
		expect   string
	}{
		{"default (unset)", false, "", ":8080"},
		{"no colon", true, "9090", ":9090"},
		{"with colon", true, ":9091", ":9091"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev, had := os.LookupEnv("KYBER_SERVER_PORT")
			if had {
				defer os.Setenv("KYBER_SERVER_PORT", prev)
			} else {
				defer os.Unsetenv("KYBER_SERVER_PORT")
			}
			if tt.setEnv {
				os.Setenv("KYBER_SERVER_PORT", tt.envValue)
			} else {
				os.Unsetenv("KYBER_SERVER_PORT")
			}
			cfg := LoadConfig()
			assert.Equal(t, tt.expect, cfg.Port)
		})
	}
}
