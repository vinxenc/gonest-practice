package core

import "go.uber.org/fx"

// Server provides the HTTP server (Fiber + Huma), registers the routes of every
// feature module supplied as []Module, and manages the server lifecycle.
var Server = fx.Module("core",
	fx.Provide(
		NewFiber,
		NewHumaAPI,
	),
	fx.Invoke(registerRoutes),
	fx.Invoke(startServer),
)
