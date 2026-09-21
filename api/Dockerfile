FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum (if present)
COPY go.mod ./
RUN go mod download || true

# Copy source code
COPY . .

# Build single API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /server .

EXPOSE 8080

CMD ["./server"]
