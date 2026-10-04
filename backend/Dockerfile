# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /calculator-api ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /calculator-api /app/calculator-api

EXPOSE 8080

USER 1000:1000

ENTRYPOINT ["/app/calculator-api"]
