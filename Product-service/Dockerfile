FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build product-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY  --from=builder /app/product-service .
COPY .env .
EXPOSE 50054
CMD ["/product-service"]
