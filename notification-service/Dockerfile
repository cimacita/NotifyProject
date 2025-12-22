FROM golang:1.25.1-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -ldflags '-s -w' -o /app/notification-service ./cmd/

FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/notification-service .
CMD ["./notification-service"]
