package gormModule

import (
	"testing"

	"gonest-practice/src/config"
)

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

	// The lifecycle wrapper closes the (never-connected) connection pool cleanly
	// when gonest runs its OnApplicationShutdown hook on shutdown.
	if err := newConnection(db).OnApplicationShutdown(""); err != nil {
		t.Fatalf("OnApplicationShutdown returned error: %v", err)
	}
}
