package server

import (
	"fmt"
	"load_balancer/internal/services"
	"net/http"
)

type HTTPServer struct {
	Addr string
}

func (h *HTTPServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Load Balancer HTTP Server Running!")
	})

	mux.HandleFunc("/services", services.ServicesListHandler)
	mux.HandleFunc("/services/", services.ServiceHandler)

	return http.ListenAndServe(h.Addr, mux)
}
