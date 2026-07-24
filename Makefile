.PHONY: help build run dev test lint clean docker-up docker-down migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make dev          - Run with live reload (air)"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linter"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make docker-up    - Start all services"
	@echo "  make docker-down  - Stop all services"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./... -v -cover

lint:
	golangci-lint run

clean:
	rm -rf bin/ tmp/

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

seed:
	go run ./cmd/seed

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
