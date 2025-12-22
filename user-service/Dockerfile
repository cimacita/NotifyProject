FROM golang:1.25.1-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o user-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/user-service .
CMD ["./user-service"]