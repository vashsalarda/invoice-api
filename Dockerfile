# syntax=docker/dockerfile:1

FROM alpine AS certs
RUN apk update && apk add ca-certificates

# Build stage (added) - builds a Linux static Go binary
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build (adjust build target if your main is elsewhere)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /invoice-api ./cmd/api

FROM busybox:stable
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

# Copy the binary produced by the builder stage
COPY --from=builder /invoice-api invoice-api

EXPOSE 3000
CMD ["./invoice-api"]

