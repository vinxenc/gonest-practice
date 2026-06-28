package employeeModule

import "github.com/0xfurai/gonest"

// EmployeeModule wires the employee feature's providers (repository, service) and
// its controller. The repository is bound to the EmployeeReader interface so the
// service depends on the abstraction (and the data layer can be faked in tests),
// and gonest resolves the Repository -> Service -> Controller graph by type,
// registering the controller's routes automatically.
//
// The repository depends on *gorm.DB, provided globally by gormModule, so simply
// importing both modules at the composition root is enough — gonest resolves the
// shared connection with no extra wiring.
var EmployeeModule = gonest.NewModule(gonest.ModuleOptions{
	Controllers: []any{EmployeeController},
	Providers: []any{
		gonest.Bind[EmployeeReader](EmployeeRepository),
		EmployeeService,
	},
})
