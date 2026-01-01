package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ServicesListHandler handles GET /services
func ServicesListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	svcs := ListServices()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svcs)
}

// ServiceHandler handles CRUD for /services/{name}
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/services/"):]
	if name == "" {
		http.Error(w, "Service name required", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		instances := GetService(name)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(instances)
	case http.MethodPost, http.MethodPut:
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
	case http.MethodDelete:
		DeleteService(name)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
