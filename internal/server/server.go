package server

// Server defines the interface for our TCP server.
type Server interface {
	Start() error
}
