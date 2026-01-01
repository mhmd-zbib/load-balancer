package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Addr     string
	Services map[string][]string // service name -> list of instance addresses
}

func LoadConfig() *Config {
	return &Config{
		Addr:     loadAddr(),
		Services: loadServices(),
	}
}

func loadAddr() string {
	addr := os.Getenv("LB_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return addr
}

func loadServices() map[string][]string {
	servicesEnv := os.Getenv("LB_SERVICES")
	services := make(map[string][]string)
	if servicesEnv != "" {
		// Expecting JSON: {"service1": ["host1:port", "host2:port"], "service2": ["host3:port"]}
		_ = json.Unmarshal([]byte(servicesEnv), &services)
	}
	return services
}
