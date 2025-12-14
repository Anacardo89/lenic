FROM golang:1.25.0-trixie AS builder
ARG APP_PATH=/lenic
WORKDIR $APP_PATH
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o lenic ./cmd/main


FROM debian:trixie-slim
ARG APP_PATH=/lenic
WORKDIR $APP_PATH
RUN apt-get update && apt-get install -y ca-certificates netcat-openbsd
RUN rm -rf /var/lib/apt/lists/*
COPY --from=builder $APP_PATH .
COPY frontend/templates/ $APP_PATH/templates
COPY frontend/static/ $APP_PATH/static
ENV APP_ENV=docker
ENV ROOT_PATH=/lenic

ENV PORT=8080

EXPOSE ${PORT}

ENTRYPOINT ["/lenic/lenic"]
