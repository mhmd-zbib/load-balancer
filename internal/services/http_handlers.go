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
