FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o profile-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/profile-service .
EXPOSE 8082
CMD ["./profile-service"]
