package employeeModule

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/0xfurai/gonest"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

// TestEmployeeModule_ServesEmployeesRoute verifies the module wiring: importing
// the exported EmployeeModule alongside a (stubbed) *gorm.DB must resolve the
// Repository -> Service -> Controller graph and serve GET /employees backed by
// that database.
func TestEmployeeModule_ServesEmployeesRoute(t *testing.T) {
	db, mock := newMockDB(t)

	// A global module supplies the stubbed *gorm.DB the employee repository
	// depends on, mirroring how gormModule provides the real connection.
	dbModule := gonest.NewModule(gonest.ModuleOptions{
		Providers: []any{gonest.ProvideValue[*gorm.DB](db)},
		Exports:   []any{(*gorm.DB)(nil)},
		Global:    true,
	})
	root := gonest.NewModule(gonest.ModuleOptions{
		Imports: []*gonest.Module{dbModule, EmployeeModule},
	})

	app := gonest.Create(root)
	if err := app.Init(); err != nil {
		t.Fatalf("initializing app: %v", err)
	}

	// The controller must serve /employees, backed by the stubbed database.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "employees"."employee"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectQuery(`SELECT \* FROM "employees"\."employee"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "gender"}))

	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/employees", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("module controller GET /employees = %d, want %d", rec.Code, http.StatusOK)
	}

	// Confirm the wired controller actually queried the injected *gorm.DB, so this
	// test fails if it stops hitting the database.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
