package core

import (
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

// Route is implemented by any controller that registers Huma routes.
type Route interface {
	RegisterRoutes(api huma.API)
}

var routeType = reflect.TypeOf((*Route)(nil)).Elem()

// AsRoute annotates a controller constructor so its result joins the "routes"
// value group as a Route. Wrap a module's controller constructor with this inside
// fx.Provide and its routes are registered automatically — no central wiring.
//
// The constructor must be a function whose result implements Route; otherwise it
// panics at startup so misuse fails fast with a clear message rather than as
// opaque fx wiring errors.
func AsRoute(constructor any) any {
	t := reflect.TypeOf(constructor)
	if t == nil || t.Kind() != reflect.Func {
		panic(fmt.Sprintf("core.AsRoute: constructor must be a function, got %T", constructor))
	}
	for i := 0; i < t.NumOut(); i++ {
		if t.Out(i).Implements(routeType) {
			return fx.Annotate(
				constructor,
				fx.As(new(Route)),
				fx.ResultTags(`group:"routes"`),
			)
		}
	}
	panic(fmt.Sprintf("core.AsRoute: constructor %s must return a type implementing core.Route", t))
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
