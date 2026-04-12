# Multi-stage build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build tools
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN make build VERSION=$VERSION

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/go-template .

# User
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

USER app

ENTRYPOINT ["./go-template"]
CMD ["--help"]
