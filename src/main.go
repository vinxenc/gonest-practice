package main

import (
	"fmt"
	"log"

	"gonest-practice/src/core"
	"gonest-practice/src/modules/employeeModule"
	"gonest-practice/src/modules/healthModule"
)

func main() {
	// core.New composes the application: it loads configuration, wires the shared
	// infrastructure (config + the GORM connection feature modules depend on) and
	// Swagger docs, and imports the feature modules below.
	app, settings, err := core.New(
		healthModule.HealthModule,
		employeeModule.EmployeeModule,
	)
	if err != nil {
		log.Fatal(err)
	}

	// ListenAndServeWithGracefulShutdown compiles the module tree, registers all
	// routes, starts the HTTP server, and shuts it down cleanly on SIGINT/SIGTERM
	// (closing the database pool via the framework's shutdown hooks).
	addr := fmt.Sprintf(":%d", settings.Port)
	if err := app.ListenAndServeWithGracefulShutdown(addr); err != nil {
		log.Fatal(err)
	}
}
