FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN  go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd

FROM alpine:latest

COPY --from=builder /app/auth-service .
COPY .env .

EXPOSE 50051

RUN ["./auth-service"]


