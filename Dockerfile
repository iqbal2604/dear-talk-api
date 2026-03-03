# ─── Build Stage ──────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git untuk go mod
RUN apk add --no-cache git

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate wire
RUN go install github.com/google/wire/cmd/wire@latest
RUN cd cmd/server && wire

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -o dear-talk \
    ./cmd/server/main.go ./cmd/server/wire_gen.go

# ─── Run Stage ────────────────────────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

# Add CA certs dan timezone
RUN apk --no-cache add ca-certificates tzdata

# Copy binary dari builder
COPY --from=builder /app/dear-talk .

EXPOSE 8080

CMD ["./dear-talk"]