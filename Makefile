.PHONY: build run test clean deps migrate

# Build the application
build:
	go build -o bin/server main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy

# Run database migrations
migrate:
	atlas migrate apply

# Generate ent code
generate:
	go generate ./ent 