# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Runtime stage
FROM alpine

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

# Expose port matching Docker Compose configuration
EXPOSE 8080

# Set default environment variables matching .env
ENV DB_HOST=db \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=event_system \
    JWT_SECRET=${JWT_SECRET}

# Run the application
CMD ["./main"]