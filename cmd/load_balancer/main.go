package main

import (
	"context"
	"load_balancer/internal/config"
	"load_balancer/internal/server"
	"load_balancer/internal/services"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	var srv server.Server = &server.HTTPServer{Addr: cfg.Addr}
	log.Printf("Starting server on %s", cfg.Addr)
	go services.StartHealthChecker(context.Background())
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
