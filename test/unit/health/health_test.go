package health_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatformStackPulse/go-template/pkg/health"
)

type mockCheck struct {
	name string
	err  error
}

func (m mockCheck) Name() string { return m.name }
func (m mockCheck) Check(ctx context.Context) error {
	return m.err
}

func TestHandlerHealthy(t *testing.T) {
	h := health.NewHandler()
	h.Register(mockCheck{name: "db", err: nil})
	h.Register(mockCheck{name: "cache", err: nil})

	results := h.CheckAll(context.Background())
	assert.Len(t, results, 2)
	assert.True(t, h.IsHealthy(context.Background()))
}

func TestHandlerUnhealthy(t *testing.T) {
	h := health.NewHandler()
	h.Register(mockCheck{name: "db", err: errors.New("db down")})

	results := h.CheckAll(context.Background())
	assert.Error(t, results["db"])
	assert.False(t, h.IsHealthy(context.Background()))
}
