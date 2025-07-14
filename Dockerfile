# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/main .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Set environment variables (override with docker run -e)
ENV PORT=8080
ENV HUGGING_FACE_API_KEY=""

# Expose port
EXPOSE $PORT

# Start the application
CMD ["./main"]