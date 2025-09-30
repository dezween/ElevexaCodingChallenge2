package config

import (
	"os"
	"regexp"
)

// Config holds application configuration parameters.
type Config struct {
	Port string // HTTP server port, e.g. ":8080"
}

var portPattern = regexp.MustCompile(`^:[0-9]{2,5}$`)

// LoadConfig loads configuration from environment variables (with defaults).
// Validates port format (":8080", ":9090", etc). Panics on invalid port.
func LoadConfig() *Config {
	port := os.Getenv("KYBER_SERVER_PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}
	if !portPattern.MatchString(port) {
		panic("Invalid port format: must be :PORT, e.g. :8080")
	}
	return &Config{
		Port: port,
	}
}
