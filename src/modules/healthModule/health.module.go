package healthModule

import "github.com/0xfurai/gonest"

// HealthModule wires the health feature's providers (repository, service) and its
// controller. gonest resolves the Repository -> Service -> Controller graph by
// type and registers the controller's routes automatically, so the feature is
// enabled simply by importing this module at the composition root — no central
// route list.
var HealthModule = gonest.NewModule(gonest.ModuleOptions{
	Controllers: []any{HealthController},
	Providers: []any{
		HealthRepository,
		HealthService,
	},
})
