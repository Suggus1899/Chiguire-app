-- 00009_sellers_payment_methods.sql
-- Phase: Cachicamo parity — sellers, payment methods, fiscal devices

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sellers (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     uuid REFERENCES users(id) ON DELETE SET NULL,
    name        text NOT NULL,
    email       text,
    phone       text,
    commission_pct numeric(5,2) NOT NULL DEFAULT 0,
    is_active   boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_attribute_groups (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_attributes (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id    uuid NOT NULL REFERENCES product_attribute_groups(id) ON DELETE CASCADE,
    name        text NOT NULL,
    value       text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS payment_methods (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        text NOT NULL,
    type        text NOT NULL DEFAULT 'cash',
    currency    text NOT NULL DEFAULT 'USD',
    available_for_invoices   boolean NOT NULL DEFAULT true,
    available_for_change     boolean NOT NULL DEFAULT false,
    available_for_refunds    boolean NOT NULL DEFAULT false,
    is_active   boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS fiscal_devices (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id   uuid REFERENCES branches(id) ON DELETE SET NULL,
    name        text NOT NULL,
    device_type text NOT NULL,
    model       text,
    serial      text,
    is_active   boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS document_sequences (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    fiscal_device_id uuid REFERENCES fiscal_devices(id) ON DELETE CASCADE,
    doc_type    text NOT NULL,
    prefix      text NOT NULL DEFAULT '',
    suffix      text NOT NULL DEFAULT '',
    last_seq    integer NOT NULL DEFAULT 0,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS contingency_books (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id   uuid REFERENCES branches(id) ON DELETE SET NULL,
    doc_type    text NOT NULL,
    prefix      text NOT NULL,
    start_seq   integer NOT NULL,
    end_seq     integer NOT NULL,
    current_seq integer NOT NULL DEFAULT 0,
    is_active   boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS credit_notes (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id  uuid REFERENCES invoices(id) ON DELETE CASCADE,
    number      text,
    type        text NOT NULL DEFAULT 'credit',
    reason      text,
    subtotal    numeric(14,2) NOT NULL DEFAULT 0,
    tax_total   numeric(14,2) NOT NULL DEFAULT 0,
    total       numeric(14,2) NOT NULL DEFAULT 0,
    status      text NOT NULL DEFAULT 'draft',
    issued_at   timestamptz,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS credit_note_items (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    credit_note_id  uuid NOT NULL REFERENCES credit_notes(id) ON DELETE CASCADE,
    product_id      uuid REFERENCES products(id) ON DELETE SET NULL,
    description     text NOT NULL,
    qty             numeric(14,2) NOT NULL DEFAULT 1,
    unit_price      numeric(14,2) NOT NULL DEFAULT 0,
    tax_pct         numeric(5,2) NOT NULL DEFAULT 0,
    line_total      numeric(14,2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS inventory_transfers (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    from_warehouse_id uuid NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    to_warehouse_id   uuid NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    status          text NOT NULL DEFAULT 'pending',
    notes           text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS inventory_transfer_items (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    transfer_id uuid NOT NULL REFERENCES inventory_transfers(id) ON DELETE CASCADE,
    product_id  uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    qty         numeric(14,2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS manufacturing_orders (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    product_id  uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    warehouse_id uuid NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    qty_to_produce numeric(14,2) NOT NULL DEFAULT 1,
    status      text NOT NULL DEFAULT 'pending',
    started_at  timestamptz,
    completed_at timestamptz,
    notes       text,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS manufacturing_bom (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    order_id    uuid NOT NULL REFERENCES manufacturing_orders(id) ON DELETE CASCADE,
    component_product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    qty_per_unit numeric(14,4) NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS picking_lists (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    warehouse_id uuid NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    route_id    uuid REFERENCES delivery_routes(id) ON DELETE SET NULL,
    status      text NOT NULL DEFAULT 'pending',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS picking_items (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    picking_list_id uuid NOT NULL REFERENCES picking_lists(id) ON DELETE CASCADE,
    product_id  uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    qty_required numeric(14,2) NOT NULL DEFAULT 0,
    qty_picked  numeric(14,2) NOT NULL DEFAULT 0,
    is_verified boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS accounts_receivable (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id uuid NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    invoice_id  uuid REFERENCES invoices(id) ON DELETE SET NULL,
    amount      numeric(14,2) NOT NULL DEFAULT 0,
    balance     numeric(14,2) NOT NULL DEFAULT 0,
    status      text NOT NULL DEFAULT 'pending',
    due_date    date,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts_payable (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vendor_id   uuid NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    purchase_order_id uuid REFERENCES purchase_orders(id) ON DELETE SET NULL,
    amount      numeric(14,2) NOT NULL DEFAULT 0,
    balance     numeric(14,2) NOT NULL DEFAULT 0,
    status      text NOT NULL DEFAULT 'pending',
    due_date    date,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS account_payments (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    account_type text NOT NULL,
    account_id  uuid NOT NULL,
    amount      numeric(14,2) NOT NULL DEFAULT 0,
    payment_method_id uuid REFERENCES payment_methods(id) ON DELETE SET NULL,
    reference   text,
    notes       text,
    paid_at     timestamptz NOT NULL DEFAULT now(),
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS api_tokens (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        text NOT NULL,
    token_hash  text NOT NULL,
    last_used_at timestamptz,
    expires_at  timestamptz,
    is_active   boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now()
);

-- RLS policies
ALTER TABLE sellers ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_attribute_groups ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_attributes ENABLE ROW LEVEL SECURITY;
ALTER TABLE payment_methods ENABLE ROW LEVEL SECURITY;
ALTER TABLE fiscal_devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE document_sequences ENABLE ROW LEVEL SECURITY;
ALTER TABLE contingency_books ENABLE ROW LEVEL SECURITY;
ALTER TABLE credit_notes ENABLE ROW LEVEL SECURITY;
ALTER TABLE credit_note_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory_transfers ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory_transfer_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE manufacturing_orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE manufacturing_bom ENABLE ROW LEVEL SECURITY;
ALTER TABLE picking_lists ENABLE ROW LEVEL SECURITY;
ALTER TABLE picking_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE accounts_receivable ENABLE ROW LEVEL SECURITY;
ALTER TABLE accounts_payable ENABLE ROW LEVEL SECURITY;
ALTER TABLE account_payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE api_tokens ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_sellers ON sellers USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_pag ON product_attribute_groups USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_pm ON payment_methods USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_fd ON fiscal_devices USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_ds ON document_sequences USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_cb ON contingency_books USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_cn ON credit_notes USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_cni ON credit_note_items USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_it ON inventory_transfers USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_iti ON inventory_transfer_items USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_mo ON manufacturing_orders USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_mbom ON manufacturing_bom USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_pl ON picking_lists USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_pi ON picking_items USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_ar ON accounts_receivable USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_ap ON accounts_payable USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_apt ON account_payments USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
CREATE POLICY tenant_isolation_at ON api_tokens USING (tenant_id = current_setting('app.tenant_id', true)::uuid);

-- Indexes on tenant_id for all new tables
CREATE INDEX IF NOT EXISTS idx_sellers_tenant ON sellers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_pag_tenant ON product_attribute_groups(tenant_id);
CREATE INDEX IF NOT EXISTS idx_pm_tenant ON payment_methods(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fd_tenant ON fiscal_devices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ds_tenant ON document_sequences(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cb_tenant ON contingency_books(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cn_tenant ON credit_notes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cni_tenant ON credit_note_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_it_tenant ON inventory_transfers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_iti_tenant ON inventory_transfer_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_mo_tenant ON manufacturing_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_mbom_tenant ON manufacturing_bom(tenant_id);
CREATE INDEX IF NOT EXISTS idx_pl_tenant ON picking_lists(tenant_id);
CREATE INDEX IF NOT EXISTS idx_pi_tenant ON picking_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ar_tenant ON accounts_receivable(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ap_tenant ON accounts_payable(tenant_id);
CREATE INDEX IF NOT EXISTS idx_apt_tenant ON account_payments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_at_tenant ON api_tokens(tenant_id);

-- Add status column to sales_commissions for mark-as-paid functionality
ALTER TABLE sales_commissions ADD COLUMN IF NOT EXISTS status text NOT NULL DEFAULT 'unpaid';
ALTER TABLE sales_commissions ADD COLUMN IF NOT EXISTS paid_at timestamptz;
ALTER TABLE sales_commissions ADD COLUMN IF NOT EXISTS payment_reference text;
ALTER TABLE sales_commissions ADD COLUMN IF NOT EXISTS payment_notes text;

-- Add seller_id to customers (optional default seller)
ALTER TABLE customers ADD COLUMN IF NOT EXISTS default_seller_id uuid REFERENCES sellers(id) ON DELETE SET NULL;

-- Add seller_id to invoices (optional per-invoice seller)
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS seller_id uuid REFERENCES sellers(id) ON DELETE SET NULL;

-- Add seller_id to quotations
ALTER TABLE quotations ADD COLUMN IF NOT EXISTS seller_id uuid REFERENCES sellers(id) ON DELETE SET NULL;

-- set_updated_at function (idempotent — safe if already exists from other migrations)
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Updated triggers
DROP TRIGGER IF EXISTS set_updated_at_sellers ON sellers;
CREATE TRIGGER set_updated_at_sellers BEFORE UPDATE ON sellers FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_pm ON payment_methods;
CREATE TRIGGER set_updated_at_pm BEFORE UPDATE ON payment_methods FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_fd ON fiscal_devices;
CREATE TRIGGER set_updated_at_fd BEFORE UPDATE ON fiscal_devices FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_cn ON credit_notes;
CREATE TRIGGER set_updated_at_cn BEFORE UPDATE ON credit_notes FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_it ON inventory_transfers;
CREATE TRIGGER set_updated_at_it BEFORE UPDATE ON inventory_transfers FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_mo ON manufacturing_orders;
CREATE TRIGGER set_updated_at_mo BEFORE UPDATE ON manufacturing_orders FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_pl ON picking_lists;
CREATE TRIGGER set_updated_at_pl BEFORE UPDATE ON picking_lists FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_ar ON accounts_receivable;
CREATE TRIGGER set_updated_at_ar BEFORE UPDATE ON accounts_receivable FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS set_updated_at_ap ON accounts_payable;
CREATE TRIGGER set_updated_at_ap BEFORE UPDATE ON accounts_payable FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS api_tokens CASCADE;
DROP TABLE IF EXISTS account_payments CASCADE;
DROP TABLE IF EXISTS accounts_payable CASCADE;
DROP TABLE IF EXISTS accounts_receivable CASCADE;
DROP TABLE IF EXISTS picking_items CASCADE;
DROP TABLE IF EXISTS picking_lists CASCADE;
DROP TABLE IF EXISTS manufacturing_bom CASCADE;
DROP TABLE IF EXISTS manufacturing_orders CASCADE;
DROP TABLE IF EXISTS inventory_transfer_items CASCADE;
DROP TABLE IF EXISTS inventory_transfers CASCADE;
DROP TABLE IF EXISTS credit_note_items CASCADE;
DROP TABLE IF EXISTS credit_notes CASCADE;
DROP TABLE IF EXISTS contingency_books CASCADE;
DROP TABLE IF EXISTS document_sequences CASCADE;
DROP TABLE IF EXISTS fiscal_devices CASCADE;
DROP TABLE IF EXISTS payment_methods CASCADE;
DROP TABLE IF EXISTS product_attributes CASCADE;
DROP TABLE IF EXISTS product_attribute_groups CASCADE;
DROP TABLE IF EXISTS sellers CASCADE;

ALTER TABLE sales_commissions DROP COLUMN IF EXISTS status;
ALTER TABLE sales_commissions DROP COLUMN IF EXISTS paid_at;
ALTER TABLE sales_commissions DROP COLUMN IF EXISTS payment_reference;
ALTER TABLE sales_commissions DROP COLUMN IF EXISTS payment_notes;

ALTER TABLE customers DROP COLUMN IF EXISTS default_seller_id;
ALTER TABLE invoices DROP COLUMN IF EXISTS seller_id;
ALTER TABLE quotations DROP COLUMN IF EXISTS seller_id;
-- +goose StatementEnd
