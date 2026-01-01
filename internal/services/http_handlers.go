package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ServicesListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	svcs := ListServices()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svcs)
}

func RegisterServiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/services", ServicesListHandler)
	mux.HandleFunc("/services/", ServiceHandler)
	mux.HandleFunc("/route/", RouteHandler)
}
// RouteHandler forwards the request to a healthy instance of the requested service.
func RouteHandler(w http.ResponseWriter, r *http.Request) {
       serviceName := r.URL.Path[len("/route/"):]
       if serviceName == "" {
	       http.Error(w, "Service name required", http.StatusBadRequest)
	       return
       }
       svc := GetService(serviceName)
       if svc == nil || len(svc.Instances) == 0 {
	       http.Error(w, "Service not found or has no instances", http.StatusNotFound)
	       return
       }
       // Pick the first healthy instance (no selection logic yet)
       var target *Instance
       for _, inst := range svc.Instances {
	       if inst.Status == StatusUp {
		       target = inst
		       break
	       }
       }
       if target == nil {
	       http.Error(w, "No healthy instances available", http.StatusServiceUnavailable)
	       return
       }
       // Forward the request (simple GET proxy for now)
       resp, err := http.Get("http://" + target.Address)
       if err != nil {
	       http.Error(w, "Error forwarding request: "+err.Error(), http.StatusBadGateway)
	       return
       }
       defer resp.Body.Close()
       w.WriteHeader(resp.StatusCode)
       // Copy response body
       body, err := ioutil.ReadAll(resp.Body)
       if err != nil {
	       http.Error(w, "Error reading response", http.StatusInternalServerError)
	       return
       }
       w.Write(body)
}
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/services/"):]
	if name == "" {
		http.Error(w, "Service name required", http.StatusBadRequest)
		return
	}
	var methodHandlers = map[string]func(http.ResponseWriter, *http.Request, string){
		http.MethodGet:    getServiceHandler,
		http.MethodPost:   setServiceHandler,
		http.MethodPut:    setServiceHandler,
		http.MethodDelete: deleteServiceHandler,
	}
	if handler, ok := methodHandlers[r.Method]; ok {
		handler(w, r, name)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func getServiceHandler(w http.ResponseWriter, r *http.Request, name string) {
	svc := GetService(name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svc)
}

func setServiceHandler(w http.ResponseWriter, r *http.Request, name string) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	var instances []string
	if err := json.Unmarshal(body, &instances); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	SetService(name, instances)
	w.WriteHeader(http.StatusNoContent)
}

func deleteServiceHandler(w http.ResponseWriter, r *http.Request, name string) {
	DeleteService(name)
	w.WriteHeader(http.StatusNoContent)
}
