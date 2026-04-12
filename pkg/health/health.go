// Package health provides health check functionality for API servers.
//
// This package is optional and recommended only for API/HTTP server deployments.
// It provides a framework for registering and executing health checks that can be
// exposed via HTTP endpoints (e.g., GET /health, GET /ready for Kubernetes probes).
//
// For CLI applications, this package is not needed.
//
// Example usage in an API server:
//
//	handler := health.NewHandler()
//	handler.Register(&DatabaseCheck{})
//	handler.Register(&CacheCheck{})
//
//	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
//	    if handler.IsHealthy(r.Context()) {
//	        w.WriteHeader(http.StatusOK)
//	    } else {
//	        w.WriteHeader(http.StatusServiceUnavailable)
//	    }
//	})
package health

import (
	"context"
)

// Status represents health status.
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// Check performs a health check.
type Check interface {
	Name() string
	Check(ctx context.Context) error
}

// Handler manages health checks.
type Handler struct {
	checks map[string]Check
}

// NewHandler creates a new health handler.
func NewHandler() *Handler {
	return &Handler{
		checks: make(map[string]Check),
	}
}

// Register registers a health check.
func (h *Handler) Register(check Check) {
	h.checks[check.Name()] = check
}

// CheckAll performs all health checks.
func (h *Handler) CheckAll(ctx context.Context) map[string]error {
	results := make(map[string]error)
	for name, check := range h.checks {
		results[name] = check.Check(ctx)
	}
	return results
}

// IsHealthy checks if all checks pass.
func (h *Handler) IsHealthy(ctx context.Context) bool {
	for _, err := range h.CheckAll(ctx) {
		if err != nil {
			return false
		}
	}
	return true
}
