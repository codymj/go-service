.PHONY: up down refresh logs test

up:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
	fi
	docker compose up -d

down:
	docker compose down

refresh:
	@if [ ! -f .env ]; then \
		cp .env.example .env \
	fi
	docker compose down --remove-orphans
	docker compose build --no-cache
	docker compose up -d --build -V
	docker compose exec go-service go run . migrate up

migrate:
	docker compose exec go-service go run . migrate up

logs:
	docker compose logs -f

test:
	go test -v -race -cover -count=1 -failfast ./...
