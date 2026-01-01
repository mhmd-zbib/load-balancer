package services

import "sync"

type InstanceStatus int

const (
	StatusUnknown InstanceStatus = iota
	StatusUp
	StatusDown
)

func (s InstanceStatus) String() string {
	switch s {
	case StatusUp:
		return "up"
	case StatusDown:
		return "down"
	default:
		return "unknown"
	}
}

type Instance struct {
	Address  string
	Status   InstanceStatus
	ReqCount int
}

type Service struct {
	Name      string
	Instances []*Instance
}

// ServiceStore holds the mapping of service names to their Service struct.
var ServiceStore = struct {
	m map[string]*Service
	sync.RWMutex
}{m: make(map[string]*Service)}

// SetService sets the instances for a service (replaces all instances).
func SetService(name string, addresses []string) {
	ServiceStore.Lock()
	defer ServiceStore.Unlock()
	instances := make([]*Instance, 0, len(addresses))
	for _, addr := range addresses {
		instances = append(instances, &Instance{Address: addr, Status: StatusUp, ReqCount: 0})
	}
	ServiceStore.m[name] = &Service{Name: name, Instances: instances}
}

// GetService gets the Service struct for a service.
func GetService(name string) *Service {
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

// ListServices returns all services and their details.
func ListServices() map[string]*Service {
	ServiceStore.RLock()
	defer ServiceStore.RUnlock()
	copy := make(map[string]*Service)
	for k, v := range ServiceStore.m {
		// Deep copy not strictly needed for read-only listing
		copy[k] = v
	}
	return copy
}
