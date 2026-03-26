# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Run go install requirements if any native deps, alpine needs gcc for some cgo (pgvector-go does not mandate cgo strictly)
# but sql deps often do.
RUN apk add --no-cache gcc musl-dev

# Download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/api ./cmd/api

# Final minimal stage
FROM alpine:3.19

WORKDIR /app

# CA certificates are required for LLM HTTP calls
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /go/bin/api /bin/api

# We also need to copy migrations to wait for execution in entrypoint or runtime
COPY migrations /app/migrations

# Expose backend port
EXPOSE 8080

CMD ["/bin/api"]
