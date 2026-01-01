package router

import (
	"load_balancer/internal/services"
	"net/http"
	"sync/atomic"
)

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	serviceName, err := extractServiceName(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	instances := SelectInstance(serviceName)
	if len(instances) == 0 {
		http.Error(w, "No healthy instances available", http.StatusServiceUnavailable)
		return
	}
	resp, err := tryInstances(instances)
	if err != nil {
		http.Error(w, "All instances failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	writeServiceResponse(w, resp)
}

// tryInstances attempts to forward the request to each instance in order, returning the first successful response.
func tryInstances(instances []*services.Instance) (*http.Response, error) {
	var lastErr error
	for _, inst := range instances {
		incrementReqCount(inst)
		resp, err := forwardRequest(inst)
		decrementReqCount(inst)
		if err != nil {
			lastErr = err
			continue
		}
		return resp, nil
	}
	return nil, lastErr
}

func incrementReqCount(inst *services.Instance) {
	atomic.AddInt32((*int32)(&inst.ReqCount), 1)
}

func decrementReqCount(inst *services.Instance) {
	atomic.AddInt32((*int32)(&inst.ReqCount), -1)
}

func extractServiceName(r *http.Request) (string, error) {
	serviceName := r.URL.Path[len("/route/"):]
	if serviceName == "" {
		return "", http.ErrMissingFile
	}
	return serviceName, nil
}

func forwardRequest(target *services.Instance) (*http.Response, error) {
	return http.Get("http://" + target.Address)
}

func writeServiceResponse(w http.ResponseWriter, resp *http.Response) {
	w.WriteHeader(resp.StatusCode)
	copyHeaders(w, resp)
	copyBody(w, resp)
}

func copyHeaders(w http.ResponseWriter, resp *http.Response) {
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
}

func copyBody(w http.ResponseWriter, resp *http.Response) {
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
}
