// Package core is the application's composition root. It provides the HTTP
// server (Fiber + Huma), registers every feature module's routes, and manages
// the server lifecycle.
//
// # The model
//
// A feature is expressed as two small contracts:
//
//   - Controller is anything that registers Huma routes: a single
//     RegisterRoutes(huma.API) method.
//   - Module bundles the controllers a feature owns: Controllers() []Controller.
//     This mirrors a NestJS module that declares its controllers.
//
// The set of modules is assembled explicitly at the composition root (main) as
// a []Module, and core registers them with a plain nested loop — for each
// module, for each controller, call RegisterRoutes:
//
//	func registerRoutes(api huma.API, modules []Module) {
//	    for _, m := range modules {
//	        for _, c := range m.Controllers() {
//	            c.RegisterRoutes(api)
//	        }
//	    }
//	}
//
// # Why an explicit module list
//
// Registration is driven by a list the composition root owns, rather than by
// auto-discovery:
//
//   - Discoverability: every registered module is visible in one place
//     (provideModules in main), and the registration order is explicit and
//     greppable — no reflection, tags, or hidden value groups to reason about.
//   - Simplicity: registerRoutes is an ordinary loop that is trivial to read and
//     unit-test; misuse is a compile error, not a runtime surprise.
//   - Modules own their controllers: a feature can expose several controllers
//     from one Module, keeping related routes together.
//
// The trade-off is that adding a feature means editing the central list (adding
// the module to provideModules), which is an accepted cost for the explicitness.
//
// # Why still fx
//
// fx is kept for what it is good at:
//
//   - Dependency injection: fx resolves each module's
//     Repository -> Service -> Controller graph by type, so there is no
//     hand-written intra-module wiring. provideModules only declares which
//     modules exist; their dependencies are injected.
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
