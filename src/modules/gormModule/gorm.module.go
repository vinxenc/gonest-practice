package gormModule

import "go.uber.org/fx"

// GormModule provides a shared *gorm.DB (PostgreSQL) connection to the
// application. Include it once at the composition root and any other module can
// depend on *gorm.DB simply by declaring it as a constructor parameter — fx
// resolves the same shared handle for every consumer.
//
// It is the infrastructure analogue of a NestJS global database module
// (e.g. TypeOrmModule.forRoot()): a single connection pool, reused everywhere.
var GormModule = fx.Module("GormModule",
	fx.Provide(NewGorm),
)
