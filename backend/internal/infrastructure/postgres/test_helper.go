package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// GetTestDB returns a test database connection
func GetTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:admin123@localhost:5432/project_visualization_test?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return nil
	}

	return db
}

// CleanupTestDB cleans up test data from the database
func CleanupTestDB(t *testing.T, db *sqlx.DB, tables ...string) {
	t.Helper()
	ctx := context.Background()

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			t.Logf("Warning: failed to clean up table %s: %v", table, err)
		}
	}
}

// ResetSequences resets auto-increment sequences for tables
func ResetSequences(t *testing.T, db *sqlx.DB, tables ...string) {
	t.Helper()
	ctx := context.Background()

	for _, table := range tables {
		query := fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			t.Logf("Warning: failed to reset sequence for %s: %v", table, err)
		}
	}
}
