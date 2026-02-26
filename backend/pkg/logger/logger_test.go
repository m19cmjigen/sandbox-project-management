package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestNew_JSONFormat verifies that New() creates a logger with JSON output.
func TestNew_JSONFormat(t *testing.T) {
	log, err := New("info", "json")

	require.NoError(t, err)
	assert.NotNil(t, log)
}

// TestNew_TextFormat verifies that New() creates a logger with text (console) output.
func TestNew_TextFormat(t *testing.T) {
	log, err := New("debug", "text")

	require.NoError(t, err)
	assert.NotNil(t, log)
}

// TestNew_InvalidLevel verifies that an invalid log level falls back to Info
// without returning an error.
func TestNew_InvalidLevel(t *testing.T) {
	log, err := New("not-a-valid-level", "json")

	require.NoError(t, err)
	assert.NotNil(t, log)
	// Invalid level is silently treated as InfoLevel
	assert.True(t, log.Core().Enabled(zap.InfoLevel))
}

// TestWithFields verifies that WithFields() returns a new logger instance
// with additional structured fields attached.
func TestWithFields(t *testing.T) {
	base, err := New("info", "json")
	require.NoError(t, err)

	derived := base.WithFields(zap.String("request_id", "abc123"), zap.Int("user_id", 1))

	assert.NotNil(t, derived)
	assert.NotSame(t, base, derived)
}

// TestLogMethods verifies that the Info/Debug/Warn/Error wrapper methods
// are callable without panicking. Fatal is excluded as it calls os.Exit.
func TestLogMethods(t *testing.T) {
	log, err := New("debug", "json")
	require.NoError(t, err)

	// Each call exercises the corresponding wrapper method.
	assert.NotPanics(t, func() { log.Info("info message", zap.String("key", "val")) })
	assert.NotPanics(t, func() { log.Debug("debug message") })
	assert.NotPanics(t, func() { log.Warn("warn message") })
	assert.NotPanics(t, func() { log.Error("error message") })
}
