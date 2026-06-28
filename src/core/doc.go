// Package core is the application's composition root. It builds the gonest
// Application from the shared infrastructure modules and the feature modules,
// and hands the caller the validated Settings needed to start the server.
//
// # The model
//
// gonest is a progressive, NestJS-inspired framework that owns the whole HTTP
// stack: its own dependency-injection container, module system, router (a
// net/http trie adapter), request pipeline, and lifecycle hooks. A controller is
// anything that implements gonest.Controller and registers its routes through
// the Router it is handed:
//
//	type Controller interface {
//	    Register(r gonest.Router)
//	}
//
// Each feature module lists its controllers and providers; gonest resolves the
// Repository -> Service -> Controller graph by type and registers every
// controller's routes automatically, so the composition root only has to import
// the modules.
//
// # Bootstrap
//
// New is the bootstrap entry point, analogous to NestJS's
// NestFactory.create(AppModule). It loads configuration once, composes the
// infrastructure and feature modules under a single root module, and returns the
// runnable Application plus the validated Settings:
//
//	app, settings, err := core.New(
//	    healthModule.HealthModule,
//	    employeeModule.EmployeeModule,
//	)
//	app.ListenAndServeWithGracefulShutdown(fmt.Sprintf(":%d", settings.Port))
//
// # Shared infrastructure
//
//   - Config: the validated *config.Settings is loaded once and shared through DI
//     as a global value provider, so every module sees the same instance without
//     reloading the environment.
//   - Database: gormModule provides a single shared *gorm.DB, exported globally so
//     any feature module can inject it just by declaring the dependency.
//   - Docs: the gonest swagger module serves the OpenAPI document and Swagger UI
//     at /swagger, built from the route metadata controllers declare
//     (Summary/Tags/Response).
//
// # Why no central controller list
//
// Registration is driven by what modules contribute, not by a central list:
// importing a feature module in New is enough for its controllers' routes to be
// served. Adding a feature touches only that feature and the New(...) import
// list — never the routing code.
//
// # Lifecycle
//
// The server is started with ListenAndServeWithGracefulShutdown, which installs
// SIGINT/SIGTERM handlers and, on shutdown, runs the framework's shutdown hooks.
// gormModule uses one of those hooks to close the connection pool cleanly.
package core
