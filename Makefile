.PHONY: build run test clean docker-up docker-down deps generate migrate migrate-status migrate-validate migrate-diff migrate-reset dev-setup

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

# Generate ent code
generate:
	go generate ./ent

# Atlas Migration Commands
# Apply pending migrations to database
migrate: docker-up
	@echo "Waiting for MySQL to be ready..."
	@sleep 5
	atlas migrate apply --env local

# Check migration status
migrate-status:
	atlas migrate status --env local

# Validate migration files
migrate-validate:
	atlas migrate validate --env local

# Generate new migration (usage: make migrate-diff name=migration_name)
migrate-diff: generate
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide migration name. Usage: make migrate-diff name=your_migration_name"; \
		exit 1; \
	fi
	atlas migrate diff $(name) --env local --to "ent://ent/schema" --dev-url "mysql://root:password@localhost:3306/atlas_dev"

# Reset database and apply all migrations (DESTRUCTIVE - use with caution)
migrate-reset: docker-up
	@echo "WARNING: This will drop and recreate the database!"
	@echo "Waiting for MySQL to be ready..."
	@sleep 5
	@echo "Dropping database..."
	@docker exec im_mysql mysql -u root -ppassword -e "DROP DATABASE IF EXISTS im_db_dev; CREATE DATABASE im_db_dev;"
	@echo "Applying migrations..."
	atlas migrate apply --env local

# Setup fresh development environment
dev-setup: docker-up migrate
	@echo "Development environment ready!" 