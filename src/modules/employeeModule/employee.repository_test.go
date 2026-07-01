package employeeModule

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newMockDB returns a GORM handle backed by a sqlmock database so repository and
// module wiring can be tested without a live PostgreSQL.
func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("creating sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("opening gorm over sqlmock: %v", err)
	}

	// List runs the count and the page query concurrently, so they may reach
	// the mock in either order.
	mock.MatchExpectationsInOrder(false)

	return db, mock
}

// TestRepository_List verifies List issues a count and a paged select and maps
// the rows into Employee values with the reported total.
func TestRepository_List(t *testing.T) {
	db, mock := newMockDB(t)
	repo := EmployeeRepository(db)

	birth := time.Date(1953, time.September, 2, 0, 0, 0, 0, time.UTC)
	hire := time.Date(1986, time.June, 26, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "employees"."employee"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(300024)))
	mock.ExpectQuery(`SELECT \* FROM "employees"\."employee"`).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "birth_date", "first_name", "last_name", "gender", "hire_date"},
		).AddRow(int64(10001), birth, "Georgi", "Facello", "M", hire))

	employees, total, err := repo.List(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if total != 300024 {
		t.Fatalf("total = %d, want 300024", total)
	}
	if len(employees) != 1 {
		t.Fatalf("len(employees) = %d, want 1", len(employees))
	}
	if employees[0].ID != 10001 || employees[0].FirstName != "Georgi" {
		t.Fatalf("employee = %+v, want id 10001 / Georgi", employees[0])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// TestRepository_List_CountError verifies a failing count query, alongside a
// succeeding find, is wrapped and returned.
func TestRepository_List_CountError(t *testing.T) {
	db, mock := newMockDB(t)
	repo := EmployeeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "employees"."employee"`)).
		WillReturnError(context.DeadlineExceeded)
	mock.ExpectQuery(`SELECT \* FROM "employees"\."employee"`).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "birth_date", "first_name", "last_name", "gender", "hire_date"},
		))

	if _, _, err := repo.List(context.Background(), 20, 0); err == nil {
		t.Fatal("List() = nil error, want count error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

// TestRepository_List_FindError verifies a failing select, alongside a
// succeeding count, is wrapped and returned.
func TestRepository_List_FindError(t *testing.T) {
	db, mock := newMockDB(t)
	repo := EmployeeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "employees"."employee"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))
	mock.ExpectQuery(`SELECT \* FROM "employees"\."employee"`).
		WillReturnError(context.DeadlineExceeded)

	if _, _, err := repo.List(context.Background(), 20, 0); err == nil {
		t.Fatal("List() = nil error, want find error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
