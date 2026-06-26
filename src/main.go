package main

import (
	"gonest-practice/src/core"
	"gonest-practice/src/modules/healthModule"
)

func main() {
	app := core.Server(
		healthModule.HealthModule,
	)
	app.Run()
}
