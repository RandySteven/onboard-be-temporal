CONFIG ?= ./files/yaml/app.local.yml

.PHONY: up down logs migrate drop run tidy

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f mysql temporal temporal-ui

migrate:
	go run ./cmd/migration --config $(CONFIG)

drop:
	go run ./cmd/drop --config $(CONFIG)

run:
	go run ./cmd/main --config $(CONFIG)

tidy:
	go mod tidy
