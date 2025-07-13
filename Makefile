.PHONY: build run test clean docker-up docker-down deps generate migrate migrate-status migrate-validate migrate-diff migrate-reset dev-setup

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
build-run: build
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

enums:
	./go-enum --marshal --sql --flag \
		--file ./internal/domain/enums/poject_configurations.go \
		--file ./internal/domain/enums/project_status.go


# Install dependencies
deps:
	go mod tidy

# Generate ent code
generate:
	go generate ./ent

run:
	go run cmd/server/main.go

run-migration:
	go run cmd/server/main.go run-migration

# Setup fresh development environment
dev-setup: docker-up migrate
	@echo "Development environment ready!" 

video-migration:
	go run cmd/server/main.go video-migration

# Mohan's Commands, If you find use-full, move them up but please dont remove.
migrate-schema:
	go run cmd/server/main.go run-migration

seed-data:
	@echo "Implement me"

migrate: migrate-schema