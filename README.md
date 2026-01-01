# Go Gateway

A lightweight API gateway and load balancer written in Go. Routes traffic intelligently, keeps your services healthy, and doesn't let anyone spam your endpoints.

## What It Does

This is a smart reverse proxy that sits in front of your backend services. It routes requests, distributes load, monitors health, and enforces rate limits—all the stuff you need when you're running multiple service instances.

## Features

### Load Balancing
Routes requests to the least busy server. We check request count first, then ping latency as a tiebreaker. Only healthy instances get traffic.

**How:** Simple bubble sort on instance metrics. Track active requests per instance, pick the one with the lowest count. No fancy algorithms needed.

### Health Checking
Pings all service instances every 5 seconds. If an instance fails 3 times in a row, it's marked down and pulled from rotation. When it recovers, it automatically comes back.

**How:** Background goroutine hits `/health` on every instance. Track consecutive failures per instance. Sub-2-second timeout so we don't wait around.

### Rate Limiting
Token bucket algorithm per client IP and endpoint combination. 100 tokens max, refills at 1 token/second. Prevents abuse without being annoying.

**How:** In-memory map tracking tokens by IP + path. Refill based on elapsed time since last request. Middleware intercepts before routing.

### Dynamic Service Registry
Add, update, or remove backend services on the fly via REST API. No config files, no restarts.

**How:** Thread-safe map with RWMutex. POST to `/services/{name}` with a JSON array of addresses. Changes take effect immediately.

### Request Routing
Matches incoming paths to registered services and forwards to available instances. Increments request counter on the chosen instance, decrements when done.

**How:** Path-based matching. Proxy the request, track it, release when complete. Dead simple.

## Running It

Set `LB_ADDR` env var if you want something other than `:8080`.

```bash
go run cmd/load_balancer
```

There's also a dummy server in `dummy/` and simulation scripts in `scripts/` for testing.

## API

- `GET /services` - List all services
- `GET /services/{name}` - Get specific service details
- `POST /services/{name}` - Register or update service instances
- `DELETE /services/{name}` - Remove a service

Send instance addresses as JSON array: `["localhost:3001", "localhost:3002"]`

## The Stack

Pure Go stdlib for HTTP. No frameworks, no dependencies. Goroutines for concurrent health checks. RWMutex for safe concurrent access to service registry.
