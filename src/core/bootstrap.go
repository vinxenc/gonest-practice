package core

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// Server builds the application and returns a runnable fx app. It provides the
// HTTP server (Fiber + Huma), applies the given feature-module options, and
// initializes the server via initServer. Call Run on the result to start it:
//
//	app := core.Server(
//	    healthModule.HealthModule,
//	    fx.Provide(provideModules),
//	)
//	app.Run()
//
// This is the composition entry point, analogous to NestJS's
// NestFactory.create(AppModule). Feature modules are passed as options (not a
// pre-built fx.App), because their providers must be registered before fx builds
// the dependency graph.
func Server(opts ...fx.Option) *fx.App {
	return fx.New(
		fx.Provide(
			NewFiber,
			NewHumaAPI,
		),
		fx.Options(opts...),
		fx.Invoke(initServer),
	)
}

// initServer registers every module's routes onto the Huma API and ties the
// Fiber server to the fx lifecycle. fx invokes it once during startup.
func initServer(lc fx.Lifecycle, app *fiber.App, api huma.API, modules []Module) {
	registerRoutes(api, modules)
	startServer(lc, app)
}
