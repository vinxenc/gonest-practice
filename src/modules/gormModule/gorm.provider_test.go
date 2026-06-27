package gormModule

import (
	"testing"

	"gonest-practice/src/config"

	"go.uber.org/fx/fxtest"
)

// TestNewGorm verifies the provider returns a handle without requiring a live
// database (automatic pinging is disabled) and registers a lifecycle hook that
// closes the connection pool on stop.
func TestNewGorm(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	settings := &config.Settings{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "employees",
		DBSSLMode:  "disable",
	}

	db, err := NewGorm(lc, settings)
	if err != nil {
		t.Fatalf("NewGorm returned error: %v", err)
	}
	if db == nil {
		t.Fatal("NewGorm returned nil db")
	}

	// Start then stop runs the registered OnStop hook, closing the (never
	// connected) connection pool cleanly.
	lc.RequireStart().RequireStop()
}
