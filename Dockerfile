# Single stage build with all source code
FROM golang:1.23-alpine

# Install git, ca-certificates, and build tools
RUN apk add --no-cache git ca-certificates build-base

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Generate Ent code
RUN go generate ./ent

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Change ownership of all files to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 9999

# Build and run the application
CMD ["go", "run", "./cmd/server"] 