package gormModule

import (
	"fmt"
	"time"

	"gonest-practice/src/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connection-pool bounds applied to the underlying *sql.DB. Sensible defaults
// that keep connection usage bounded under load.
const (
	maxOpenConns    = 25
	maxIdleConns    = 5
	connMaxLifetime = time.Hour
	connMaxIdleTime = 30 * time.Minute
)

// NewGorm provides a GORM database handle backed by PostgreSQL.
//
// The handle is opened with automatic pinging disabled, so constructing it never
// requires a reachable database: the underlying connection pool is established
// lazily on first query. This keeps app startup (and tests that build the full
// DI graph) independent of a live database, while real queries still surface
// connection errors when they run. *config.Settings is resolved from the DI
// container, which gonest supplies from the shared configuration.
func NewGorm(settings *config.Settings) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(settings.DatabaseDSN()), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("accessing database connection pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	return db, nil
}

// connection ties the shared *gorm.DB to the application lifecycle. gonest
// resolves it alongside the other providers and invokes its OnApplicationShutdown
// hook during graceful shutdown, where it closes the underlying connection pool.
// Replaces the fx OnStop hook the previous fx-based wiring used.
type connection struct {
	db *gorm.DB
}

// newConnection wraps the shared handle so its pool is closed on shutdown.
func newConnection(db *gorm.DB) *connection {
	return &connection{db: db}
}

// OnApplicationShutdown closes the underlying connection pool. It implements the
// gonest.OnApplicationShutdown lifecycle hook (the signal name is unused).
func (c *connection) OnApplicationShutdown(string) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
