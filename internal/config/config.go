package config

import (
	"os"
)

// Config holds application configuration parameters.
type Config struct {
	Port string // HTTP server port, e.g. ":8080"
}

// LoadConfig loads configuration from environment variables (with defaults).
func LoadConfig() *Config {
	port := os.Getenv("KYBER_SERVER_PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}
	return &Config{
		Port: port,
	}
}
