package services

import "sync"

// ServiceStore holds the mapping of service names to their backend instances.
var ServiceStore = struct {
	m map[string][]string
	sync.RWMutex
}{m: make(map[string][]string)}

// SetService sets the instances for a service.
func SetService(name string, instances []string) {
	ServiceStore.Lock()
	defer ServiceStore.Unlock()
	ServiceStore.m[name] = instances
}

// GetService gets the instances for a service.
func GetService(name string) []string {
	ServiceStore.RLock()
	defer ServiceStore.RUnlock()
	return ServiceStore.m[name]
}

// DeleteService removes a service.
func DeleteService(name string) {
	ServiceStore.Lock()
	defer ServiceStore.Unlock()
	delete(ServiceStore.m, name)
}

// ListServices returns all services and their instances.
func ListServices() map[string][]string {
	ServiceStore.RLock()
	defer ServiceStore.RUnlock()
	copy := make(map[string][]string)
	for k, v := range ServiceStore.m {
		copy[k] = append([]string{}, v...)
	}
	return copy
}
