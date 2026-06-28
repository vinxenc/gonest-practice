package gormModule

import (
	"github.com/0xfurai/gonest"
	"gorm.io/gorm"
)

// GormModule provides a shared *gorm.DB (PostgreSQL) connection to the
// application. Import it once at the composition root and any other module can
// depend on *gorm.DB simply by declaring it as a constructor parameter — gonest
// resolves the same shared handle for every consumer.
//
// It is the infrastructure analogue of a NestJS global database module
// (e.g. TypeOrmModule.forRoot()): a single connection pool, reused everywhere.
// The module is global, so the handle is available to every feature module
// without an explicit import. The unexported connection provider binds the pool
// to the application lifecycle, closing it on shutdown.
var GormModule = gonest.NewModule(gonest.ModuleOptions{
	Providers: []any{
		NewGorm,
		newConnection,
	},
	Exports: []any{(*gorm.DB)(nil)},
	Global:  true,
})
