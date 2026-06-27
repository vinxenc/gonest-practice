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
