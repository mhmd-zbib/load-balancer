package main

import (
	"load_balancer/internal/config"
	"load_balancer/internal/server"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	var srv server.Server = &server.HTTPServer{Addr: cfg.Addr}
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
