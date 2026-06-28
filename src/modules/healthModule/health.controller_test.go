package healthModule

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xfurai/gonest"
)

// newTestHandler builds a minimal gonest application exposing only the health
// controller (backed by the real service and repository) and returns its HTTP
// handler for end-to-end request testing.
func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	mod := gonest.NewModule(gonest.ModuleOptions{
		Controllers: []any{HealthController},
		Providers:   []any{HealthRepository, HealthService},
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
	svc := HealthService(HealthRepository())
	c := HealthController(svc)
	if c == nil {
		t.Fatal("HealthController returned nil")
	}
	if c.service != svc {
		t.Fatal("HealthController did not store the provided service")
	}
}

// TestController_RegisterRoutes_HealthEndpoint registers the controller's routes
// via gonest and exercises GET /health end to end, asserting the status code and
// JSON body.
func TestController_RegisterRoutes_HealthEndpoint(t *testing.T) {
	h := newTestHandler(t)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("response status = %q, want %q", body.Status, "ok")
	}
}
