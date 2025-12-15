FROM golang:1.25.0-trixie AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o lenic ./cmd/main


FROM debian:trixie-slim
ARG APP_PATH=/opt/lenic
WORKDIR $APP_PATH
RUN apt-get update && apt-get install -y ca-certificates netcat-openbsd
RUN rm -rf /var/lib/apt/lists/*
RUN mkdir -p ./bin
COPY --from=builder /app/lenic ./bin
COPY config/config.yaml $APP_PATH/config/config.yaml
COPY frontend/templates/ $APP_PATH/templates
COPY frontend/static/ $APP_PATH/static

ENV PORT=8080

EXPOSE ${PORT}

ENTRYPOINT ["./bin/lenic"]
