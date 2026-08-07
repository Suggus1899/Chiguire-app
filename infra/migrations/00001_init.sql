-- +goose Up
-- +goose StatementBegin

-- Extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Multi-tenant foundation
CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    plan        TEXT NOT NULL DEFAULT 'trial' CHECK (plan IN ('trial', 'starter', 'pyme', 'gold', 'enterprise')),
    trial_ends_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name     TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_tenants (
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role      TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, tenant_id)
);

CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON refresh_tokens(token_hash);
CREATE INDEX ON refresh_tokens(user_id);

-- Shared reference data (no tenant_id — all tenants receive this via sync)
CREATE TABLE exchange_rates (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    currency   TEXT NOT NULL,          -- e.g. 'USD', 'EUR'
    rate_to_ves NUMERIC(20, 8) NOT NULL,
    source     TEXT NOT NULL DEFAULT 'bcv',
    effective_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON exchange_rates(currency, effective_at);

-- Proof-of-concept table for Phase 0 sync test
CREATE TABLE test_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    label       TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON test_items(tenant_id);

-- Enable RLS
ALTER TABLE test_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE exchange_rates ENABLE ROW LEVEL SECURITY;

-- RLS policies (API sets 'app.tenant_id' at request start)
CREATE POLICY tenant_isolation ON test_items
    USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- exchange_rates: readable by any authenticated context
CREATE POLICY allow_read ON exchange_rates FOR SELECT USING (true);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS test_items CASCADE;
DROP TABLE IF EXISTS exchange_rates CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS user_tenants CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;
DROP EXTENSION IF EXISTS "pgcrypto";
-- +goose StatementEnd
