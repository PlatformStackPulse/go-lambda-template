package cli_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	appcli "github.com/PlatformStackPulse/go-template/internal/cli"
	"github.com/PlatformStackPulse/go-template/internal/logger"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

func TestRootCommandMetadata(t *testing.T) {
	log := logger.NewLogger(false)
	cmd := appcli.NewRootCommand(log)

	assert.Equal(t, "go-template", cmd.Use)
	assert.NotNil(t, cmd.PersistentPreRun)
}

func TestExampleCommandOutput(t *testing.T) {
	log := logger.NewLogger(false)
	cmd := appcli.NewExampleCommand(log)

	out := captureStdout(t, func() {
		cmd.SetArgs([]string{"--name", "Dev"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	assert.True(t, strings.Contains(out, "Hello, Dev!"))
}
