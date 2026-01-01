package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

// TCPServer implements the Server interface.
type TCPServer struct {
	Addr string
}

// Start begins listening for TCP connections and handles them.
func (s *TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.Addr, err)
	}
	defer ln.Close()
	log.Printf("Server listening on %s", s.Addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accepted connection from %s", conn.RemoteAddr())
	if _, err := io.Copy(conn, conn); err != nil {
		log.Printf("error echoing data: %v", err)
	}
}
