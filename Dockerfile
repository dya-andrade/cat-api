# Etapa de build
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /cats-api ./cmd/api

# Etapa final
FROM gcr.io/distroless/base-debian12
ENV APP_ADDR=0.0.0.0:8080
ENV DB_DSN=postgres://postgres:postgres@db:5432/cats?sslmode=disable
ENV DB_MAX_CONNS=10
ENV DB_MIN_CONNS=2
ENV DB_MAX_IDLE_TIME=30s
ENV WORKER_CONCURRENCY=4
ENV REQUEST_TIMEOUT=10s
COPY --from=builder /cats-api /cats-api
EXPOSE 8080
ENTRYPOINT ["/cats-api"]