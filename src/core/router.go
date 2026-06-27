package core

import (
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

// Controller is implemented by anything that registers Huma routes.
type Controller interface {
	RegisterRoutes(api huma.API)
}

// AsController annotates a controller constructor so its result joins the
// "controllers" value group as a Controller. A module wraps its controller
// constructor with this inside fx.Provide, so the controller is collected and
// its routes registered automatically just by including the module — the
// composition root needs no central list.
func AsController(constructor any) any {
	return fx.Annotate(
		constructor,
		fx.As(new(Controller)),
		fx.ResultTags(`group:"controllers"`),
	)
}

// registerRoutes mounts every collected controller onto the Huma API.
func registerRoutes(api huma.API, controllers []Controller) {
	for _, c := range controllers {
		c.RegisterRoutes(api)
	}
}
