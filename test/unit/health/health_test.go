package health_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/PlatformStackPulse/go-lambda-template/pkg/health"
)

func TestHealth(t *testing.T) {
	status := health.Health(time.Date(2026, time.April, 13, 9, 0, 0, 0, time.UTC))
	assert.Equal(t, "ok", status.Status)
	assert.Equal(t, "liveness", status.Check)
	assert.Equal(t, "2026-04-13T09:00:00Z", status.Timestamp)
}

func TestReady(t *testing.T) {
	status := health.Ready(time.Date(2026, time.April, 13, 9, 0, 0, 0, time.UTC))
	assert.Equal(t, "ok", status.Status)
	assert.Equal(t, "readiness", status.Check)
	assert.Equal(t, "2026-04-13T09:00:00Z", status.Timestamp)
}
