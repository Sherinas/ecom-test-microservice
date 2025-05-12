# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o admin-service ./cmd
# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/admin-service .
COPY cmd/.env .
EXPOSE 50053
CMD ["./admin-service"]