-- +goose Up
-- +goose StatementBegin

-- Fase 5: Compras, cotizaciones, comisiones, rutas

-- Órdenes de compra
CREATE TABLE purchase_orders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vendor_id   UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    number      TEXT,
    status      TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'approved', 'partial', 'received', 'void')),
    currency    TEXT NOT NULL DEFAULT 'USD',
    subtotal    NUMERIC(20,4) NOT NULL DEFAULT 0,
    tax_total   NUMERIC(20,4) NOT NULL DEFAULT 0,
    total       NUMERIC(20,4) NOT NULL DEFAULT 0,
    approved_at TIMESTAMPTZ,
    received_at TIMESTAMPTZ,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON purchase_orders(tenant_id);
CREATE INDEX ON purchase_orders(vendor_id);
ALTER TABLE purchase_orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON purchase_orders USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

CREATE TABLE purchase_order_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    product_id      UUID REFERENCES products(id) ON DELETE SET NULL,
    description     TEXT NOT NULL,
    qty_ordered     NUMERIC(20,4) NOT NULL,
    qty_received    NUMERIC(20,4) NOT NULL DEFAULT 0,
    unit_cost       NUMERIC(20,4) NOT NULL,
    line_total      NUMERIC(20,4) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON purchase_order_items(tenant_id);
CREATE INDEX ON purchase_order_items(purchase_order_id);
ALTER TABLE purchase_order_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON purchase_order_items USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Recepciones parciales
CREATE TABLE purchase_order_receipts (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    purchase_order_id   UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    po_item_id          UUID NOT NULL REFERENCES purchase_order_items(id) ON DELETE CASCADE,
    warehouse_id        UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    qty_received        NUMERIC(20,4) NOT NULL,
    received_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON purchase_order_receipts(tenant_id);
ALTER TABLE purchase_order_receipts ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON purchase_order_receipts USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Cotizaciones
CREATE TABLE quotations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    number      TEXT,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    status      TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'accepted', 'rejected', 'expired', 'converted')),
    currency    TEXT NOT NULL DEFAULT 'USD',
    subtotal    NUMERIC(20,4) NOT NULL DEFAULT 0,
    discount_total NUMERIC(20,4) NOT NULL DEFAULT 0,
    tax_total   NUMERIC(20,4) NOT NULL DEFAULT 0,
    total       NUMERIC(20,4) NOT NULL DEFAULT 0,
    valid_days  INTEGER NOT NULL DEFAULT 15,
    valid_until TIMESTAMPTZ,
    converted_invoice_id UUID REFERENCES invoices(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON quotations(tenant_id);
CREATE INDEX ON quotations(customer_id);
ALTER TABLE quotations ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON quotations USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

CREATE TABLE quotation_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    quotation_id    UUID NOT NULL REFERENCES quotations(id) ON DELETE CASCADE,
    product_id      UUID REFERENCES products(id) ON DELETE SET NULL,
    description     TEXT NOT NULL,
    qty             NUMERIC(20,4) NOT NULL,
    unit_price      NUMERIC(20,4) NOT NULL,
    discount_pct    NUMERIC(5,2) NOT NULL DEFAULT 0,
    line_total      NUMERIC(20,4) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON quotation_items(tenant_id);
CREATE INDEX ON quotation_items(quotation_id);
ALTER TABLE quotation_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON quotation_items USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Comisiones de vendedores
CREATE TABLE sales_commissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vendor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    invoice_id  UUID REFERENCES invoices(id) ON DELETE SET NULL,
    category_id UUID REFERENCES product_categories(id) ON DELETE SET NULL,
    rate_pct    NUMERIC(5,2) NOT NULL,
    base_amount NUMERIC(20,4) NOT NULL,  -- venta cobrada
    commission_amount NUMERIC(20,4) NOT NULL,
    period      TEXT NOT NULL,            -- YYYY-MM
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON sales_commissions(tenant_id);
CREATE INDEX ON sales_commissions(vendor_user_id);
ALTER TABLE sales_commissions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sales_commissions USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Rutas de reparto
CREATE TABLE delivery_routes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    driver_id   UUID REFERENCES users(id) ON DELETE SET NULL,
    date        DATE NOT NULL,
    status      TEXT NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON delivery_routes(tenant_id);
ALTER TABLE delivery_routes ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON delivery_routes USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

CREATE TABLE delivery_stops (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    route_id        UUID NOT NULL REFERENCES delivery_routes(id) ON DELETE CASCADE,
    customer_id     UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    invoice_id      UUID REFERENCES invoices(id) ON DELETE SET NULL,
    stop_order      INTEGER NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'delivered', 'rejected', 'skipped')),
    signature_data  TEXT,  -- firma digital base64
    delivered_at    TIMESTAMPTZ,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON delivery_stops(tenant_id);
CREATE INDEX ON delivery_stops(route_id);
ALTER TABLE delivery_stops ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON delivery_stops USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS delivery_stops CASCADE;
DROP TABLE IF EXISTS delivery_routes CASCADE;
DROP TABLE IF EXISTS sales_commissions CASCADE;
DROP TABLE IF EXISTS quotation_items CASCADE;
DROP TABLE IF EXISTS quotations CASCADE;
DROP TABLE IF EXISTS purchase_order_receipts CASCADE;
DROP TABLE IF EXISTS purchase_order_items CASCADE;
DROP TABLE IF EXISTS purchase_orders CASCADE;
-- +goose StatementEnd
