package core

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gonest-practice/src/modules/employeeModule"
	"gonest-practice/src/modules/healthModule"
)

// TestNew_LoadsSettings verifies New loads and returns the validated Settings,
// honoring the PORT environment variable.
func TestNew_LoadsSettings(t *testing.T) {
	t.Setenv("PORT", "8080")

	_, settings, err := New(healthModule.HealthModule)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if settings.Port != 8080 {
		t.Fatalf("settings.Port = %d, want 8080", settings.Port)
	}
}

// TestNew_ServesFeatureAndSwaggerRoutes builds the full application via New with
// both feature modules and exercises it end to end: the /health route is served,
// and the Swagger module mounts an OpenAPI document carrying the configured API
// identity and documenting every feature route. Building the app does not require
// a live database (GORM is opened with pinging disabled), so this stays a pure
// composition-root test — hence it asserts /employees is documented rather than
// issuing a request that would query the database.
func TestNew_ServesFeatureAndSwaggerRoutes(t *testing.T) {
	app, _, err := New(healthModule.HealthModule, employeeModule.EmployeeModule)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if err := app.Init(); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}
	h := app.Handler()

	// The feature route is served.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health = %d, want %d", rec.Code, http.StatusOK)
	}

	// The Swagger module serves an OpenAPI document that carries the API identity
	// and documents every feature route, proving controller route metadata flows
	// into the generated spec.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/swagger/json", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /swagger/json = %d, want %d", rec.Code, http.StatusOK)
	}
	var spec struct {
		Info struct {
			Title   string `json:"title"`
			Version string `json:"version"`
		} `json:"info"`
		Paths map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("decoding OpenAPI spec: %v", err)
	}
	if spec.Info.Title != apiTitle || spec.Info.Version != apiVersion {
		t.Fatalf("OpenAPI info = %q %q, want %q %q", spec.Info.Title, spec.Info.Version, apiTitle, apiVersion)
	}
	for _, path := range []string{"/health", "/employees"} {
		if _, ok := spec.Paths[path]; !ok {
			t.Fatalf("OpenAPI paths missing %s: %v", path, spec.Paths)
		}
	}
}
