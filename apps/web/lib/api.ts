const BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:3001';

type TokenStore = {
  access: string;
  refresh: string;
  tenantId?: string;
};

// ponytail: sessionStorage for dev; swap to httpOnly cookie in prod
function getTokens(): TokenStore | null {
  if (typeof window === 'undefined') return null;
  const raw = sessionStorage.getItem('chiguire_tokens');
  return raw ? JSON.parse(raw) : null;
}

function setTokens(t: TokenStore) {
  sessionStorage.setItem('chiguire_tokens', JSON.stringify(t));
}

function clearTokens() {
  sessionStorage.removeItem('chiguire_tokens');
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const tokens = getTokens();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init.headers as Record<string, string>),
  };
  if (tokens?.access) headers['Authorization'] = `Bearer ${tokens.access}`;

  const res = await fetch(`${BASE}${path}`, { ...init, headers });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`${res.status}: ${text}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

export const api = {
  auth: {
    async register(email: string, password: string, fullName: string) {
      return request<{ user_id: string }>('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password, full_name: fullName }),
      });
    },
    async login(email: string, password: string, tenantId?: string) {
      const data = await request<{ access_token: string; refresh_token: string }>(
        '/auth/login',
        {
          method: 'POST',
          body: JSON.stringify({ email, password, tenant_id: tenantId }),
        },
      );
      setTokens({ access: data.access_token, refresh: data.refresh_token, tenantId });
      return data;
    },
    logout() {
      clearTokens();
    },
    getAccessToken() {
      return getTokens()?.access ?? null;
    },
  },
  tenants: {
    list() {
      return request<{ id: string; name: string; slug: string; plan: string }[]>('/tenants');
    },
    create(name: string, slug: string) {
      return request<{ id: string; name: string; slug: string }>('/tenants', {
        method: 'POST',
        body: JSON.stringify({ name, slug }),
      });
    },
  },
  powersync: {
    getToken() {
      return request<{ token: string; expires_at: number }>('/powersync/token');
    },
  },
  testItems: {
    list() {
      return request<{ id: string; label: string; created_at: string }[]>('/test-items');
    },
    create(label: string) {
      return request<{ id: string; label: string; created_at: string }>('/test-items', {
        method: 'POST',
        body: JSON.stringify({ label }),
      });
    },
  },
  customers: {
    list() {
      return request<Customer[]>('/customers');
    },
    create(data: Partial<Customer>) {
      return request<Customer>('/customers', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<Customer>(`/customers/${id}`);
    },
    update(id: string, data: Partial<Customer>) {
      return request<Customer>(`/customers/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
    },
    delete(id: string) {
      return request<void>(`/customers/${id}`, { method: 'DELETE' });
    },
  },
  vendors: {
    list() {
      return request<Vendor[]>('/vendors');
    },
    create(data: Partial<Vendor>) {
      return request<Vendor>('/vendors', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<Vendor>(`/vendors/${id}`);
    },
    update(id: string, data: Partial<Vendor>) {
      return request<Vendor>(`/vendors/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
    },
    delete(id: string) {
      return request<void>(`/vendors/${id}`, { method: 'DELETE' });
    },
  },
  products: {
    list() {
      return request<Product[]>('/products');
    },
    create(data: Partial<Product>) {
      return request<Product>('/products', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<Product>(`/products/${id}`);
    },
    update(id: string, data: Partial<Product>) {
      return request<Product>(`/products/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
    },
    delete(id: string) {
      return request<void>(`/products/${id}`, { method: 'DELETE' });
    },
  },
  categories: {
    list() {
      return request<Category[]>('/product-categories');
    },
    create(data: Partial<Category>) {
      return request<Category>('/product-categories', { method: 'POST', body: JSON.stringify(data) });
    },
  },
  units: {
    list() {
      return request<Unit[]>('/units');
    },
    create(data: Partial<Unit>) {
      return request<Unit>('/units', { method: 'POST', body: JSON.stringify(data) });
    },
  },
  productPrices: {
    list() {
      return request<ProductPrice[]>('/product-prices');
    },
    create(data: Partial<ProductPrice>) {
      return request<ProductPrice>('/product-prices', { method: 'POST', body: JSON.stringify(data) });
    },
  },
  branches: {
    list() {
      return request<Branch[]>('/branches');
    },
    create(data: Partial<Branch>) {
      return request<Branch>('/branches', { method: 'POST', body: JSON.stringify(data) });
    },
  },
  warehouses: {
    list() {
      return request<Warehouse[]>('/warehouses');
    },
    create(data: Partial<Warehouse>) {
      return request<Warehouse>('/warehouses', { method: 'POST', body: JSON.stringify(data) });
    },
  },
  stock: {
    listMovements() {
      return request<StockMovement[]>('/stock/movements');
    },
    createMovement(data: Partial<StockMovement>) {
      return request<StockMovement>('/stock/movements', { method: 'POST', body: JSON.stringify(data) });
    },
    getStock(productId: string) {
      return request<{ product_id: string; quantity: number }[]>(`/stock/${productId}`);
    },
  },
  invoices: {
    list() {
      return request<Invoice[]>('/invoices');
    },
    create(data: Partial<Invoice>) {
      return request<Invoice>('/invoices', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<Invoice>(`/invoices/${id}`);
    },
    update(id: string, data: Partial<Invoice>) {
      return request<Invoice>(`/invoices/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
    },
    void(id: string) {
      return request<void>(`/invoices/${id}/void`, { method: 'POST' });
    },
    emit(id: string) {
      return request<Invoice>(`/invoices/${id}/emit`, { method: 'POST' });
    },
    addPayment(id: string, data: Partial<InvoicePayment>) {
      return request<InvoicePayment>(`/invoices/${id}/payments`, { method: 'POST', body: JSON.stringify(data) });
    },
    addItems(id: string, items: Partial<InvoiceItem>[]) {
      return request<InvoiceItem[]>(`/invoices/${id}/items`, { method: 'POST', body: JSON.stringify({ items }) });
    },
  },
  fiscal: {
    listTaxCategories() {
      return request<TaxCategory[]>('/fiscal/tax-categories');
    },
    createWithholding(data: Partial<Withholding>) {
      return request<Withholding>('/fiscal/withholdings', { method: 'POST', body: JSON.stringify(data) });
    },
    listWithholdings() {
      return request<Withholding[]>('/fiscal/withholdings');
    },
    getExchangeRate(currency: string) {
      return request<{ currency: string; rate: number; date: string }>(`/fiscal/exchange-rate/${currency}`);
    },
    setExchangeRate(currency: string, rate: number) {
      return request<{ currency: string; rate: number }>(`/fiscal/exchange-rate/${currency}`, {
        method: 'POST',
        body: JSON.stringify({ rate }),
      });
    },
    generateFiscalBook(period: string) {
      return request<{ id: string }>(`/fiscal/books/${period}`, { method: 'POST' });
    },
    getFiscalBook(id: string) {
      return request<FiscalBook>(`/fiscal/books/${id}`);
    },
  },
  payments: {
    createLink(data: Partial<PaymentLink>) {
      return request<PaymentLink>('/payment-links', { method: 'POST', body: JSON.stringify(data) });
    },
    getLink(id: string) {
      return request<PaymentLink>(`/payment-links/${id}`);
    },
    listLinks() {
      return request<PaymentLink[]>('/payment-links');
    },
  },
  purchases: {
    list() {
      return request<PurchaseOrder[]>('/purchase-orders');
    },
    create(data: Partial<PurchaseOrder>) {
      return request<PurchaseOrder>('/purchase-orders', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<PurchaseOrder>(`/purchase-orders/${id}`);
    },
    approve(id: string) {
      return request<PurchaseOrder>(`/purchase-orders/${id}/approve`, { method: 'POST' });
    },
    receiveItem(id: string, itemId: string, data: { quantity: number }) {
      return request<PurchaseOrderItem>(`/purchase-orders/${id}/items/${itemId}/receive`, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },
    addPOItem(id: string, data: Partial<PurchaseOrderItem>) {
      return request<PurchaseOrderItem>(`/purchase-orders/${id}/items`, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },
  },
  quotations: {
    list() {
      return request<Quotation[]>('/quotations');
    },
    create(data: Partial<Quotation>) {
      return request<Quotation>('/quotations', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<Quotation>(`/quotations/${id}`);
    },
    convertToInvoice(id: string) {
      return request<Invoice>(`/quotations/${id}/convert`, { method: 'POST' });
    },
    send(id: string, data: { to: string }) {
      return request<void>(`/quotations/${id}/send`, { method: 'POST', body: JSON.stringify(data) });
    },
  },
  commissions: {
    list() {
      return request<Commission[]>('/commissions');
    },
    calculate(data: { from: string; to: string }) {
      return request<Commission[]>('/commissions/calculate', { method: 'POST', body: JSON.stringify(data) });
    },
  },
  delivery: {
    listRoutes() {
      return request<DeliveryRoute[]>('/delivery/routes');
    },
    createRoute(data: Partial<DeliveryRoute>) {
      return request<DeliveryRoute>('/delivery/routes', { method: 'POST', body: JSON.stringify(data) });
    },
    getRoute(id: string) {
      return request<DeliveryRoute>(`/delivery/routes/${id}`);
    },
    addStop(routeId: string, data: Partial<DeliveryStop>) {
      return request<DeliveryStop>(`/delivery/routes/${routeId}/stops`, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },
    updateStopStatus(routeId: string, stopId: string, status: string) {
      return request<DeliveryStop>(`/delivery/routes/${routeId}/stops/${stopId}/status`, {
        method: 'POST',
        body: JSON.stringify({ status }),
      });
    },
  },
  saas: {
    getSubscription() {
      return request<Subscription>('/saas/subscription');
    },
    updatePlan(plan: string) {
      return request<Subscription>('/saas/subscription', { method: 'PATCH', body: JSON.stringify({ plan }) });
    },
    getPlanLimits() {
      return request<Record<string, number>>('/saas/plan-limits');
    },
  },
  webhooks: {
    list() {
      return request<OutgoingWebhook[]>('/webhooks');
    },
    create(data: Partial<OutgoingWebhook>) {
      return request<OutgoingWebhook>('/webhooks', { method: 'POST', body: JSON.stringify(data) });
    },
    delete(id: string) {
      return request<void>(`/webhooks/${id}`, { method: 'DELETE' });
    },
    test(id: string) {
      return request<{ success: boolean; status: number }>(`/webhooks/${id}/test`, { method: 'POST' });
    },
    listDeliveries() {
      return request<WebhookDelivery[]>('/webhooks/deliveries');
    },
  },
  sellers: {
    list() {
      return request<Seller[]>('/sellers');
    },
    create(data: Partial<Seller>) {
      return request<Seller>('/sellers', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<Seller>(`/sellers/${id}`);
    },
    update(id: string, data: Partial<Seller>) {
      return request<Seller>(`/sellers/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
    },
    delete(id: string) {
      return request<void>(`/sellers/${id}`, { method: 'DELETE' });
    },
    listCommissions(sellerId: string) {
      return request<SellerCommission[]>(`/sellers/${sellerId}/commissions`);
    },
    markCommissionsPaid(sellerId: string, data: { commission_ids: string[]; reference: string; notes?: string }) {
      return request<void>(`/sellers/${sellerId}/commissions/pay`, { method: 'POST', body: JSON.stringify(data) });
    },
  },
  paymentMethods: {
    list() {
      return request<PaymentMethod[]>('/payment-methods');
    },
    create(data: Partial<PaymentMethod>) {
      return request<PaymentMethod>('/payment-methods', { method: 'POST', body: JSON.stringify(data) });
    },
    update(id: string, data: Partial<PaymentMethod>) {
      return request<PaymentMethod>(`/payment-methods/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
    },
    delete(id: string) {
      return request<void>(`/payment-methods/${id}`, { method: 'DELETE' });
    },
  },
  fiscalDevices: {
    list() {
      return request<FiscalDevice[]>('/fiscal-devices');
    },
    create(data: Partial<FiscalDevice>) {
      return request<FiscalDevice>('/fiscal-devices', { method: 'POST', body: JSON.stringify(data) });
    },
    update(id: string, data: Partial<FiscalDevice>) {
      return request<FiscalDevice>(`/fiscal-devices/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
    },
    delete(id: string) {
      return request<void>(`/fiscal-devices/${id}`, { method: 'DELETE' });
    },
    listSequences(deviceId: string) {
      return request<DocumentSequence[]>(`/fiscal-devices/${deviceId}/sequences`);
    },
    createSequence(deviceId: string, data: Partial<DocumentSequence>) {
      return request<DocumentSequence>(`/fiscal-devices/${deviceId}/sequences`, { method: 'POST', body: JSON.stringify(data) });
    },
    listContingency(deviceId: string) {
      return request<ContingencyBook[]>(`/fiscal-devices/${deviceId}/contingency`);
    },
    createContingency(deviceId: string, data: Partial<ContingencyBook>) {
      return request<ContingencyBook>(`/fiscal-devices/${deviceId}/contingency`, { method: 'POST', body: JSON.stringify(data) });
    },
  },
  creditNotes: {
    list() {
      return request<CreditNote[]>('/credit-notes');
    },
    create(data: Partial<CreditNote>) {
      return request<CreditNote>('/credit-notes', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<CreditNote>(`/credit-notes/${id}`);
    },
    void(id: string) {
      return request<void>(`/credit-notes/${id}/void`, { method: 'POST' });
    },
  },
  transfers: {
    list() {
      return request<Transfer[]>('/transfers');
    },
    create(data: Partial<Transfer>) {
      return request<Transfer>('/transfers', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<Transfer>(`/transfers/${id}`);
    },
    ship(id: string) {
      return request<Transfer>(`/transfers/${id}/ship`, { method: 'POST' });
    },
    receive(id: string) {
      return request<Transfer>(`/transfers/${id}/receive`, { method: 'POST' });
    },
  },
  manufacturing: {
    list() {
      return request<ManufacturingOrder[]>('/manufacturing-orders');
    },
    create(data: Partial<ManufacturingOrder>) {
      return request<ManufacturingOrder>('/manufacturing-orders', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<ManufacturingOrder>(`/manufacturing-orders/${id}`);
    },
    start(id: string) {
      return request<ManufacturingOrder>(`/manufacturing-orders/${id}/start`, { method: 'POST' });
    },
    complete(id: string) {
      return request<ManufacturingOrder>(`/manufacturing-orders/${id}/complete`, { method: 'POST' });
    },
  },
  picking: {
    list() {
      return request<PickingList[]>('/picking-lists');
    },
    create(data: Partial<PickingList>) {
      return request<PickingList>('/picking-lists', { method: 'POST', body: JSON.stringify(data) });
    },
    get(id: string) {
      return request<PickingList>(`/picking-lists/${id}`);
    },
    verifyItem(id: string, itemId: string, data: { qty_picked: number }) {
      return request<PickingItem>(`/picking-lists/${id}/items/${itemId}/verify`, { method: 'POST', body: JSON.stringify(data) });
    },
    complete(id: string) {
      return request<PickingList>(`/picking-lists/${id}/complete`, { method: 'POST' });
    },
  },
  accountsPayable: {
    listReceivable() {
      return request<AccountItem[]>('/accounts/receivable');
    },
    listPayable() {
      return request<AccountItem[]>('/accounts/payable');
    },
    createPayment(data: Partial<AccountPayment>) {
      return request<AccountPayment>('/accounts/payments', { method: 'POST', body: JSON.stringify(data) });
    },
    listPayments() {
      return request<AccountPayment[]>('/accounts/payments');
    },
    sendReminder(id: string, data: { channel: 'sms' | 'email' }) {
      return request<void>(`/accounts/${id}/reminder`, { method: 'POST', body: JSON.stringify(data) });
    },
  },
  reports: {
    salesBook(from: string, to: string) {
      return request<Record<string, unknown>[]>(`/reports/sales-book?from=${from}&to=${to}`);
    },
    purchasesBook(from: string, to: string) {
      return request<Record<string, unknown>[]>(`/reports/purchases-book?from=${from}&to=${to}`);
    },
    inventoryCurrent() {
      return request<Record<string, unknown>[]>('/reports/inventory-current');
    },
    inventoryValued() {
      return request<Record<string, unknown>[]>('/reports/inventory-valued');
    },
    kardex(productId: string, from: string, to: string) {
      return request<Record<string, unknown>[]>(`/reports/kardex?product_id=${productId}&from=${from}&to=${to}`);
    },
    art177(from: string, to: string) {
      return request<Record<string, unknown>[]>(`/reports/art-177?from=${from}&to=${to}`);
    },
    igtfReport(from: string, to: string) {
      return request<Record<string, unknown>[]>(`/reports/igtf?from=${from}&to=${to}`);
    },
  },
  apiTokens: {
    list() {
      return request<ApiToken[]>('/api-tokens');
    },
    create(data: { name: string; expires_at?: string }) {
      return request<{ id: string; token: string; name: string; expires_at: string }>('/api-tokens', { method: 'POST', body: JSON.stringify(data) });
    },
    revoke(id: string) {
      return request<void>(`/api-tokens/${id}`, { method: 'DELETE' });
    },
  },
  importExport: {
    async importCustomers(file: File) {
      const form = new FormData();
      form.append('file', file);
      const tokens = getTokens();
      const headers: Record<string, string> = {};
      if (tokens?.access) headers['Authorization'] = `Bearer ${tokens.access}`;
      const res = await fetch(`${BASE}/import/customers`, { method: 'POST', headers, body: form });
      if (!res.ok) throw new Error(`${res.status}: ${await res.text()}`);
      return res.json() as Promise<{ imported: number }>;
    },
    async importProducts(file: File) {
      const form = new FormData();
      form.append('file', file);
      const tokens = getTokens();
      const headers: Record<string, string> = {};
      if (tokens?.access) headers['Authorization'] = `Bearer ${tokens.access}`;
      const res = await fetch(`${BASE}/import/products`, { method: 'POST', headers, body: form });
      if (!res.ok) throw new Error(`${res.status}: ${await res.text()}`);
      return res.json() as Promise<{ imported: number }>;
    },
    async exportInvoices() {
      const tokens = getTokens();
      const headers: Record<string, string> = {};
      if (tokens?.access) headers['Authorization'] = `Bearer ${tokens.access}`;
      const res = await fetch(`${BASE}/export/invoices`, { headers });
      if (!res.ok) throw new Error(`${res.status}: ${await res.text()}`);
      return res.blob();
    },
    async exportCustomers() {
      const tokens = getTokens();
      const headers: Record<string, string> = {};
      if (tokens?.access) headers['Authorization'] = `Bearer ${tokens.access}`;
      const res = await fetch(`${BASE}/export/customers`, { headers });
      if (!res.ok) throw new Error(`${res.status}: ${await res.text()}`);
      return res.blob();
    },
  },
};

export type Customer = {
  id: string; tenant_id: string; name: string; tax_id: string; email: string; phone: string;
  address: string; created_at: string;
};
export type Vendor = {
  id: string; tenant_id: string; name: string; tax_id: string; email: string; phone: string;
  contact_name: string; created_at: string;
};
export type Product = {
  id: string; tenant_id: string; sku: string; name: string; description: string;
  category_id: string; unit_id: string; cost: number; is_active: boolean; created_at: string;
};
export type Category = { id: string; tenant_id: string; name: string; created_at: string };
export type Unit = { id: string; tenant_id: string; name: string; symbol: string; created_at: string };
export type ProductPrice = {
  id: string; tenant_id: string; product_id: string; price_list: string;
  currency: string; price: number; created_at: string;
};
export type Branch = {
  id: string; tenant_id: string; name: string; address: string; phone: string;
  is_active: boolean; created_at: string;
};
export type Warehouse = {
  id: string; tenant_id: string; branch_id: string; name: string; code: string;
  created_at: string;
};
export type StockMovement = {
  id: string; tenant_id: string; product_id: string; warehouse_id: string;
  movement_type: string; quantity: number; reference: string; created_at: string;
};
export type Invoice = {
  id: string; tenant_id: string; customer_id: string; number: string; status: string;
  currency: string; subtotal: number; tax: number; total: number; exchange_rate: number;
  issue_date: string; due_date: string; created_at: string;
};
export type InvoiceItem = {
  id: string; invoice_id: string; product_id: string; description: string;
  quantity: number; unit_price: number; tax_rate: number; line_total: number;
};
export type InvoicePayment = {
  id: string; invoice_id: string; amount: number; currency: string;
  payment_method: string; reference: string; created_at: string;
};
export type TaxCategory = {
  id: string; tenant_id: string; name: string; rate: number; is_retention: boolean;
  created_at: string;
};
export type Withholding = {
  id: string; tenant_id: string; invoice_id: string; type: string; base_amount: number;
  rate: number; amount: number; created_at: string;
};
export type FiscalBook = {
  id: string; tenant_id: string; period: string; book_type: string; status: string;
  generated_at: string;
};
export type PaymentLink = {
  id: string; tenant_id: string; invoice_id: string; provider: string; url: string;
  amount: number; currency: string; status: string; created_at: string;
};
export type PurchaseOrder = {
  id: string; tenant_id: string; vendor_id: string; number: string; status: string;
  total: number; currency: string; created_at: string;
};
export type PurchaseOrderItem = {
  id: string; purchase_order_id: string; product_id: string; quantity: number;
  unit_cost: number; received_quantity: number; line_total: number;
};
export type Quotation = {
  id: string; tenant_id: string; customer_id: string; number: string; status: string;
  total: number; currency: string; valid_until: string; created_at: string;
};
export type Commission = {
  id: string; tenant_id: string; salesperson_id: string; invoice_id: string;
  rate: number; amount: number; status: string; created_at: string;
};
export type DeliveryRoute = {
  id: string; tenant_id: string; driver_name: string; date: string; status: string;
  created_at: string;
};
export type DeliveryStop = {
  id: string; route_id: string; sequence: number; customer_id: string; address: string;
  status: string; created_at: string;
};
export type OutgoingWebhook = {
  id: string; tenant_id: string; url: string; event: string; secret: string;
  is_active: boolean; created_at: string;
};
export type WebhookDelivery = {
  id: string; webhook_id: string; event: string; status: number; attempt: number;
  created_at: string;
};
export type Subscription = {
  id: string; tenant_id: string; plan: string; status: string; seats: number;
  current_period_end: string; created_at: string;
};
export type Seller = {
  id: string; tenant_id: string; name: string; email: string; phone: string;
  commission_pct: number; is_active: boolean; created_at: string;
};
export type SellerCommission = {
  id: string; seller_id: string; invoice_id: string; doc_type: string;
  base_amount: number; commission_amount: number; rate: number;
  payment_status: string; cxc_customer: string; created_at: string;
};
export type PaymentMethod = {
  id: string; tenant_id: string; name: string; type: string; currency: string;
  allow_invoices: boolean; allow_change: boolean; allow_refunds: boolean;
  is_active: boolean; created_at: string;
};
export type FiscalDevice = {
  id: string; tenant_id: string; name: string; type: string; model: string;
  serial: string; branch_id: string; status: string; created_at: string;
};
export type DocumentSequence = {
  id: string; device_id: string; doc_type: string; prefix: string;
  suffix: string; last_seq: number; created_at: string;
};
export type ContingencyBook = {
  id: string; device_id: string; doc_type: string; prefix: string;
  start_seq: number; end_seq: number; current_seq: number; created_at: string;
};
export type CreditNote = {
  id: string; tenant_id: string; number: string; type: string;
  invoice_id: string; reason: string; total: number; status: string;
  created_at: string;
};
export type Transfer = {
  id: string; tenant_id: string; from_warehouse_id: string; to_warehouse_id: string;
  status: string; created_at: string;
};
export type ManufacturingOrder = {
  id: string; tenant_id: string; product_id: string; warehouse_id: string;
  quantity: number; status: string; started_at: string | null;
  completed_at: string | null; created_at: string;
};
export type PickingList = {
  id: string; tenant_id: string; warehouse_id: string; route_id: string | null;
  status: string; created_at: string;
};
export type PickingItem = {
  id: string; picking_list_id: string; product_id: string; qty_required: number;
  qty_picked: number; verified: boolean;
};
export type AccountItem = {
  id: string; tenant_id: string; party_name: string; amount: number;
  balance: number; status: string; due_date: string; type: string; created_at: string;
};
export type AccountPayment = {
  id: string; tenant_id: string; account_id: string; amount: number;
  payment_method: string; reference: string; notes: string; created_at: string;
};
export type ApiToken = {
  id: string; tenant_id: string; name: string; token_preview: string;
  last_used_at: string | null; expires_at: string; is_active: boolean;
  created_at: string;
};
