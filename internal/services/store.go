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
	Address     string
	Status      InstanceStatus
	ReqCount    int32
	FailCount   int
	PingLatency int64
}

type Service struct {
	Name      string
	Instances []*Instance
}

var ServiceStore = struct {
	m map[string]*Service
	sync.RWMutex
}{m: make(map[string]*Service)}

// SetService replaces all instances for a service.
func SetService(name string, addresses []string) {
	ServiceStore.Lock()
	instances := make([]*Instance, 0, len(addresses))
	for _, addr := range addresses {
		instances = append(instances, &Instance{Address: addr, Status: StatusUp, ReqCount: 0})
	}
	ServiceStore.m[name] = &Service{Name: name, Instances: instances}
	ServiceStore.Unlock()
	PingServiceNow(name)
}

// GetService retrieves the Service struct for a service.
func GetService(name string) *Service {
	ServiceStore.RLock()
	defer ServiceStore.RUnlock()
	return ServiceStore.m[name]
}

// DeleteService removes a service from the store.
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
		copy[k] = v
	}
	return copy
}
