FROM golang:1.24.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/user-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/user-service .

COPY --from=builder /app/config ./config

COPY --from=builder /app/internal/lib/migrator/migrations ./internal/lib/migrator/migrations

EXPOSE 8080

CMD ["./user-service"]
