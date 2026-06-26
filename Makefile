.PHONY: test test-cover

# Run the unit tests with a quick coverage summary per package.
test:
	go test -cover ./...

# Run the unit tests, write a coverage profile, and print total coverage.
test-cover:
	go test -coverprofile=coverage.out -covermode=atomic ./... && go tool cover -func=coverage.out | tail -1
