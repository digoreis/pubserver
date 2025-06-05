# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /pubserver ./cmd/server

# Production image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /pubserver /pubserver
COPY config ./config
EXPOSE 8080
ENTRYPOINT ["/pubserver"]