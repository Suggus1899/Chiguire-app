-- +goose Up
-- +goose StatementBegin

-- Fase 6: SaaS, webhooks salientes, onboarding

-- Webhooks salientes (configurados por el tenant)
CREATE TABLE outgoing_webhooks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    url         TEXT NOT NULL,
    events      TEXT[] NOT NULL DEFAULT '{}',  -- 'invoice.created', 'invoice.paid', 'stock.low'
    secret      TEXT,                           -- para firmar payloads
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON outgoing_webhooks(tenant_id);
ALTER TABLE outgoing_webhooks ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON outgoing_webhooks USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Cola de envío de webhooks con reintentos
CREATE TABLE webhook_deliveries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    webhook_id      UUID NOT NULL REFERENCES outgoing_webhooks(id) ON DELETE CASCADE,
    event_type      TEXT NOT NULL,
    payload         JSONB NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'delivered', 'failed')),
    attempts        INTEGER NOT NULL DEFAULT 0,
    last_response   TEXT,
    next_retry_at   TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON webhook_deliveries(tenant_id);
CREATE INDEX ON webhook_deliveries(status, next_retry_at);
ALTER TABLE webhook_deliveries ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON webhook_deliveries USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Suscripciones SaaS (Chiguire como producto)
CREATE TABLE subscriptions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan            TEXT NOT NULL CHECK (plan IN ('trial', 'emprendedor', 'pyme', 'gold', 'enterprise')),
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'past_due', 'canceled', 'trialing', 'read_only')),
    stripe_customer_id  TEXT,
    stripe_subscription_id TEXT,
    current_period_start TIMESTAMPTZ,
    current_period_end   TIMESTAMPTZ,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON subscriptions(tenant_id);
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON subscriptions USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Límites por plan
CREATE TABLE plan_limits (
    plan            TEXT PRIMARY KEY,
    max_invoices_month INTEGER NOT NULL,
    max_branches    INTEGER NOT NULL,
    max_users       INTEGER NOT NULL,
    features        JSONB NOT NULL DEFAULT '{}'
);

INSERT INTO plan_limits (plan, max_invoices_month, max_branches, max_users, features) VALUES
    ('trial',       10,   1, 1, '{"fiscal": false, "payments": false}'::jsonb),
    ('emprendedor', 100,  1, 2, '{"fiscal": true,  "payments": false}'::jsonb),
    ('pyme',        1000, 3, 5, '{"fiscal": true,  "payments": true}'::jsonb),
    ('gold',        10000,10, 20,'{"fiscal": true, "payments": true, "webhooks": true}'::jsonb),
    ('enterprise',  0,    0, 0, '{"fiscal": true, "payments": true, "webhooks": true, "custom": true}'::jsonb);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions CASCADE;
DROP TABLE IF EXISTS webhook_deliveries CASCADE;
DROP TABLE IF EXISTS outgoing_webhooks CASCADE;
DROP TABLE IF EXISTS plan_limits CASCADE;
-- +goose StatementEnd
