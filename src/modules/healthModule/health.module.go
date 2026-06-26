package healthModule

import (
	"go.uber.org/fx"

	"gonest-practice/src/core"
)

// HealthModule wires the health feature's providers (repository, service,
// controller) for dependency injection. The controller is contributed to the
// "controllers" group via core.AsController, so its routes register
// automatically just by including this module — no central list.
var HealthModule = fx.Module("HealthModule",
	fx.Provide(
		HealthRepository,
		HealthService,
		core.AsController(HealthController),
	),
)
