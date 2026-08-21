# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy module files first so dependency downloads are cached
# separately from source changes (faster rebuilds during dev)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# ---- Runtime stage ----
FROM alpine:3.20

WORKDIR /app

# Copy just the compiled binary from the build stage — none of the
# Go toolchain or source ends up in the final image
COPY --from=builder /build/main .

# data/ and client/ are created here so the image runs correctly even
# before any volumes are mounted; docker-compose will override them
# with persistent volumes at runtime (see docker-compose.yml)
RUN mkdir -p data client

EXPOSE 8080

CMD ["./main"]