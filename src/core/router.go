package core

import (
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

// Route is implemented by any controller that registers Huma routes.
type Route interface {
	RegisterRoutes(api huma.API)
}

// AsRoute annotates a controller constructor so its result joins the "routes"
// value group as a Route. Wrap a module's controller constructor with this inside
// fx.Provide and its routes are registered automatically — no central wiring.
func AsRoute(constructor any) any {
	return fx.Annotate(
		constructor,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

// routeParams collects every Route contributed to the "routes" group.
type routeParams struct {
	fx.In

	API    huma.API
	Routes []Route `group:"routes"`
}

// registerRoutes mounts all collected module routes onto the Huma API.
func registerRoutes(p routeParams) {
	for _, route := range p.Routes {
		route.RegisterRoutes(p.API)
	}
}
