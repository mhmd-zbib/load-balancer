package services

import (
	"fmt"
	"net/http"
	"time"
)

// PingServices pings each instance of each service and prints their status.
func PingServices(services map[string][]string) {
	for service, instances := range services {
		for _, addr := range instances {
			go func(srv, inst string) {
				client := http.Client{Timeout: 2 * time.Second}
				resp, err := client.Get("http://" + inst + "/health")
				if err != nil {
					fmt.Printf("[HEALTH] %s (%s) is DOWN: %v\n", srv, inst, err)
					return
				}
				defer resp.Body.Close()
				fmt.Printf("[HEALTH] %s (%s) is UP (status: %d)\n", srv, inst, resp.StatusCode)
			}(service, addr)
		}
	}
}
