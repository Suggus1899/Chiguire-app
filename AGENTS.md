# Chiguire — ERP offline-first multi-tenant

## Stack
- **API**: Go 1.26 + chi + pgx/v5 | `apps/api/`
- **Web**: Next.js 16 (App Router) + TS + Tailwind + PowerSync | `apps/web/`
- **Mobile**: Flutter + Riverpod + PowerSync | `apps/mobile/`
- **Desktop**: Tauri 2 (Phase 3) | `apps/desktop/`
- **DB**: PostgreSQL 16 (WAL logical) + Row-Level Security
- **Sync**: PowerSync (self-hosted)
- **Migrations**: goose | `infra/migrations/`

## Ports (local dev)
- API: http://localhost:3001
- Postgres: localhost:5433 (Docker) or localhost:5432 (native)
- PowerSync: http://localhost:8080

## Setup
```
# Start infra
make infra-up

# Copy env
cp apps/api/.env.example apps/api/.env

# Run migrations
DATABASE_URL=postgresql://chiguire:chiguire_dev@localhost:5433/chiguire make migrate

# Start API
make api-dev

# Start web
make web-dev
```

## Commands
| Task | Command |
|---|---|
| Start Docker infra | `make infra-up` |
| Run migrations | `make migrate` |
| New migration | `make migrate-create name=add_customers` |
| Build API | `cd apps/api && go build ./...` |
| Test API | `cd apps/api && go test ./...` |
| Dev web | `pnpm --filter web dev` |
| Build web | `pnpm --filter web build` |
| Flutter analyze | `cd apps/mobile && flutter analyze` |
| Flutter test | `cd apps/mobile && flutter test` |

## Multi-tenant rules
- Every business table has `tenant_id UUID NOT NULL`.
- RLS policy: `USING (tenant_id = current_setting('app.tenant_id', true)::uuid)`.
- API sets `SET LOCAL app.tenant_id = $1` at request start.
- PowerSync sync rules replicate only rows matching `token_parameters.tenant_id`.
- **Never** bypass RLS for business queries. Use a superuser connection only in migrations.

## Offline-first rules
- UI reads always from local SQLite (PowerSync).
- Simple writes (drafts, edits): local first, PowerSync syncs to Postgres.
- Authoritative writes (emit invoice number, process payment): POST to API, mark `pending_emission` locally while offline, reconcile on reconnect.
- Stock tracked via `stock_movements` (append-only). Never update absolute quantity directly.

## Git Hooks
A pre-commit hook (`scripts/pre-commit.sh`) runs Go vet, Go build, and the
web build before each commit is accepted.

Install it once per clone by pointing Git at the `scripts` directory:
```bash
# Install hooks (run once)
git config core.hooksPath scripts
```
The hook file is `scripts/pre-commit.sh`. Ensure it is executable
(`chmod +x scripts/pre-commit.sh` on Unix). On Windows, Git for Windows
honors `core.hooksPath` as well.

## CI
GitHub Actions runs on every push/PR to `main` (see
`.github/workflows/ci.yml`):
- **Go API**: build, test, vet (with a Postgres 16 service container).
- **Web (Next.js)**: install with pnpm and build.
- **Flutter Mobile**: analyze and test.

## Phases
| Phase | What |
|---|---|
| 0 (current) | Monorepo, auth, multi-tenant, PowerSync sync test |
| 1 | Invoicing + inventory (offline on all 3 platforms) |
| 2 | Multi-currency + Venezuelan fiscal (IGTF, IVA, ISLR) |
| 3 | Fiscal printer (Tauri desktop) |
| 4 | Payment integrations (Cashea, Spidi, WayuPay, Biopago) |
| 5 | Purchases, quotations, commissions, routes |
| 6 | N8N, Excel import/export, subscription billing |
