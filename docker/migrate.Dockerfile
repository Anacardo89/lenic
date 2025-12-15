FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o migrate ./cmd/migrate

FROM alpine:3.18
WORKDIR /opt/migrate_svc
RUN mkdir -p ./bin
COPY --from=builder /app/migrate ./bin
COPY db/migrations ./migrations
COPY config/config.yaml ./config/config.yaml

ENTRYPOINT ["./bin/migrate"]