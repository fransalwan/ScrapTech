# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and build essentials
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/scrapflow-api ./cmd/api/main.go

# Runtime Stage
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/scrapflow-api /app/scrapflow-api
COPY .env.example /app/.env

EXPOSE 8080

CMD ["/app/scrapflow-api"]
