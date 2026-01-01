package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: dummy_server <port>")
	}
	port := os.Args[1]

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","port":%q}`, port)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello from dummy backend on port %s\n", port)
	})

	log.Printf("Dummy backend server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
