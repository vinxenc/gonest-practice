package gormModule

import (
	"testing"

	"gonest-practice/src/config"

	"github.com/0xfurai/gonest"
	"gorm.io/gorm"
)

// TestGormModule_ClosesPoolOnShutdown verifies the full lifecycle wiring: when an
// application that imports GormModule is closed, the connection provider's
// OnApplicationShutdown hook actually runs and closes the shared pool. This
// guards against the hook silently never firing (e.g. if the connection provider
// stopped being instantiated), which would leak the pool on shutdown.
func TestGormModule_ClosesPoolOnShutdown(t *testing.T) {
	// Compose config + GormModule the same way core.New does, so the provider
	// graph matches production.
	cfgModule := gonest.NewModule(gonest.ModuleOptions{
		Providers: []any{gonest.ProvideValue[*config.Settings](&config.Settings{
			DBHost:     "localhost",
			DBPort:     5432,
			DBUser:     "postgres",
			DBPassword: "postgres",
			DBName:     "employees",
			DBSSLMode:  "disable",
		})},
		Exports: []any{(*config.Settings)(nil)},
		Global:  true,
	})
	root := gonest.NewModule(gonest.ModuleOptions{
		Imports: []*gonest.Module{cfgModule, GormModule},
	})

	app := gonest.Create(root)
	if err := app.Init(); err != nil {
		t.Fatalf("initializing app: %v", err)
	}

	db, err := gonest.Resolve[*gorm.DB](app.GetContainer())
	if err != nil {
		t.Fatalf("resolving *gorm.DB: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("accessing connection pool: %v", err)
	}

	// Closing the app runs the shutdown hooks; the connection provider must close
	// the pool, with no live database required.
	if err := app.Close(); err != nil {
		t.Fatalf("app.Close returned error: %v", err)
	}
	if err := sqlDB.Ping(); err == nil {
		t.Fatal("connection pool still usable after app.Close(), want it closed by the shutdown hook")
	}
}
