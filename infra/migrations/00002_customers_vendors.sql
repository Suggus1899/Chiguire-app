-- +goose Up
-- +goose StatementBegin

-- Fase 1: Clientes y proveedores
CREATE TABLE customers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    rif         TEXT NOT NULL,
    name        TEXT NOT NULL,
    address     TEXT,
    phone       TEXT,
    email       TEXT,
    is_special_taxpayer BOOLEAN NOT NULL DEFAULT FALSE, -- contribuyente especial
    credit_limit NUMERIC(20,2) NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON customers(tenant_id);
CREATE UNIQUE INDEX ON customers(tenant_id, rif);
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON customers USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

CREATE TABLE vendors (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    rif         TEXT NOT NULL,
    name        TEXT NOT NULL,
    address     TEXT,
    phone       TEXT,
    email       TEXT,
    is_special_taxpayer BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON vendors(tenant_id);
CREATE UNIQUE INDEX ON vendors(tenant_id, rif);
ALTER TABLE vendors ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON vendors USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vendors CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
-- +goose StatementEnd
