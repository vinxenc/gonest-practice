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
// as a PostgreSQL connection URL.
func TestSettings_DatabaseDSN(t *testing.T) {
	s := Settings{
		DBHost:     "db.example.com",
		DBPort:     5433,
		DBUser:     "alice",
		DBPassword: "s3cret",
		DBName:     "employees",
		DBSSLMode:  "require",
	}
	want := "postgres://alice:s3cret@db.example.com:5433/employees?sslmode=require"
	if got := s.DatabaseDSN(); got != want {
		t.Fatalf("DatabaseDSN() = %q, want %q", got, want)
	}
}

// TestSettings_DatabaseDSN_EscapesCredentials verifies special characters in the
// user, password, and database name are percent-encoded so the DSN stays valid.
func TestSettings_DatabaseDSN_EscapesCredentials(t *testing.T) {
	s := Settings{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "foo bar",
		DBPassword: "p@ss:w/rd?",
		DBName:     "my db",
		DBSSLMode:  "disable",
	}
	want := "postgres://foo%20bar:p%40ss%3Aw%2Frd%3F@localhost:5432/my%20db?sslmode=disable"
	if got := s.DatabaseDSN(); got != want {
		t.Fatalf("DatabaseDSN() = %q, want %q", got, want)
	}
}

// TestLoad_DatabaseDefaults verifies the database settings fall back to their
// defaults and produce a localhost DSN when nothing is set.
func TestLoad_DatabaseDefaults(t *testing.T) {
	// Clear any ambient config so the defaults are what's exercised, regardless
	// of the developer's shell or CI environment.
	for _, k := range []string{
		"PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
	} {
		t.Setenv(k, "")
	}

	s, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	want := "postgres://postgres:postgres@localhost:5432/employees?sslmode=disable"
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
