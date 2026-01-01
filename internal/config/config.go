package config

import "os"

// Config holds server configuration.
type Config struct {
	Addr string
}

// LoadConfig loads configuration from environment variables or defaults.
func LoadConfig() *Config {
	addr := os.Getenv("LB_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return &Config{Addr: addr}
}
