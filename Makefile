.PHONY: run dev build test migrate-up migrate-down

run:
	go run ./cmd/api

build:
	go build -o bin/cats-api ./cmd/api

dev:
	docker compose up -d db

migrate-up:
	psql "$$DB_DSN" -f db/migrations/0001_init.sql && \
	psql "$$DB_DSN" -f db/migrations/0002_cat_thumbs.sql

migrate-down:
	psql "$$DB_DSN" -c "DROP TABLE IF EXISTS cat_thumbnails;" && \
	psql "$$DB_DSN" -c "DROP TABLE IF EXISTS cats;"

test:
	go test ./...