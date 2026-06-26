// Package core is the application's composition root. It provides the HTTP
// server (Fiber + Huma) and wires every feature module's routes automatically,
// so feature modules stay decoupled from the bootstrap code.
//
// # The problem
//
// In a naive setup, main (or some central function) imports every controller
// and registers each route by hand:
//
//	health.RegisterRoutes(api)
//	users.RegisterRoutes(api)
//	orders.RegisterRoutes(api)  // ...and one more line for every new module
//
// That central list is a magnet for merge conflicts and is easy to forget to
// update. The bootstrap code has to know about every module, so modules are not
// really independent — adding one means editing shared code (it violates the
// open/closed principle: the system is not open for extension without
// modification).
//
// # The solution
//
// core inverts the dependency using an fx value group. Instead of core
// reaching out to each module, every module contributes its controller into a
// shared "routes" group, and core consumes the whole group in one place:
//
//   - Route is the contract a controller satisfies to be auto-registered:
//     a single RegisterRoutes(huma.API) method.
//
//   - AsRoute wraps a controller constructor so its result is provided into the
//     "routes" group as a Route (via fx.As + fx.ResultTags). A module uses it
//     inside its own fx.Provide, so the module — not core — declares membership:
//
//     var HealthModule = fx.Module("HealthModule",
//     fx.Provide(
//     NewRepository,
//     NewService,
//     core.AsRoute(NewController),
//     ),
//     )
//
//   - registerRoutes receives every Route in the group (via the routeParams
//     fx.In struct, tagged group:"routes") and mounts them all onto the Huma
//     API. core.Module runs this once with fx.Invoke.
//
// The result: main only lists modules, and adding a feature requires zero edits
// to core. The wiring direction is module -> core, never core -> module.
//
// # Why fx (and value groups)
//
//   - Decoupling / extensibility: new modules plug in without touching shared
//     code, which removes the central-list merge-conflict hot spot.
//   - Construction by type: fx resolves the Repository -> Service -> Controller
//     graph automatically, so there is no hand-written wiring to keep in sync.
//   - Lifecycle: fx hooks start and gracefully stop the Fiber server alongside
//     the rest of the application (see startServer).
//
// # Why AsRoute validates eagerly
//
// fx.As and fx.ResultTags map positionally to a constructor's first result, so
// only func(...) Route and func(...) (Route, error) can be annotated correctly.
// AsRoute checks the constructor's shape with reflection and panics at startup
// with a clear, domain-specific message on misuse, instead of letting it surface
// later as an opaque fx wiring error.
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
