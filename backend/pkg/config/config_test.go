package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoad_Defaults verifies that Load() returns the correct defaults
// when no environment variables are set.
func TestLoad_Defaults(t *testing.T) {
	// Clear env vars that might be set in CI
	for _, key := range []string{
		"PORT", "GIN_MODE",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"LOG_LEVEL", "LOG_FORMAT",
		"JWT_SECRET", "CORS_ALLOWED_ORIGINS",
	} {
		t.Setenv(key, "")
	}

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.GinMode)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "admin", cfg.Database.User)
	assert.Equal(t, "admin123", cfg.Database.Password)
	assert.Equal(t, "project_visualization", cfg.Database.DBName)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "dev-secret-change-in-production", cfg.Auth.JWTSecret)
	assert.Equal(t, "", cfg.Auth.AllowedOrigins)
}

// TestLoad_EnvOverrides verifies that environment variables correctly
// override the default configuration values.
func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("GIN_MODE", "release")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "appuser")
	t.Setenv("DB_PASSWORD", "s3cr3t")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("LOG_FORMAT", "text")
	t.Setenv("JWT_SECRET", "my-jwt-secret")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "release", cfg.Server.GinMode)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
	assert.Equal(t, "appuser", cfg.Database.User)
	assert.Equal(t, "s3cr3t", cfg.Database.Password)
	assert.Equal(t, "mydb", cfg.Database.DBName)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, "warn", cfg.Log.Level)
	assert.Equal(t, "text", cfg.Log.Format)
	assert.Equal(t, "my-jwt-secret", cfg.Auth.JWTSecret)
	assert.Equal(t, "https://example.com", cfg.Auth.AllowedOrigins)
}

// TestLoad_InvalidDBPort verifies that Load() returns an error when DB_PORT
// is not a valid integer.
func TestLoad_InvalidDBPort(t *testing.T) {
	t.Setenv("DB_PORT", "not-a-number")

	_, err := Load()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PORT")
}

// TestGetDSN verifies that GetDSN() formats the connection string correctly.
func TestGetDSN(t *testing.T) {
	db := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "admin",
		Password: "secret",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expected := fmt.Sprintf(
		"host=localhost port=5432 user=admin password=secret dbname=testdb sslmode=disable",
	)
	assert.Equal(t, expected, db.GetDSN())
}
