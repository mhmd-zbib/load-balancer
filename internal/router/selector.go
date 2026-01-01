package router

import "load_balancer/internal/services"

func SelectInstance(serviceName string) []*services.Instance {
	svc := getService(serviceName)
	if svc == nil {
		return nil
	}
	healthy := getHealthyInstances(svc)
	if len(healthy) == 0 {
		return nil
	}
	sortInstancesByLoad(healthy)
	return healthy
}

func getService(serviceName string) *services.Service {
	svc := services.GetService(serviceName)
	if svc == nil || len(svc.Instances) == 0 {
		return nil
	}
	return svc
}

func getHealthyInstances(svc *services.Service) []*services.Instance {
	healthy := make([]*services.Instance, 0, len(svc.Instances))
	for _, inst := range svc.Instances {
		if inst.Status == services.StatusUp {
			healthy = append(healthy, inst)
		}
	}
	return healthy
}

func sortInstancesByLoad(instances []*services.Instance) {
	for i := 0; i < len(instances)-1; i++ {
		for j := i + 1; j < len(instances); j++ {
			if instances[j].ReqCount < instances[i].ReqCount {
				instances[i], instances[j] = instances[j], instances[i]
			}
		}
	}
}

// ...existing code...
