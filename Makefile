.PHONY: build run test test-integration test-unit clean docker-up docker-down deps generate migrate migrate-status migrate-validate migrate-diff migrate-reset dev-setup db-reset test-setup reset

# Build the application
build:
	go build -o bin/server cmd/server/main.go


# Seed testimonials into static_site_data
seed-testimonials:
	docker cp sql/seed_testimonials.sql im_postgres_db:/tmp/seed_testimonials.sql
	docker exec im_postgres_db psql -U im_db_dev -d mydb -f /tmp/seed_testimonials.sql

# Seed amenities into static_site_data
seed-amenities:
	PGPASSWORD=password psql -h localhost -p 5434 -U im_db_dev -d mydb -f sql/seed_amenities.sql

# Run the application
build-run: build
	./bin/server

# Run all tests (no external dependencies needed - testcontainers handles DB)
test:
	go test -v ./...

# Run only unit tests (excluding integration tests)
test-unit:
	go test -v ./... -short

# Run integration tests (self-contained with testcontainers)
test-integration:
	go test -v ./tests/integration/...

# Run integration tests with coverage
test-integration-coverage:
	go test -v -coverprofile=coverage.out ./tests/integration/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific integration test suites
test-auth:
	go test -v ./tests/integration/ -run TestAuthEndpoints

test-projects:
	go test -v ./tests/integration/ -run TestProjectEndpoints

test-properties:
	go test -v ./tests/integration/ -run TestPropertyEndpoints

test-leads:
	go test -v ./tests/integration/ -run TestLeadEndpoints

# Legacy test-setup command (now just a message since testcontainers handles everything)
test-setup:
	@echo "Tests now use testcontainers - no manual setup required!"
	@echo "Docker containers are created automatically during test execution."

# Clean build artifacts
clean:
	rm -rf bin/

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Reset database with clean slate (removes all data and recreates)
reset:
	@echo "🔄 Resetting database with clean slate..."
	@echo "⚠️  This will destroy ALL existing data!"
	@sleep 3
	@echo "📦 Stopping and removing containers and volumes..."
	docker-compose down -v
	@echo "🚀 Starting fresh database container..."
	docker-compose up -d
	@echo "⏳ Waiting for database to be ready..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker exec im_postgres_db pg_isready -U im_db_dev -d mydb >/dev/null 2>&1; then \
			echo "✅ Database is ready!"; \
			break; \
		fi; \
		echo "⏱️  Waiting for database... ($$timeout seconds remaining)"; \
		sleep 2; \
		timeout=$$((timeout-2)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "❌ Database failed to start within 60 seconds"; \
		exit 1; \
	fi
	@echo "🔧 Running migrations..."
	make migrate
	@echo "🌱 Seeding data..."
	make seed-data
	make seed-projects
	make seed-testimonials
	@echo "🎉 Database reset complete! Clean slate ready for development."

# Reset database by dropping all tables
db-reset: 
	@echo "Dropping all tables from database..."
	@sleep 3
	@docker exec -i im_postgres_db psql -U im_db_dev -d mydb -c "\
		DO \$$\$$ \
		DECLARE \
			r RECORD; \
		BEGIN \
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP \
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE'; \
			END LOOP; \
		END \$$\$$;"
	@echo "All tables dropped successfully!"
	@echo "You can now run 'make run' to apply fresh migrations."

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
	go run cmd/server/main.go run-migration ./migration_jobs/database_export
	make seed-testimonials
	make seed-amenities
	


# Setup fresh development environment
dev-setup: docker-up migrate
	@echo "Development environment ready!" 

video-migration:
	go run cmd/server/main.go video-migration

# Mohan's Commands, If you find use-full, move them up but please dont remove.
migrate-schema: run-migration

seed-data:
	go run cmd/server/main.go seed-admin

migrate: migrate-schema

export-database:
	go run cmd/server/main.go export-database

export-specific-tables:
	go run cmd/server/main.go export-specific-tables

initialize-json-loader:
	go run cmd/server/main.go initialize-json-loader ./migration_jobs/database_export

seed-projects:
	go run cmd/server/main.go seed-projects

