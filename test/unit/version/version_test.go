package version_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatformStackPulse/go-template/pkg/version"
)

func TestGetReturnsVersionInfo(t *testing.T) {
	origVersion := version.Version
	origCommit := version.Commit
	origBuildTime := version.BuildTime
	origGoVersion := version.GoVersion
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.BuildTime = origBuildTime
		version.GoVersion = origGoVersion
	}()

	version.Version = "v1.2.3"
	version.Commit = "abc123"
	version.BuildTime = "2026-04-06"
	version.GoVersion = "go1.22"

	info := version.Get()
	assert.Equal(t, "v1.2.3", info.Version)
	assert.Equal(t, "abc123", info.Commit)
	assert.Equal(t, "2026-04-06", info.BuildTime)
	assert.Equal(t, "go1.22", info.GoVersion)
}

func TestInfoStringContainsFields(t *testing.T) {
	info := version.Info{Version: "v1.0.0", Commit: "deadbeef", BuildTime: "today", GoVersion: "go1.22"}
	out := info.String()

	assert.True(t, strings.Contains(out, "Version: v1.0.0"))
	assert.True(t, strings.Contains(out, "Commit: deadbeef"))
	assert.True(t, strings.Contains(out, "Build Time: today"))
	assert.True(t, strings.Contains(out, "Go Version: go1.22"))
}
