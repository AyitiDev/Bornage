FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency definition files
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download

# Copy application source code
COPY . .

# Build single static API binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /server .

EXPOSE 8080

CMD ["./server"]
