# backend/Dockerfile

# ✅ Base image
FROM golang:1.23-alpine

# ✅ Install tools
RUN apk add --no-cache git ca-certificates build-base unzip

# ✅ Working directory
WORKDIR /app

# ✅ Copy Go mod files & download deps
COPY go.mod go.sum ./
RUN go mod download

# ✅ Copy all backend source
COPY . .

# ✅ Extract frontend.zip into ./build (match Go code path)
COPY ./static-zip/frontend.zip ./frontend.zip
RUN unzip -o ./frontend.zip && \
    mkdir -p ./build && \
    mv build/* ./build/ && \
    rm -rf build frontend.zip

# ✅ Generate Ent code
RUN go generate ./ent

# ✅ Add non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

RUN chown -R appuser:appgroup /app
USER appuser

# ✅ Expose backend port
EXPOSE 9999

# ✅ Start the server
CMD ["go", "run", "./cmd/server"]

