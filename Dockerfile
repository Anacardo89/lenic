# Stage 1: Build the Go application
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download
COPY . .
RUN go build -o lenic cmd/main.go

# Stage 2: Create the runtime image
FROM debian:12-slim 
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/lenic .

# Copy necessary assets
COPY cmd/ssl/ /root/ssl/
COPY cmd/templates/ /root/templates
COPY cmd/static/ /root/static

# Add the wait-for-it script
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

# Expose ports
EXPOSE 8081
EXPOSE 8082

# Use wait-for-it to check RabbitMQ readiness before starting the app
CMD ["wait-for-it", "rabbitmq:5672", "--timeout=30", "--strict", "--", "./lenic"]
