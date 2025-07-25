# Dockerfile

### 1) Builder stage ###
FROM golang:1.23-alpine AS builder

# Install git for module fetches
RUN apk add --no-cache git

WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

### 2) Runner stage ###
FROM alpine:latest

# Add certificates so HTTPS calls work
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy built binary and config folder
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

# Expose your HTTP port
EXPOSE 8080

# Default config path (you can override via env)
ENV APP_CONFIG=/app/config/dev.yaml

# Entrypoint
ENTRYPOINT ["./server"]