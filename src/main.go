package main

import (
	"gonest-practice/src/core"
	"gonest-practice/src/modules/employeeModule"
	"gonest-practice/src/modules/gormModule"
	"gonest-practice/src/modules/healthModule"
)

func main() {
	app := core.Server(
		// Shared infrastructure module: provides the *gorm.DB that feature
		// modules (e.g. employeeModule) depend on.
		gormModule.GormModule,

		healthModule.HealthModule,
		employeeModule.EmployeeModule,
	)
	app.Run()
}
