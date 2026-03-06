FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o profile-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/profile-service .
EXPOSE 8082
CMD ["./profile-service"]
