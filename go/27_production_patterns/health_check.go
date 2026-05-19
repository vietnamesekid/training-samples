package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"time"
)

// HealthStatus is the result of a health check
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
	Memory    MemoryInfo        `json:"memory"`
}

type MemoryInfo struct {
	Alloc      string `json:"alloc"`
	TotalAlloc string `json:"total_alloc"`
	NumGC      uint32 `json:"num_gc"`
}

// Checker defines a specific health check
type Checker interface {
	Name() string
	Check() error
}

// DatabaseChecker checks the DB connection
type DatabaseChecker struct {
	dsn string
}

func (d *DatabaseChecker) Name() string { return "database" }
func (d *DatabaseChecker) Check() error {
	// In practice: db.PingContext(ctx)
	if d.dsn == "" {
		return fmt.Errorf("no database configured")
	}
	return nil // assume OK
}

// CacheChecker checks the cache (Redis, etc.)
type CacheChecker struct{}

func (c *CacheChecker) Name() string { return "cache" }
func (c *CacheChecker) Check() error {
	// In practice: redis.Ping()
	return nil
}

// HealthHandler returns a handler for the health check endpoint
func HealthHandler(checkers []Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := HealthStatus{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Checks:    make(map[string]string),
		}

		// Run all checkers
		for _, checker := range checkers {
			if err := checker.Check(); err != nil {
				status.Checks[checker.Name()] = "unhealthy: " + err.Error()
				status.Status = "degraded"
			} else {
				status.Checks[checker.Name()] = "healthy"
			}
		}

		// Memory stats
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		status.Memory = MemoryInfo{
			Alloc:      fmt.Sprintf("%d KB", ms.Alloc/1024),
			TotalAlloc: fmt.Sprintf("%d KB", ms.TotalAlloc/1024),
			NumGC:      ms.NumGC,
		}

		w.Header().Set("Content-Type", "application/json")
		if status.Status != "ok" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(status)
	}
}

// ReadinessHandler — Kubernetes readiness probe
// Returns 200 when the app is ready to receive traffic
func ReadinessHandler(ready *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if *ready {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "ready")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "not ready")
		}
	}
}

// LivenessHandler — Kubernetes liveness probe
// Returns 200 when the app is still running (not deadlocked)
func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "alive")
	}
}

func demoHealthCheck() {
	checkers := []Checker{
		&DatabaseChecker{dsn: "postgres://localhost/mydb"},
		&CacheChecker{},
	}

	handler := HealthHandler(checkers)
	ready := true

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler)
	mux.HandleFunc("GET /ready", ReadinessHandler(&ready))
	mux.HandleFunc("GET /live", LivenessHandler())

	// Test health endpoint
	req := httptest.NewRequest("GET", "/health", nil)
	rw := httptest.NewRecorder()
	handler(rw, req)

	fmt.Printf("  /health status: %d\n", rw.Code)
	fmt.Printf("  Response: %s", rw.Body.String())

	// Test with a failing checker
	failCheckers := []Checker{
		&DatabaseChecker{dsn: ""}, // will fail
		&CacheChecker{},
	}
	failHandler := HealthHandler(failCheckers)
	req2 := httptest.NewRequest("GET", "/health", nil)
	rw2 := httptest.NewRecorder()
	failHandler(rw2, req2)
	fmt.Printf("  /health (with failure) status: %d\n", rw2.Code)

	fmt.Println("\n  Kubernetes probe endpoints:")
	fmt.Println("  GET /health  — detailed health with all checker statuses")
	fmt.Println("  GET /ready   — readiness probe (200=ready, 503=not ready)")
	fmt.Println("  GET /live    — liveness probe (200=alive)")
}
