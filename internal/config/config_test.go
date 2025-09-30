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
				defer func() {
					if err := os.Setenv("KYBER_SERVER_PORT", prev); err != nil {
						t.Logf("failed to restore env var: %v", err)
					}
				}()
			} else {
				defer func() {
					if err := os.Unsetenv("KYBER_SERVER_PORT"); err != nil {
						t.Logf("failed to unset env var: %v", err)
					}
				}()
			}
			if tt.setEnv {
				if err := os.Setenv("KYBER_SERVER_PORT", tt.envValue); err != nil {
					t.Fatalf("failed to set env var: %v", err)
				}
			} else {
				if err := os.Unsetenv("KYBER_SERVER_PORT"); err != nil {
					t.Fatalf("failed to unset env var: %v", err)
				}
			}
			cfg := LoadConfig()
			assert.Equal(t, tt.expect, cfg.Port)
		})
	}
}
