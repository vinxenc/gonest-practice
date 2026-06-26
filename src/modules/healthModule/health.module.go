package healthModule

import (
	"go.uber.org/fx"

	"gonest-practice/src/core"
)

// HealthModule wires the health feature's providers (repository, service,
// controller) for dependency injection, analogous to a NestJS module.
var HealthModule = fx.Module("HealthModule",
	fx.Provide(
		HealthRepository,
		HealthService,
		HealthController,
	),
)

// Module bundles the health feature's controllers so the composition root can
// register their routes. It implements core.Module.
type Module struct {
	controllers []core.Controller
}

// NewModule builds the health Module from its controllers.
func NewModule(controller *Controller) *Module {
	return &Module{
		controllers: []core.Controller{controller},
	}
}

// Controllers returns the controllers owned by the health module.
func (m *Module) Controllers() []core.Controller {
	return m.controllers
}
