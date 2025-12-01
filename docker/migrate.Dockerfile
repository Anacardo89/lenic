FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o migrate ./cmd/migrate

FROM alpine:3.18
WORKDIR /migrate_svc
COPY --from=builder /app/migrate .
COPY db/migrations ./db/migrations
COPY config/config.yaml ./config/config.yaml
COPY .env ./
ENV APP_ENV=docker
ENV ROOT_PATH=/migrate_svc

ENTRYPOINT ["./migrate"]