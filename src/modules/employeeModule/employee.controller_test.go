package employeeModule

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

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

// TestController_ListEmployees registers the controller's routes on a test API
// and exercises GET /employees end to end, asserting the status code and the
// JSON body (employees plus pagination metadata).
func TestController_ListEmployees(t *testing.T) {
	repo := &fakeReader{
		employees: []Employee{{ID: 10001, FirstName: "Georgi", LastName: "Facello", Gender: "M"}},
		total:     300024,
	}
	c := EmployeeController(EmployeeService(repo))

	_, api := humatest.New(t)
	c.RegisterRoutes(api)

	resp := api.Get("/employees")
	if resp.Code != http.StatusOK {
		t.Fatalf("GET /employees status = %d, want %d", resp.Code, http.StatusOK)
	}

	var body struct {
		Employees []EmployeeDTO `json:"employees"`
		Limit     int           `json:"limit"`
		Offset    int           `json:"offset"`
		Total     int64         `json:"total"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	if len(body.Employees) != 1 || body.Employees[0].ID != 10001 {
		t.Fatalf("employees = %+v, want one employee with id 10001", body.Employees)
	}
	if body.Total != 300024 {
		t.Fatalf("total = %d, want 300024", body.Total)
	}
	// Huma applies the schema defaults for the query params.
	if body.Limit != defaultLimit || body.Offset != 0 {
		t.Fatalf("pagination = limit %d / offset %d, want %d / 0", body.Limit, body.Offset, defaultLimit)
	}
	// The default limit must have reached the data layer.
	if repo.gotLimit != defaultLimit {
		t.Fatalf("repo called with limit %d, want default %d", repo.gotLimit, defaultLimit)
	}
}

// TestController_ListEmployees_WithParams verifies explicit limit/offset query
// params are honored.
func TestController_ListEmployees_WithParams(t *testing.T) {
	repo := &fakeReader{}
	c := EmployeeController(EmployeeService(repo))

	_, api := humatest.New(t)
	c.RegisterRoutes(api)

	if resp := api.Get("/employees?limit=5&offset=15"); resp.Code != http.StatusOK {
		t.Fatalf("GET /employees status = %d, want %d", resp.Code, http.StatusOK)
	}
	if repo.gotLimit != 5 || repo.gotOffset != 15 {
		t.Fatalf("repo called with limit=%d offset=%d, want 5/15", repo.gotLimit, repo.gotOffset)
	}
}

// TestController_ListEmployees_ServiceError verifies a data-layer failure
// surfaces as a 500.
func TestController_ListEmployees_ServiceError(t *testing.T) {
	c := EmployeeController(EmployeeService(&fakeReader{err: errors.New("db down")}))

	_, api := humatest.New(t)
	c.RegisterRoutes(api)

	if resp := api.Get("/employees"); resp.Code != http.StatusInternalServerError {
		t.Fatalf("GET /employees status = %d, want %d", resp.Code, http.StatusInternalServerError)
	}
}
