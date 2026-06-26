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

var (
	routeType = reflect.TypeOf((*Route)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

// AsRoute annotates a controller constructor so its result joins the "routes"
// value group as a Route. Wrap a module's controller constructor with this inside
// fx.Provide and its routes are registered automatically — no central wiring.
//
// The constructor must return either Route or (Route, error). fx.As and
// fx.ResultTags map positionally to the first result, so any other shape would
// be annotated incorrectly; AsRoute panics at startup on such misuse so it fails
// fast with a clear message rather than as opaque fx wiring errors.
func AsRoute(constructor any) any {
	t := reflect.TypeOf(constructor)
	if t == nil || t.Kind() != reflect.Func {
		panic(fmt.Sprintf("core.AsRoute: constructor must be a function, got %T", constructor))
	}
	if !validRouteConstructor(t) {
		panic(fmt.Sprintf("core.AsRoute: constructor %s must return Route or (Route, error)", t))
	}
	return fx.Annotate(
		constructor,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

// validRouteConstructor reports whether t is a function returning Route or
// (Route, error): the first result must implement Route and an optional second
// result must be the error interface.
func validRouteConstructor(t reflect.Type) bool {
	if t.NumOut() < 1 || t.NumOut() > 2 {
		return false
	}
	if !t.Out(0).Implements(routeType) {
		return false
	}
	return t.NumOut() == 1 || t.Out(1) == errorType
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
