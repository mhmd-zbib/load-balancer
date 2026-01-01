package server

import (
	"fmt"
	"net/http"
)

type HTTPServer struct {
	Addr string
}

func (h *HTTPServer) Start() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Load Balancer HTTP Server Running!")
	})
	return http.ListenAndServe(h.Addr, handler)
}
