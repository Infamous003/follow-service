FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o follow-service ./cmd/main.go

# second stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/follow-service .
COPY migrations ./migrations

EXPOSE 8080
CMD ["./follow-service"]
