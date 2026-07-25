.PHONY: help build run dev test lint clean docker-up docker-down migrate-up migrate-down migrate-up-docker migrate seed seed-docker migrate-seed db-bootstrap db-reset-seed

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
	@echo "  make migrate      - Run migrations in Docker (no local deps)"
	@echo "  make seed-docker  - Run seeder in Docker (no local deps)"
	@echo "  make migrate-seed - Run migrate + seed in Docker"
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

migrate-up-docker:
	docker run --rm --network finances-go_default -v "$$PWD/migrations:/migrations" migrate/migrate -path=/migrations -database "postgres://root:root@db:5432/root?sslmode=disable" up

migrate:
	docker compose up -d db
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' $$(docker compose ps -q db))" = "healthy" ]; do true; done
	docker run --rm --network finances-go_default -v "$$PWD/migrations:/migrations" migrate/migrate -path=/migrations -database "postgres://root:root@db:5432/root?sslmode=disable" up

seed:
	go run ./cmd/seed

seed-docker:
	docker compose up -d db
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' $$(docker compose ps -q db))" = "healthy" ]; do true; done
	docker run --rm --network finances-go_default -v "$$PWD:/app" -w /app -e DATABASE_URL="postgres://root:root@db:5432/root?sslmode=disable" -e DATABASE_HOST=db -e DATABASE_PORT=5432 -e DATABASE_USER=root -e DATABASE_PASSWORD=root -e DATABASE_NAME=root golang:1.26-alpine go run ./cmd/seed

migrate-seed: migrate seed-docker

db-bootstrap:
	docker compose up -d db
	docker run --rm --network finances-go_default -v "$$PWD/migrations:/migrations" migrate/migrate -path=/migrations -database "postgres://root:root@db:5432/root?sslmode=disable" up
	docker run --rm --network finances-go_default -v "$$PWD:/app" -w /app -e DATABASE_URL="postgres://root:root@db:5432/root?sslmode=disable" -e DATABASE_HOST=db -e DATABASE_PORT=5432 -e DATABASE_USER=root -e DATABASE_PASSWORD=root -e DATABASE_NAME=root golang:1.26-alpine go run ./cmd/seed

db-reset-seed:
	docker compose down -v
	docker compose up -d db
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' $$(docker compose ps -q db))" = "healthy" ]; do true; done
	docker run --rm --network finances-go_default -v "$$PWD/migrations:/migrations" migrate/migrate -path=/migrations -database "postgres://root:root@db:5432/root?sslmode=disable" up
	docker run --rm --network finances-go_default -v "$$PWD:/app" -w /app -e DATABASE_URL="postgres://root:root@db:5432/root?sslmode=disable" -e DATABASE_HOST=db -e DATABASE_PORT=5432 -e DATABASE_USER=root -e DATABASE_PASSWORD=root -e DATABASE_NAME=root golang:1.26-alpine go run ./cmd/seed

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
