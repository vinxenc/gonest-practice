package core

import "github.com/danielgtaylor/huma/v2"

// Controller is implemented by anything that registers Huma routes.
type Controller interface {
	RegisterRoutes(api huma.API)
}

// Module bundles one or more controllers, mirroring a NestJS module that
// declares the controllers it owns.
type Module interface {
	Controllers() []Controller
}

// registerRoutes mounts every controller of every module onto the Huma API.
// The module list is assembled explicitly at the composition root, so route
// registration is a plain loop with no reflection or hidden wiring.
func registerRoutes(api huma.API, modules []Module) {
	for _, m := range modules {
		for _, c := range m.Controllers() {
			c.RegisterRoutes(api)
		}
	}
}
