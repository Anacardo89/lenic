#!/bin/sh

DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

migrate -path=/migrations -database "$DB_DSN" up
echo "Migrations complete."