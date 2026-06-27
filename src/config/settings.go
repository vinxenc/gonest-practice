// Package config holds the application's validated runtime configuration. It
// reads values from environment variables, applies defaults, and validates them
// up front so the rest of the application can depend on a single, already-valid
// Settings value.
package config

import (
	envvalidator "github.com/philiprehberger/go-env-validator"
)

// Settings is the validated application configuration. Each field is sourced
// from an environment variable via the `env` struct tag; go-env-validator
// applies defaults and reports every problem at once.
type Settings struct {
	// Port is the TCP port the HTTP server listens on.
	Port int `env:"PORT,default=3000"`
}

// Load reads configuration from the environment, applies defaults, and
// validates it. On failure it returns a *envvalidator.ValidationError that
// describes every offending field, so misconfiguration surfaces in one shot.
//
// Its (*Settings, error) signature doubles as an fx provider: fx constructs the
// Settings value lazily and fails app startup if validation fails.
func Load() (*Settings, error) {
	var s Settings
	if err := envvalidator.Validate(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
