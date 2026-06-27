// Package core is the application's composition root. It provides the HTTP
// server (Fiber + Huma), registers every feature module's routes, and manages
// the server lifecycle.
//
// # The model
//
// A controller is anything that registers Huma routes:
//
//	type Controller interface {
//	    RegisterRoutes(api huma.API)
//	}
//
// Each feature module contributes its controller(s) to a "controllers" fx value
// group with AsController, and core collects the whole group and registers each
// one with a plain loop:
//
//	func registerRoutes(api huma.API, controllers []Controller) {
//	    for _, c := range controllers {
//	        c.RegisterRoutes(api)
//	    }
//	}
//
// # Bootstrap
//
// Server is the bootstrap entry point, analogous to NestJS's
// NestFactory.create(AppModule). It takes feature-module options, provides the
// Fiber + Huma server, and wires fx to invoke initServer:
//
//	app := core.Server(
//	    healthModule.HealthModule,
//	)
//	app.Run()
//
// initServer is the single fx.Invoke. It collects every controller in the
// "controllers" group (via the serverParams fx.In struct), triggers
// registerRoutes, and ties the server to the fx lifecycle (startServer).
// Modules are passed as options rather than a pre-built fx.App, because their
// providers must be registered before fx builds the dependency graph.
//
// # Why a value group
//
// Registration is driven by what modules contribute, not by a central list:
//
//   - Self-contained modules: a module registers its routes simply by being
//     included in Server(...). Including healthModule.HealthModule is enough;
//     there is no separate place to also list its controller.
//   - Open/closed: adding a feature touches only that feature and the Server(...)
//     option list — never core, and never a central registration function.
//   - AsController marks intent: the one-line wrapper at the provider makes a
//     controller's membership explicit and greppable, while the collection and
//     ordering are handled generically by fx.
//
// # Why still fx
//
// fx is kept for what it is good at:
//
//   - Dependency injection: fx resolves each module's
//     Repository -> Service -> Controller graph by type, so there is no
//     hand-written intra-module wiring.
//   - Collection: the "controllers" value group gathers controllers from every
//     module without the composition root knowing their concrete types.
//   - Lifecycle: fx hooks start and gracefully stop the Fiber server alongside
//     the rest of the application (see startServer).
//
// # Why the listener is opened synchronously
//
// app.Listen both binds and serves and would block, so it must run in a
// goroutine. But returning from the OnStart hook before the bind happened would
// hide bind failures (e.g. the port is already in use) from fx's startup path.
// startServer therefore binds with net.Listen synchronously and returns that
// error from OnStart — so fx can fail startup and roll back — then serves the
// already-bound listener in the goroutine.
package core
