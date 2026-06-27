.PHONY: install-tools install-hooks test test-cover

# Install the dev tools (lefthook + golangci-lint) without brew, via go install.
install-tools:
	go install github.com/evilmartians/lefthook/v2@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

# Install the git hooks (run once after install-tools).
install-hooks:
	lefthook install

# Run the unit tests with a quick coverage summary per package.
test:
	go test -cover ./...

# Run the unit tests, write a coverage profile, and print total coverage.
test-cover:
	go test -coverprofile=coverage.out -covermode=atomic ./... && go tool cover -func=coverage.out | tail -1
