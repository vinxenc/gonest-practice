package core

import "go.uber.org/fx"

// Module provides the HTTP server (Fiber + Huma) and automatically registers
// every Route contributed to the "routes" group by feature modules.
var Module = fx.Module("core",
	fx.Provide(
		NewFiber,
		NewHumaAPI,
	),
	fx.Invoke(registerRoutes),
	fx.Invoke(startServer),
)
