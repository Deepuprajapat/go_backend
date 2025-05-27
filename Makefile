.PHONY: build run test clean docker-up docker-down

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run: build
	./bin/server

# Run tests
test: docker-up
	@echo "Waiting for MySQL to be ready..."
	@sleep 5
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Install dependencies
deps:
	go mod tidy

# Run database migrations
migrate:
	atlas migrate apply

# Generate ent code
generate:
	go generate ./ent 