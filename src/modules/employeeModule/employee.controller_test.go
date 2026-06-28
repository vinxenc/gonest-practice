package employeeModule

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xfurai/gonest"
)

// newControllerHandler builds a minimal gonest application exposing only the
// employee controller, backed by the given fake reader, and returns its HTTP
// handler for end-to-end request testing.
func newControllerHandler(t *testing.T, repo *fakeReader) http.Handler {
	t.Helper()
	mod := gonest.NewModule(gonest.ModuleOptions{
		Controllers: []any{EmployeeController},
		Providers: []any{
			EmployeeService,
			gonest.Bind[EmployeeReader](func() *fakeReader { return repo }),
		},
	})
	app := gonest.Create(mod)
	if err := app.Init(); err != nil {
		t.Fatalf("initializing app: %v", err)
	}
	return app.Handler()
}

// TestController_Construction verifies the controller constructor wires in the
// provided service.
func TestController_Construction(t *testing.T) {
	svc := EmployeeService(&fakeReader{})
	c := EmployeeController(svc)
	if c == nil {
		t.Fatal("EmployeeController returned nil")
	}
	if c.service != svc {
		t.Fatal("EmployeeController did not store the provided service")
	}
}

// TestController_ListEmployees exercises GET /employees end to end, asserting the
// status code and the JSON body (employees plus pagination metadata).
func TestController_ListEmployees(t *testing.T) {
	repo := &fakeReader{
		employees: []Employee{{ID: 10001, FirstName: "Georgi", LastName: "Facello", Gender: "M"}},
		total:     300024,
	}
	h := newControllerHandler(t, repo)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/employees", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /employees status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body ListEmployeesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	if len(body.Employees) != 1 || body.Employees[0].ID != 10001 {
		t.Fatalf("employees = %+v, want one employee with id 10001", body.Employees)
	}
	if body.Total != 300024 {
		t.Fatalf("total = %d, want 300024", body.Total)
	}
	// With no query params, the service applies its default page size.
	if body.Limit != defaultLimit || body.Offset != 0 {
		t.Fatalf("pagination = limit %d / offset %d, want %d / 0", body.Limit, body.Offset, defaultLimit)
	}
	// The default limit must have reached the data layer.
	if repo.gotLimit != defaultLimit {
		t.Fatalf("repo called with limit %d, want default %d", repo.gotLimit, defaultLimit)
	}
}

// TestController_ListEmployees_WithParams verifies explicit limit/offset query
// params are honored end to end: they reach the data layer and are echoed back
// in the response pagination metadata.
func TestController_ListEmployees_WithParams(t *testing.T) {
	repo := &fakeReader{}
	h := newControllerHandler(t, repo)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/employees?limit=5&offset=15", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /employees status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body ListEmployeesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	if body.Limit != 5 || body.Offset != 15 {
		t.Fatalf("response pagination = limit %d / offset %d, want 5 / 15", body.Limit, body.Offset)
	}
	if repo.gotLimit != 5 || repo.gotOffset != 15 {
		t.Fatalf("repo called with limit=%d offset=%d, want 5/15", repo.gotLimit, repo.gotOffset)
	}
}

// TestController_ListEmployees_ClampsOutOfRange verifies an out-of-range limit is
// clamped by the service (rather than rejected) and the clamped value is what
// reaches the data layer and is echoed in the response.
func TestController_ListEmployees_ClampsOutOfRange(t *testing.T) {
	repo := &fakeReader{}
	h := newControllerHandler(t, repo)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/employees?limit=1000", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /employees status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body ListEmployeesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	if body.Limit != maxLimit {
		t.Fatalf("response limit = %d, want clamped %d", body.Limit, maxLimit)
	}
	if repo.gotLimit != maxLimit {
		t.Fatalf("repo called with limit %d, want clamped %d", repo.gotLimit, maxLimit)
	}
}

// TestController_ListEmployees_ServiceError verifies a data-layer failure
// surfaces as a 500.
func TestController_ListEmployees_ServiceError(t *testing.T) {
	h := newControllerHandler(t, &fakeReader{err: errors.New("db down")})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/employees", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("GET /employees status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
