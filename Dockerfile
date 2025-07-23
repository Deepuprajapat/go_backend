FROM golang:1.23-alpine

RUN apk add --no-cache git ca-certificates build-base unzip

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY ./static-zip/frontend.zip ./frontend.zip

# ✅ Extract and copy build files safely
RUN unzip -o ./frontend.zip && \
    mkdir -p ./build && \
    cp -r frontend/build/* ./build/ && \
    rm -rf frontend frontend.zip

RUN go generate ./ent

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 9999

CMD ["go", "run", "./cmd/server"]

