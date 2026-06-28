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
	// "sql: database is closed" on use, with no live database required.
	if err := sqlDB.Ping(); err == nil {
		t.Fatal("connection pool still usable after OnApplicationShutdown, want it closed")
	}
}
