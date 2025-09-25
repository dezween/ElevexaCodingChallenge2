package config

import (
	"os"
	"testing"
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
	if cfg.Port != ":8080" {
		t.Fatalf("expected default port :8080, got %q", cfg.Port)
	}
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
	if cfg.Port != ":9090" {
		t.Fatalf("expected port :9090 for KYBER_SERVER_PORT=9090, got %q", cfg.Port)
	}
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
	if cfg.Port != ":9091" {
		t.Fatalf("expected port :9091 for KYBER_SERVER_PORT=:9091, got %q", cfg.Port)
	}
}
