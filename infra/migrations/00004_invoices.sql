-- +goose Up
-- +goose StatementBegin

-- Fase 1: Facturación y pagos básicos

CREATE TABLE invoices (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    number          TEXT,                          -- XX-YYYYMMDD-NNNNN, asignado por API al emitir
    customer_id     UUID REFERENCES customers(id) ON DELETE SET NULL,
    vendor_user_id  UUID REFERENCES users(id) ON DELETE SET NULL, -- vendedor
    branch_id       UUID REFERENCES branches(id) ON DELETE SET NULL,
    status          TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'issued', 'paid', 'partial', 'void')),
    currency        TEXT NOT NULL DEFAULT 'USD',
    subtotal        NUMERIC(20,4) NOT NULL DEFAULT 0,
    discount_total  NUMERIC(20,4) NOT NULL DEFAULT 0,
    tax_total       NUMERIC(20,4) NOT NULL DEFAULT 0,
    total           NUMERIC(20,4) NOT NULL DEFAULT 0,
    exchange_rate_id UUID REFERENCES exchange_rates(id), -- tasa al momento de emisión
    issued_at       TIMESTAMPTZ,
    voided_at       TIMESTAMPTZ,
    notes           TEXT,
    pending_emission BOOLEAN NOT NULL DEFAULT FALSE, -- offline: encola emisión
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON invoices(tenant_id);
CREATE INDEX ON invoices(customer_id);
CREATE INDEX ON invoices(status);
CREATE UNIQUE INDEX ON invoices(tenant_id, number) WHERE number IS NOT NULL;
ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON invoices USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Secuencia de numeración por sucursal
CREATE TABLE invoice_sequences (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id   UUID REFERENCES branches(id) ON DELETE CASCADE,
    prefix      TEXT NOT NULL,          -- e.g. '01'
    last_seq    INTEGER NOT NULL DEFAULT 0,
    last_date   DATE NOT NULL,
    UNIQUE(tenant_id, branch_id, prefix, last_date)
);

ALTER TABLE invoice_sequences ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON invoice_sequences USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

CREATE TABLE invoice_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id      UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    product_id      UUID REFERENCES products(id) ON DELETE SET NULL,
    description     TEXT NOT NULL,
    qty             NUMERIC(20,4) NOT NULL,
    unit_price      NUMERIC(20,4) NOT NULL,
    discount_pct    NUMERIC(5,2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(20,4) NOT NULL DEFAULT 0,
    tax_rate        NUMERIC(5,2) NOT NULL DEFAULT 16,
    tax_amount      NUMERIC(20,4) NOT NULL DEFAULT 0,
    line_total      NUMERIC(20,4) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON invoice_items(tenant_id);
CREATE INDEX ON invoice_items(invoice_id);
ALTER TABLE invoice_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON invoice_items USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Pagos básicos (efectivo, transferencia)
CREATE TABLE invoice_payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id      UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    method          TEXT NOT NULL CHECK (method IN ('cash_usd', 'cash_ves', 'transfer', 'card', 'digital')),
    provider        TEXT,  -- 'cashea', 'spidi', etc. (Fase 4)
    amount_usd      NUMERIC(20,4) NOT NULL DEFAULT 0,
    amount_ves      NUMERIC(20,4) NOT NULL DEFAULT 0,
    rate_at_payment NUMERIC(20,8) NOT NULL DEFAULT 1,
    igt_amount      NUMERIC(20,4) NOT NULL DEFAULT 0, -- IGTF si aplica
    reference       TEXT,
    status          TEXT NOT NULL DEFAULT 'completed' CHECK (status IN ('pending', 'completed', 'failed')),
    pending_payment BOOLEAN NOT NULL DEFAULT FALSE,  -- offline
    paid_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON invoice_payments(tenant_id);
CREATE INDEX ON invoice_payments(invoice_id);
ALTER TABLE invoice_payments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON invoice_payments USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invoice_payments CASCADE;
DROP TABLE IF EXISTS invoice_items CASCADE;
DROP TABLE IF EXISTS invoice_sequences CASCADE;
DROP TABLE IF EXISTS invoices CASCADE;
-- +goose StatementEnd
