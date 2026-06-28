package gormModule

import (
	"strings"
	"testing"

	"gonest-practice/src/config"
)

// errPoolClosed is the error database/sql reports for any use of a *sql.DB after
// Close. Asserting it (rather than any error) distinguishes "the pool was closed"
// from "the database is merely unreachable".
const errPoolClosed = "sql: database is closed"

// TestNewGorm verifies the provider returns a handle without requiring a live
// database (automatic pinging is disabled).
func TestNewGorm(t *testing.T) {
	settings := &config.Settings{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "employees",
		DBSSLMode:  "disable",
	}

	db, err := NewGorm(settings)
	if err != nil {
		t.Fatalf("NewGorm returned error: %v", err)
	}
	if db == nil {
		t.Fatal("NewGorm returned nil db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("accessing connection pool: %v", err)
	}

	// The lifecycle wrapper closes the (never-connected) connection pool cleanly
	// when gonest runs its OnApplicationShutdown hook on shutdown.
	if err := newConnection(db).OnApplicationShutdown(""); err != nil {
		t.Fatalf("OnApplicationShutdown returned error: %v", err)
	}

	// And the pool must actually be closed afterwards: a closed *sql.DB reports
	// errPoolClosed on use, with no live database required. Asserting that
	// specific error (not just any failure) ensures the test would catch the hook
	// no longer closing the pool, rather than passing because Postgres is down.
	if err := sqlDB.Ping(); err == nil || !strings.Contains(err.Error(), errPoolClosed) {
		t.Fatalf("Ping after OnApplicationShutdown = %v, want %q", err, errPoolClosed)
	}
}
