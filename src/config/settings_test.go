package config

import (
	"testing"
)

// TestLoad_Default verifies that with no PORT set the default is applied and the
// resulting Settings is valid.
func TestLoad_Default(t *testing.T) {
	t.Setenv("PORT", "")

	s, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if s.Port != 3000 {
		t.Fatalf("default Port = %d, want 3000", s.Port)
	}
}

// TestLoad_ValidPort verifies an in-range PORT is read from the environment.
func TestLoad_ValidPort(t *testing.T) {
	t.Setenv("PORT", "8080")

	s, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if s.Port != 8080 {
		t.Fatalf("Port = %d, want 8080", s.Port)
	}
}

// TestLoad_OutOfRangePort verifies the semantic validate() check rejects a port
// outside [1, 65535].
func TestLoad_OutOfRangePort(t *testing.T) {
	for _, port := range []string{"0", "70000"} {
		t.Run(port, func(t *testing.T) {
			t.Setenv("PORT", port)

			if _, err := Load(); err == nil {
				t.Fatalf("Load() with PORT=%s = nil error, want validation error", port)
			}
		})
	}
}

// TestLoad_NonNumericPort verifies go-env-validator rejects a non-integer PORT
// before the semantic check runs.
func TestLoad_NonNumericPort(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	if _, err := Load(); err == nil {
		t.Fatal("Load() with non-numeric PORT = nil error, want parse error")
	}
}

// TestSettings_DatabaseDSN verifies the DSN is built from the database settings
// in libpq keyword/value form.
func TestSettings_DatabaseDSN(t *testing.T) {
	s := Settings{
		DBHost:     "db.example.com",
		DBPort:     5433,
		DBUser:     "alice",
		DBPassword: "s3cret",
		DBName:     "employees",
		DBSSLMode:  "require",
	}
	want := "host=db.example.com port=5433 user=alice password=s3cret dbname=employees sslmode=require"
	if got := s.DatabaseDSN(); got != want {
		t.Fatalf("DatabaseDSN() = %q, want %q", got, want)
	}
}

// TestLoad_DatabaseDefaults verifies the database settings fall back to their
// defaults and produce a localhost DSN when nothing is set.
func TestLoad_DatabaseDefaults(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	want := "host=localhost port=5432 user=postgres password=postgres dbname=employees sslmode=disable"
	if got := s.DatabaseDSN(); got != want {
		t.Fatalf("default DatabaseDSN() = %q, want %q", got, want)
	}
}

// TestSettings_Validate verifies the validate method directly across boundary
// values.
func TestSettings_Validate(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"min valid", 1, false},
		{"max valid", maxPort, false},
		{"typical", 3000, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"above max", maxPort + 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Settings{Port: tt.port}.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validate(Port=%d) err = %v, wantErr = %v", tt.port, err, tt.wantErr)
			}
		})
	}
}
