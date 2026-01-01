package ratelimiter

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type endpointData struct {
	TokenCount int
	Timestamp  time.Time
}

type clientData struct {
	Endpoints map[string]*endpointData // key: apiPath
}

var (
	mu         sync.Mutex
	userTokens = make(map[string]*clientData) // key: client IP
)

const maxTokens = 100
const refillRate = 1

func getOrCreateEndpointData(ip, apiPath string, now time.Time) *endpointData {
	cData, exists := userTokens[ip]
	if !exists {
		cData = &clientData{Endpoints: make(map[string]*endpointData)}
		userTokens[ip] = cData
	}
	eData, exists := cData.Endpoints[apiPath]
	if !exists {
		eData = &endpointData{TokenCount: maxTokens, Timestamp: now}
		cData.Endpoints[apiPath] = eData
	}
	return eData
}

func refillTokens(eData *endpointData, now time.Time) {
	elapsed := int(now.Sub(eData.Timestamp).Seconds())
	if elapsed > 0 {
		newTokens := eData.TokenCount + elapsed*refillRate
		if newTokens > maxTokens {
			newTokens = maxTokens
		}
		if newTokens > eData.TokenCount {
			eData.TokenCount = newTokens
			eData.Timestamp = now
		}
	}
}

func checkAndConsumeToken(eData *endpointData) bool {
	if eData.TokenCount < 1 {
		return false
	}
	eData.TokenCount--
	return true
}

// RateLimitMiddleware is a stub middleware for demonstration. Replace with your implementation.
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		apiPath := r.URL.Path
		now := time.Now()

		mu.Lock()
		eData := getOrCreateEndpointData(clientIP, apiPath, now)
		refillTokens(eData, now)
		allowed := checkAndConsumeToken(eData)
		mu.Unlock()

		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for %s on %s", clientIP, apiPath)
			return
		}
		next(w, r)
	}
}
