package healthModule

import "testing"

// TestRepository_Status verifies the data-access layer reports the expected
// health status.
func TestRepository_Status(t *testing.T) {
	repo := HealthRepository()
	if repo == nil {
		t.Fatal("HealthRepository returned nil")
	}
	if got := repo.Status(); got != "ok" {
		t.Fatalf("Status() = %q, want %q", got, "ok")
	}
}
