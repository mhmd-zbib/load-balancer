package config

import (
	"os"
)

type Config struct {
	Addr string
}

// LoadConfig loads the configuration for the load balancer.
func LoadConfig() *Config {
	return &Config{
		Addr: loadAddr(),
	}
}

func loadAddr() string {
	addr := os.Getenv("LB_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return addr
}
