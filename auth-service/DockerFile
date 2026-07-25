FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/app

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/auth-service .
EXPOSE 8081

CMD ["./auth-service"]