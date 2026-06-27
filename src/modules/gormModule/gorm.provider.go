package gormModule

import (
	"context"
	"fmt"
	"time"

	"gonest-practice/src/config"

	"go.uber.org/fx"
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

// NewGorm provides a GORM database handle backed by PostgreSQL and ties its
// connection pool to the fx lifecycle, closing it on shutdown.
//
// The handle is opened with automatic pinging disabled, so constructing it never
// requires a reachable database: the underlying connection pool is established
// lazily on first query. This keeps app startup (and tests that build the full
// graph) independent of a live database, while real queries still surface
// connection errors when they run.
func NewGorm(lc fx.Lifecycle, settings *config.Settings) (*gorm.DB, error) {
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

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return sqlDB.Close()
		},
	})

	return db, nil
}
