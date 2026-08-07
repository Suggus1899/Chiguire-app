-- +goose Up
-- +goose StatementBegin

-- Fase 2: Fiscal venezolano (IGTF, IVA, ISLR, retenciones, libros)

-- Categorías fiscales
CREATE TABLE tax_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,           -- 'general', 'reduced', 'exempt'
    iva_rate    NUMERIC(5,2) NOT NULL,   -- 16, 8, 0
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON tax_categories(tenant_id);
ALTER TABLE tax_categories ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON tax_categories USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Retenciones (IVA e ISLR)
CREATE TABLE tax_withholdings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id      UUID REFERENCES invoices(id) ON DELETE CASCADE,
    vendor_id       UUID REFERENCES vendors(id) ON DELETE SET NULL,
    customer_id     UUID REFERENCES customers(id) ON DELETE SET NULL,
    type            TEXT NOT NULL CHECK (type IN ('iva', 'islr')),
    base            NUMERIC(20,4) NOT NULL,
    rate            NUMERIC(5,2) NOT NULL,    -- 75, 100 (IVA); variable (ISLR)
    amount          NUMERIC(20,4) NOT NULL,
    document_number TEXT,                     -- comprobante
    period          TEXT NOT NULL,             -- YYYY-MM
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON tax_withholdings(tenant_id);
CREATE INDEX ON tax_withholdings(invoice_id);
ALTER TABLE tax_withholdings ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON tax_withholdings USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Períodos fiscales
CREATE TABLE fiscal_periods (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    period      TEXT NOT NULL,           -- 'YYYY-MM'
    is_closed   BOOLEAN NOT NULL DEFAULT FALSE,
    closed_at   TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, period)
);

ALTER TABLE fiscal_periods ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON fiscal_periods USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Libros fiscales (ventas/compras)
CREATE TABLE fiscal_books (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    type        TEXT NOT NULL CHECK (type IN ('sales', 'purchases')),
    period      TEXT NOT NULL,           -- 'YYYY-MM'
    entries_json JSONB NOT NULL DEFAULT '[]', -- entradas del libro
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, type, period)
);

CREATE INDEX ON fiscal_books(tenant_id);
ALTER TABLE fiscal_books ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON fiscal_books USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fiscal_books CASCADE;
DROP TABLE IF EXISTS fiscal_periods CASCADE;
DROP TABLE IF EXISTS tax_withholdings CASCADE;
DROP TABLE IF EXISTS tax_categories CASCADE;
-- +goose StatementEnd
