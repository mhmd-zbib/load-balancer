package config

import (
	"os"
	"strings"
)

type Config struct {
	Addr           string
	BackendServers []string
}

func LoadConfig() *Config {
	return &Config{
		Addr:           loadAddr(),
		BackendServers: loadBackends(),
	}
}

func loadAddr() string {
	addr := os.Getenv("LB_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return addr
}

func loadBackends() []string {
	backendsEnv := os.Getenv("LB_BACKENDS")
	var backends []string
	if backendsEnv != "" {
		for _, b := range strings.Split(backendsEnv, ",") {
			trimmed := strings.TrimSpace(b)
			if trimmed != "" {
				backends = append(backends, trimmed)
			}
		}
	}
	return backends
}
