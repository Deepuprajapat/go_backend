# Single stage build with all source code
FROM golang:1.23-alpine

# Install git, ca-certificates, build tools, unzip
RUN apk add --no-cache git ca-certificates build-base unzip

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all backend source code
COPY . .

# ✅ Unzip frontend.zip (copied via CI/CD)
RUN unzip -o ./static-zip/frontend.zip -d ./frontend && rm -f ./static-zip/frontend.zip

# Generate Ent code
RUN go generate ./ent

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Change ownership of all files to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose backend port
EXPOSE 9999

# ✅ Start the server
CMD ["go", "run", "./cmd/server"]

