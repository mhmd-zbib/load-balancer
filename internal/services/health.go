package services

import (
	"context"
	"log"
	"net/http"
	"time"
)

func StartHealthChecker(ctx context.Context) {
	log.Println("[HEALTH] Health checker started (interval: 5s)")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[HEALTH] Health checker stopped")
			return
		case <-ticker.C:
			log.Println("[HEALTH] Health check tick: checking all services...")
			pingAllServices()
		}
	}
}

// pingAllServices checks all instances in the store and updates their status.
func pingAllServices() {
	ServiceStore.RLock()
	defer ServiceStore.RUnlock()
	for _, svc := range ServiceStore.m {
		for _, inst := range svc.Instances {
			go pingAndUpdateInstance(svc.Name, inst)
		}
	}
}

// pingAndUpdateInstance pings an instance and updates its status.
// Instance is defined in store.go in the same package.
// If you see errors, ensure both files are in the same package and compiled together.
func pingAndUpdateInstance(serviceName string, inst *Instance) {
	const failureThreshold = 3
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + inst.Address + "/health")
	ServiceStore.Lock()
	defer ServiceStore.Unlock()
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		inst.FailCount++
		if inst.FailCount >= failureThreshold {
			if inst.Status != StatusDown {
				log.Printf("[HEALTH] %s (%s) is DOWN (failures: %d)", serviceName, inst.Address, inst.FailCount)
			}
			inst.Status = StatusDown
		}
		if err != nil {
			log.Printf("[HEALTH] %s (%s) health check error: %v", serviceName, inst.Address, err)
		} else {
			log.Printf("[HEALTH] %s (%s) health check failed: status %d", serviceName, inst.Address, resp.StatusCode)
		}
	} else {
		if inst.Status != StatusUp {
			log.Printf("[HEALTH] %s (%s) is UP", serviceName, inst.Address)
		}
		inst.Status = StatusUp
		inst.FailCount = 0
	}
	if resp != nil {
		resp.Body.Close()
	}
}

func PingServiceNow(serviceName string) {
	ServiceStore.RLock()
	svc := ServiceStore.m[serviceName]
	ServiceStore.RUnlock()
	if svc == nil {
		return
	}
	for _, inst := range svc.Instances {
		pingAndUpdateInstance(serviceName, inst)
	}
}
