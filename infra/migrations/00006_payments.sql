-- +goose Up
-- +goose StatementBegin

-- Fase 4: Pagos digitales (links de pago, webhooks)

CREATE TABLE payment_links (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id  UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    provider    TEXT NOT NULL CHECK (provider IN ('cashea', 'spidi', 'wayupay', 'biopago')),
    short_code  TEXT NOT NULL,           -- chiguire.app/pay/{code}
    amount_usd  NUMERIC(20,4) NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'expired', 'failed')),
    provider_ref TEXT,                   -- referencia externa
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ,
    paid_at     TIMESTAMPTZ
);

CREATE INDEX ON payment_links(tenant_id);
CREATE INDEX ON payment_links(invoice_id);
CREATE UNIQUE INDEX ON payment_links(short_code);
ALTER TABLE payment_links ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON payment_links USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Webhooks entrantes de pasarelas
CREATE TABLE payment_webhook_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    provider    TEXT NOT NULL,
    event_type  TEXT NOT NULL,
    payload     JSONB NOT NULL,
    signature   TEXT,
    processed   BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON payment_webhook_events(tenant_id);
CREATE INDEX ON payment_webhook_events(provider, event_type);
ALTER TABLE payment_webhook_events ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON payment_webhook_events USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_webhook_events CASCADE;
DROP TABLE IF EXISTS payment_links CASCADE;
-- +goose StatementEnd
