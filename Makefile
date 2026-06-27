.PHONY: install-tools install-hooks lint test test-cover cover-check

# Minimum total statement coverage (percent) enforced by cover-check.
COVERAGE_THRESHOLD ?= 90

# Install the dev tools (lefthook + golangci-lint) without brew, via go install.
install-tools:
	go install github.com/evilmartians/lefthook/v2@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

# Install the git hooks (run once after install-tools).
install-hooks:
	lefthook install

# Run golangci-lint (lint + formatting checks, no auto-fix).
lint:
	golangci-lint run ./...

# Run the unit tests with a quick coverage summary per package.
test:
	go test -cover ./...

# Run the unit tests, write a coverage profile, and print total coverage.
test-cover:
	go test -coverprofile=coverage.out -covermode=atomic ./... && go tool cover -func=coverage.out | tail -1

# Run tests with coverage and fail if total coverage is below COVERAGE_THRESHOLD.
cover-check: test-cover
	@total=$$(go tool cover -func=coverage.out | awk 'END { gsub(/%/, "", $$NF); print $$NF }'); \
	awk -v total="$$total" -v min="$(COVERAGE_THRESHOLD)" 'BEGIN { \
		printf "total coverage: %s%% (minimum: %s%%)\n", total, min; \
		if (total + 0 < min + 0) { \
			printf "FAIL: coverage %s%% is below the %s%% threshold\n", total, min; \
			exit 1; \
		} \
	}'
