FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payment-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/payment-service .
COPY migrations ./migrations
EXPOSE 50054
CMD ["./payment-service"]