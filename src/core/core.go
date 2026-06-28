package core

import (
	"fmt"

	"gonest-practice/src/config"
	"gonest-practice/src/modules/gormModule"

	"github.com/0xfurai/gonest"
	"github.com/0xfurai/gonest/swagger"
)

// apiTitle and apiVersion identify the service in the generated OpenAPI document.
const (
	apiTitle   = "Gonest Practice API"
	apiVersion = "1.0.0"
)

// New composes the application. It loads and validates configuration once, wires
// the shared infrastructure (config + GORM) and Swagger documentation, imports
// the given feature modules, and returns a runnable gonest Application together
// with the validated Settings so the caller can derive the listen address.
//
// This is the composition entry point, analogous to NestJS's
// NestFactory.create(AppModule): feature modules are passed in and composed
// under a single root module. The configuration and database connection are
// shared through gonest's DI container, so every module resolves the same
// validated Settings and the same *gorm.DB with no hand-written wiring.
func New(features ...*gonest.Module) (*gonest.Application, *config.Settings, error) {
	settings, err := config.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("loading configuration: %w", err)
	}

	// Infrastructure first (config, then GORM which depends on it), then the
	// feature modules, then Swagger last. Imports compile in order, and each
	// module's exports become resolvable to the modules imported after it.
	imports := make([]*gonest.Module, 0, len(features)+3)
	imports = append(imports, configModule(settings), gormModule.GormModule)
	imports = append(imports, features...)
	imports = append(imports, swaggerModule())

	root := gonest.NewModule(gonest.ModuleOptions{Imports: imports})

	return gonest.Create(root), settings, nil
}

// configModule exposes the already-validated Settings to every module via DI. It
// is global so any feature module can inject *config.Settings without importing
// this module explicitly — the infrastructure analogue of a NestJS global
// ConfigModule.
func configModule(settings *config.Settings) *gonest.Module {
	return gonest.NewModule(gonest.ModuleOptions{
		Providers: []any{gonest.ProvideValue[*config.Settings](settings)},
		Exports:   []any{(*config.Settings)(nil)},
		Global:    true,
	})
}

// swaggerModule serves the OpenAPI document and Swagger UI at /swagger, built
// from the route metadata each controller declares (Summary/Tags/Response).
func swaggerModule() *gonest.Module {
	return swagger.Module(swagger.Options{
		Title:   apiTitle,
		Version: apiVersion,
	})
}
