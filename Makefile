.PHONY: dev infra-up infra-down migrate migrate-create api-gen

# Start local infra (Postgres + PowerSync)
infra-up:
	docker compose -f infra/docker/compose.yml up -d

infra-down:
	docker compose -f infra/docker/compose.yml down

# Run DB migrations
migrate:
	cd apps/api && goose -dir ../../infra/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	cd apps/api && goose -dir ../../infra/migrations postgres "$(DATABASE_URL)" down

migrate-create:
	cd apps/api && goose -dir ../../infra/migrations create $(name) sql

# Generate sqlc
api-gen:
	cd apps/api && sqlc generate

# Run API locally
api-dev:
	cd apps/api && go run ./cmd/server

# Run web locally
web-dev:
	pnpm --filter web dev
