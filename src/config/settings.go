// Package config holds the application's validated runtime configuration. It
// reads values from environment variables, applies defaults, and validates them
// up front so the rest of the application can depend on a single, already-valid
// Settings value.
package config

import (
	"fmt"

	envvalidator "github.com/philiprehberger/go-env-validator"
)

// maxPort is the largest valid TCP port number.
const maxPort = 65535

// Settings is the validated application configuration. Each field is sourced
// from an environment variable via the `env` struct tag; go-env-validator
// applies defaults and reports every problem at once.
type Settings struct {
	// Port is the TCP port the HTTP server listens on.
	Port int `env:"PORT,default=3000"`

	// DBHost is the PostgreSQL server hostname.
	DBHost string `env:"DB_HOST,default=localhost"`
	// DBPort is the PostgreSQL server port.
	DBPort int `env:"DB_PORT,default=5432"`
	// DBUser is the PostgreSQL user to connect as.
	DBUser string `env:"DB_USER,default=postgres"`
	// DBPassword is the PostgreSQL user's password.
	DBPassword string `env:"DB_PASSWORD,default=postgres"`
	// DBName is the PostgreSQL database to connect to.
	DBName string `env:"DB_NAME,default=employees"`
	// DBSSLMode is the libpq sslmode (disable, require, verify-full, ...).
	DBSSLMode string `env:"DB_SSLMODE,default=disable"`
}

// DatabaseDSN builds a libpq-style connection string from the database settings,
// suitable for passing to the GORM PostgreSQL driver.
func (s Settings) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.DBHost, s.DBPort, s.DBUser, s.DBPassword, s.DBName, s.DBSSLMode,
	)
}

// validate performs semantic checks that the `env` struct tags cannot express,
// so misconfiguration is caught at load time rather than later at net.Listen.
func (s Settings) validate() error {
	if s.Port < 1 || s.Port > maxPort {
		return fmt.Errorf("PORT must be between 1 and %d, got %d", maxPort, s.Port)
	}
	return nil
}

// Load reads configuration from the environment, applies defaults, and
// validates it — both the per-field tag validation from go-env-validator and
// the semantic checks in validate. On failure it returns an error describing
// the problem, so misconfiguration surfaces up front.
//
// Its (*Settings, error) signature doubles as an fx provider: fx constructs the
// Settings value lazily and fails app startup if validation fails.
func Load() (*Settings, error) {
	var s Settings
	if err := envvalidator.Validate(&s); err != nil {
		return nil, err
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return &s, nil
}
