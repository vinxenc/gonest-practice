package employeeModule

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielgtaylor/huma/v2/humatest"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"gonest-practice/src/core"
)

// TestEmployeeModule_ProvidesControllerToGroup verifies the module wiring:
// building the EmployeeModule (with a stubbed *gorm.DB) must contribute a
// controller into the "controllers" value group whose routes serve GET
// /employees.
func TestEmployeeModule_ProvidesControllerToGroup(t *testing.T) {
	db, mock := newMockDB(t)

	var controllers []core.Controller
	app := fxtest.New(t,
		fx.Supply(db),
		EmployeeModule,
		fx.Invoke(fx.Annotate(
			func(cs []core.Controller) { controllers = cs },
			fx.ParamTags(`group:"controllers"`),
		)),
	)
	defer app.RequireStop()
	app.RequireStart()

	if len(controllers) != 1 {
		t.Fatalf("EmployeeModule contributed %d controllers, want 1", len(controllers))
	}

	// The contributed controller must serve the /employees route, backed by the
	// stubbed database.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "employees"."employee"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectQuery(`SELECT \* FROM "employees"\."employee"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "gender"}))

	_, api := humatest.New(t)
	controllers[0].RegisterRoutes(api)
	if resp := api.Get("/employees"); resp.Code != http.StatusOK {
		t.Fatalf("module controller GET /employees = %d, want %d", resp.Code, http.StatusOK)
	}
}
