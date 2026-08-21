# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# ---- Runtime stage ----
FROM alpine:3.20

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /build/main .

# Copy your client assets and create an empty data directory for runtime storage
COPY src/client ./client
RUN mkdir -p data

EXPOSE 8080

CMD ["./main"]