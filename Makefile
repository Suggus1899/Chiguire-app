.PHONY: db-create db-drop db-wal migrate migrate-down migrate-create api-gen api-dev web-dev infra-up powersync-up powersync-down

# ─── Local Postgres (no Docker) ─────────────────────────────────────────
DB_URL ?= postgresql://postgres:1234@localhost:5432/chiguire

db-create:
	psql -U postgres -c "CREATE DATABASE chiguire;" 2>/dev/null || echo "DB already exists"

db-drop:
	psql -U postgres -c "DROP DATABASE IF EXISTS chiguire;"

# Run once to enable logical replication (requires Postgres restart)
db-wal:
	psql -U postgres -c "ALTER SYSTEM SET wal_level = logical;"
	psql -U postgres -c "ALTER SYSTEM SET max_replication_slots = 10;"
	psql -U postgres -c "ALTER SYSTEM SET max_wal_senders = 10;"
	@echo "Restart Postgres for changes to take effect."

# ─── Migrations ──────────────────────────────────────────────────────────
migrate:
	goose -dir infra/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir infra/migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir infra/migrations postgres "$(DB_URL)" status

migrate-create:
	goose -dir infra/migrations create $(name) sql

# ─── API ─────────────────────────────────────────────────────────────────
api-gen:
	cd apps/api && sqlc generate

api-dev:
	cd apps/api && go run ./cmd/server

api-build:
	cd apps/api && go build -o bin/server ./cmd/server

# ─── Web ─────────────────────────────────────────────────────────────────
web-dev:
	pnpm --filter web dev

web-build:
	pnpm --filter web build

# ─── PowerSync (only service that uses Docker) ───────────────────────────
infra-up: powersync-up

powersync-up:
	docker compose -f infra/docker/compose.yml up -d

powersync-down:
	docker compose -f infra/docker/compose.yml down

# ─── Flutter ─────────────────────────────────────────────────────────────
mobile-analyze:
	cd apps/mobile && flutter analyze

mobile-test:
	cd apps/mobile && flutter test
