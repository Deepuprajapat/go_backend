.PHONY: build run test clean deps migrate

# Build the application
build:
	go build -o bin/server main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	@if ! docker ps --format '{{.Names}}' | grep -q '^im_mysql$$'; then \
		echo 'Starting MySQL with docker-compose...'; \
		docker-compose up -d; \
	else \
		echo 'MySQL container already running.'; \
	fi
	go test -v

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