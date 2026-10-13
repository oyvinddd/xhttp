.PHONY: format build issues test

# Format source code
format:
	go fmt ./...

# Compiles all packages without producing an executable
build:
	go build ./...

# Runs static code analysis
issues:
	go vet ./...

# Compiles packages and runs tests
test:
	go test ./...
