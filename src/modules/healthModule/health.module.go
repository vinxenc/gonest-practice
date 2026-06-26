package healthModule

import (
	"go.uber.org/fx"

	"gonest-practice/src/core"
)

// HealthModule groups the health module's providers, analogous to a NestJS module.
// The controller is registered as a Route, so its endpoints are wired up
// automatically just by including this module.
var HealthModule = fx.Module("HealthModule",
	fx.Provide(
		HealthRepository,
		HealthService,
		core.AsRoute(HealthController),
	),
)
