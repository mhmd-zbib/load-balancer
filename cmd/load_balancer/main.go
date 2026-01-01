package main


import (
	"log"
	"load_balancer/internal/server"
	"load_balancer/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	srv := &server.TCPServer{Addr: cfg.Addr}
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
