package main

import (
	"go.uber.org/fx"

	"gonest-practice/src/core"
	"gonest-practice/src/modules/healthModule"
)

// provideModules is the explicit list of feature modules registered with the
// server. fx injects each module's dependencies; to add a feature, add its
// module here (and its fx.Module to fx.New below).
func provideModules(health *healthModule.Controller) []core.Module {
	return []core.Module{
		healthModule.NewModule(health),
	}
}

func main() {
	app := core.Server(
		healthModule.HealthModule,
		fx.Provide(provideModules),
	)
	app.Run()
}
