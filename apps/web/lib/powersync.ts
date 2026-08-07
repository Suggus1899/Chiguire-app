import { PowerSyncDatabase, Schema, Table, column } from '@powersync/web';

const customers = new Table({
  tenant_id: column.text, name: column.text, tax_id: column.text, email: column.text,
  phone: column.text, address: column.text, created_at: column.text,
});
const vendors = new Table({
  tenant_id: column.text, name: column.text, tax_id: column.text, email: column.text,
  phone: column.text, contact_name: column.text, created_at: column.text,
});
const productCategories = new Table({
  tenant_id: column.text, name: column.text, created_at: column.text,
});
const unitsOfMeasure = new Table({
  tenant_id: column.text, name: column.text, symbol: column.text, created_at: column.text,
});
const products = new Table({
  tenant_id: column.text, sku: column.text, name: column.text, description: column.text,
  category_id: column.text, unit_id: column.text, cost: column.real, is_active: column.integer,
  created_at: column.text,
});
const productPrices = new Table({
  tenant_id: column.text, product_id: column.text, price_list: column.text,
  currency: column.text, price: column.real, created_at: column.text,
});
const branches = new Table({
  tenant_id: column.text, name: column.text, address: column.text, phone: column.text,
  is_active: column.integer, created_at: column.text,
});
const warehouses = new Table({
  tenant_id: column.text, branch_id: column.text, name: column.text, code: column.text,
  created_at: column.text,
});
const stockMovements = new Table({
  tenant_id: column.text, product_id: column.text, warehouse_id: column.text,
  movement_type: column.text, quantity: column.real, reference: column.text, created_at: column.text,
});
const invoices = new Table({
  tenant_id: column.text, customer_id: column.text, number: column.text, status: column.text,
  currency: column.text, subtotal: column.real, tax: column.real, total: column.real,
  exchange_rate: column.real, issue_date: column.text, due_date: column.text, created_at: column.text,
});
const invoiceItems = new Table({
  invoice_id: column.text, product_id: column.text, description: column.text,
  quantity: column.real, unit_price: column.real, tax_rate: column.real, line_total: column.real,
});
const invoicePayments = new Table({
  invoice_id: column.text, amount: column.real, currency: column.text,
  payment_method: column.text, reference: column.text, created_at: column.text,
});
const invoiceSequences = new Table({
  tenant_id: column.text, prefix: column.text, last_number: column.integer, created_at: column.text,
});
const taxCategories = new Table({
  tenant_id: column.text, name: column.text, rate: column.real, is_retention: column.integer,
  created_at: column.text,
});
const taxWithholdings = new Table({
  tenant_id: column.text, invoice_id: column.text, type: column.text, base_amount: column.real,
  rate: column.real, amount: column.real, created_at: column.text,
});
const fiscalPeriods = new Table({
  tenant_id: column.text, period: column.text, status: column.text, created_at: column.text,
});
const fiscalBooks = new Table({
  tenant_id: column.text, period: column.text, book_type: column.text, status: column.text,
  generated_at: column.text,
});
const paymentLinks = new Table({
  tenant_id: column.text, invoice_id: column.text, provider: column.text, url: column.text,
  amount: column.real, currency: column.text, status: column.text, created_at: column.text,
});
const paymentWebhookEvents = new Table({
  tenant_id: column.text, link_id: column.text, event_type: column.text, payload: column.text,
  created_at: column.text,
});
const purchaseOrders = new Table({
  tenant_id: column.text, vendor_id: column.text, number: column.text, status: column.text,
  total: column.real, currency: column.text, created_at: column.text,
});
const purchaseOrderItems = new Table({
  purchase_order_id: column.text, product_id: column.text, quantity: column.real,
  unit_cost: column.real, received_quantity: column.real, line_total: column.real,
});
const quotations = new Table({
  tenant_id: column.text, customer_id: column.text, number: column.text, status: column.text,
  total: column.real, currency: column.text, valid_until: column.text, created_at: column.text,
});
const quotationItems = new Table({
  quotation_id: column.text, product_id: column.text, description: column.text,
  quantity: column.real, unit_price: column.real, line_total: column.real,
});
const salesCommissions = new Table({
  tenant_id: column.text, salesperson_id: column.text, invoice_id: column.text,
  rate: column.real, amount: column.real, status: column.text, created_at: column.text,
});
const deliveryRoutes = new Table({
  tenant_id: column.text, driver_name: column.text, date: column.text, status: column.text,
  created_at: column.text,
});
const deliveryStops = new Table({
  route_id: column.text, sequence: column.integer, customer_id: column.text, address: column.text,
  status: column.text, created_at: column.text,
});
const outgoingWebhooks = new Table({
  tenant_id: column.text, url: column.text, event: column.text, secret: column.text,
  is_active: column.integer, created_at: column.text,
});
const webhookDeliveries = new Table({
  webhook_id: column.text, event: column.text, status: column.integer, attempt: column.integer,
  created_at: column.text,
});
const subscriptions = new Table({
  tenant_id: column.text, plan: column.text, status: column.text, seats: column.integer,
  current_period_end: column.text, created_at: column.text,
});

export const AppSchema = new Schema({
  test_items: new Table({
    tenant_id: column.text, label: column.text, created_at: column.text,
  }),
  customers,
  vendors,
  product_categories: productCategories,
  units_of_measure: unitsOfMeasure,
  products,
  product_prices: productPrices,
  branches,
  warehouses,
  stock_movements: stockMovements,
  invoices,
  invoice_items: invoiceItems,
  invoice_payments: invoicePayments,
  invoice_sequences: invoiceSequences,
  tax_categories: taxCategories,
  tax_withholdings: taxWithholdings,
  fiscal_periods: fiscalPeriods,
  fiscal_books: fiscalBooks,
  payment_links: paymentLinks,
  payment_webhook_events: paymentWebhookEvents,
  purchase_orders: purchaseOrders,
  purchase_order_items: purchaseOrderItems,
  quotations,
  quotation_items: quotationItems,
  sales_commissions: salesCommissions,
  delivery_routes: deliveryRoutes,
  delivery_stops: deliveryStops,
  outgoing_webhooks: outgoingWebhooks,
  webhook_deliveries: webhookDeliveries,
  subscriptions,
});

export type AppDatabase = (typeof AppSchema)['types'];

let _db: PowerSyncDatabase | null = null;

export function getPowerSync(): PowerSyncDatabase {
  if (!_db) {
    _db = new PowerSyncDatabase({
      schema: AppSchema,
      database: { dbFilename: 'chiguire.db' },
    });
  }
  return _db;
}
