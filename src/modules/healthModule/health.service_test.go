package healthModule

import "testing"

// TestService_Check verifies the service delegates to its repository and returns
// the status it reports.
func TestService_Check(t *testing.T) {
	svc := HealthService(HealthRepository())
	if svc == nil {
		t.Fatal("HealthService returned nil")
	}
	if got := svc.Check(); got != "ok" {
		t.Fatalf("Check() = %q, want %q", got, "ok")
	}
}
