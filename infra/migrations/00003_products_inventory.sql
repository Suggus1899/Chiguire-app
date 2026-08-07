-- +goose Up
-- +goose StatementBegin

-- Fase 1: Productos, catálogo e inventario

-- Categorías de productos
CREATE TABLE product_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    parent_id   UUID REFERENCES product_categories(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON product_categories(tenant_id);
ALTER TABLE product_categories ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON product_categories USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Unidades de medida
CREATE TABLE units_of_measure (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,        -- e.g. 'unidad', 'kg', 'litro'
    abbreviation TEXT NOT NULL,       -- e.g. 'und', 'kg', 'lt'
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON units_of_measure(tenant_id);
ALTER TABLE units_of_measure ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON units_of_measure USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Productos
CREATE TABLE products (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    sku             TEXT NOT NULL,
    barcode         TEXT,
    name            TEXT NOT NULL,
    description     TEXT,
    category_id     UUID REFERENCES product_categories(id) ON DELETE SET NULL,
    unit_id         UUID REFERENCES units_of_measure(id) ON DELETE SET NULL,
    tax_category    TEXT NOT NULL DEFAULT 'general' CHECK (tax_category IN ('general', 'reduced', 'exempt')),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    min_stock       NUMERIC(20,4) NOT NULL DEFAULT 0,
    avg_cost_usd    NUMERIC(20,4) NOT NULL DEFAULT 0, -- costo promedio ponderado
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON products(tenant_id);
CREATE UNIQUE INDEX ON products(tenant_id, sku);
ALTER TABLE products ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON products USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Precios por moneda
CREATE TABLE product_prices (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    currency    TEXT NOT NULL,           -- 'USD', 'VES'
    amount      NUMERIC(20,4) NOT NULL,
    valid_from  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON product_prices(tenant_id);
CREATE INDEX ON product_prices(product_id);
ALTER TABLE product_prices ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON product_prices USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Almacenes por sucursal
CREATE TABLE branches (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    address     TEXT,
    phone       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON branches(tenant_id);
ALTER TABLE branches ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON branches USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

CREATE TABLE warehouses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id   UUID REFERENCES branches(id) ON DELETE SET NULL,
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON warehouses(tenant_id);
ALTER TABLE warehouses ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON warehouses USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Movimientos de stock (append-only, nunca UPDATE)
CREATE TABLE stock_movements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    warehouse_id    UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    product_id      UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type            TEXT NOT NULL CHECK (type IN ('in', 'out', 'transfer', 'adjustment')),
    qty             NUMERIC(20,4) NOT NULL,  -- positivo para in, negativo para out
    reference_type  TEXT,                     -- 'invoice', 'purchase_order', 'manual'
    reference_id    UUID,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON stock_movements(tenant_id);
CREATE INDEX ON stock_movements(warehouse_id, product_id);
CREATE INDEX ON stock_movements(product_id);
ALTER TABLE stock_movements ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON stock_movements USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Vista materializada de stock por ubicación
CREATE MATERIALIZED VIEW stock_by_location AS
SELECT
    sm.tenant_id,
    sm.warehouse_id,
    sm.product_id,
    COALESCE(SUM(sm.qty), 0) AS qty_on_hand
FROM stock_movements sm
GROUP BY sm.tenant_id, sm.warehouse_id, sm.product_id
WITH DATA;

CREATE UNIQUE INDEX ON stock_by_location(tenant_id, warehouse_id, product_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS stock_by_location CASCADE;
DROP TABLE IF EXISTS stock_movements CASCADE;
DROP TABLE IF EXISTS warehouses CASCADE;
DROP TABLE IF EXISTS branches CASCADE;
DROP TABLE IF EXISTS product_prices CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS units_of_measure CASCADE;
DROP TABLE IF EXISTS product_categories CASCADE;
-- +goose StatementEnd
