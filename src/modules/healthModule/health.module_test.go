package healthModule

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xfurai/gonest"
)

// TestHealthModule_RegistersHealthRoute verifies the module wiring: compiling the
// exported HealthModule must resolve the Repository -> Service -> Controller graph
// and register the GET /health route, which then serves a 200.
func TestHealthModule_RegistersHealthRoute(t *testing.T) {
	app := gonest.Create(HealthModule)
	if err := app.Init(); err != nil {
		t.Fatalf("initializing HealthModule app: %v", err)
	}

	// The module must have contributed exactly the /health route.
	routes := app.GetRoutes()
	var found bool
	for _, r := range routes {
		if r.Method == http.MethodGet && r.Path == "/health" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("HealthModule did not register GET /health; routes = %v", routes)
	}

	// And that route must be served by the compiled module.
	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("module GET /health = %d, want %d", rec.Code, http.StatusOK)
	}
}
