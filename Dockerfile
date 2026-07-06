FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/api/main.go

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
