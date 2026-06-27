package employeeModule

import (
	"go.uber.org/fx"

	"gonest-practice/src/core"
)

// EmployeeModule wires the employee feature's providers (repository, service,
// controller) for dependency injection. The repository is exposed as the
// EmployeeReader interface so the service depends on the abstraction, and the
// controller is contributed to the "controllers" group via core.AsController so
// its routes register automatically just by including this module.
//
// The repository depends on *gorm.DB, which core provides (core.NewDatabase), so
// including this module is enough — no extra wiring.
var EmployeeModule = fx.Module("EmployeeModule",
	fx.Provide(
		fx.Annotate(EmployeeRepository, fx.As(new(EmployeeReader))),
		EmployeeService,
		core.AsController(EmployeeController),
	),
)
