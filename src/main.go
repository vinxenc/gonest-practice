package main

import (
	"go.uber.org/fx"

	"gonest-practice/src/core"
	"gonest-practice/src/modules/healthModule"
)

func main() {
	fx.New(
		core.Module,
		healthModule.HealthModule,
	).Run()
}
