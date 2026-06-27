package healthModule

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

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
// on a test API and exercises the GET /health endpoint end to end, asserting the
// status code and JSON body.
func TestController_RegisterRoutes_HealthEndpoint(t *testing.T) {
	c := HealthController(HealthService(HealthRepository()))

	_, api := humatest.New(t)
	c.RegisterRoutes(api)

	resp := api.Get("/health")
	if resp.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", resp.Code, http.StatusOK)
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decoding response body: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("response status = %q, want %q", body.Status, "ok")
	}
}
