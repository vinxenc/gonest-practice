package core

import (
	"gonest-practice/src/config"

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
			config.Load,
			NewFiber,
			NewHumaAPI,
		),
		fx.Options(opts...),
		fx.Invoke(initServer),
	)
}

// serverParams collects the dependencies initServer needs, including every
// controller contributed by a module to the "controllers" value group.
type serverParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Shutdowner  fx.Shutdowner
	App         *fiber.App
	API         huma.API
	Settings    *config.Settings
	Controllers []Controller `group:"controllers"`
}

// initServer registers every collected controller's routes onto the Huma API
// and ties the Fiber server to the fx lifecycle. fx invokes it once at startup.
func initServer(p serverParams) {
	registerRoutes(p.API, p.Controllers)
	startServer(p.Lifecycle, p.App, p.Shutdowner, p.Settings)
}
