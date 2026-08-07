import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:powersync/powersync.dart';
import 'api.dart';

/// PowerSync service endpoint. Configure via the POWERSYNC_URL dart-define.
const _powerSyncEndpoint = String.fromEnvironment(
  'POWERSYNC_URL',
  defaultValue: 'http://10.0.2.2:8080',
);

/// Backend connector that bridges the local PowerSync database with the
/// Chiguire API. It fetches sync credentials from the API and (for now)
/// treats local-change uploads as a no-op — reads are the priority.
class _ChiguireConnector extends PowerSyncBackendConnector {
  @override
  Future<PowerSyncCredentials?> fetchCredentials() async {
    final token = await powerSyncApi.getToken();
    return PowerSyncCredentials(
      endpoint: _powerSyncEndpoint,
      token: token,
    );
  }

  @override
  Future<void> uploadData(PowerSyncDatabase database) async {
    // No-op for now: reads are the priority. Local writes will be synced in a
    // later phase once the upload endpoint is implemented.
  }
}

const schema = Schema([
  Table('test_items', [
    Column.text('tenant_id'),
    Column.text('label'),
    Column.text('created_at'),
  ]),
  Table('exchange_rates', [
    Column.text('currency'),
    Column.real('rate_to_ves'),
    Column.text('source'),
    Column.text('effective_at'),
  ]),
  Table('customers', [
    Column.text('tenant_id'),
    Column.text('name'),
    Column.text('tax_id'),
    Column.text('email'),
    Column.text('phone'),
    Column.text('address'),
    Column.text('created_at'),
    Column.text('updated_at'),
  ]),
  Table('vendors', [
    Column.text('tenant_id'),
    Column.text('name'),
    Column.text('tax_id'),
    Column.text('email'),
    Column.text('phone'),
    Column.text('address'),
    Column.text('created_at'),
    Column.text('updated_at'),
  ]),
  Table('product_categories', [
    Column.text('tenant_id'),
    Column.text('name'),
    Column.text('parent_id'),
    Column.text('created_at'),
  ]),
  Table('units_of_measure', [
    Column.text('tenant_id'),
    Column.text('name'),
    Column.text('symbol'),
    Column.text('created_at'),
  ]),
  Table('products', [
    Column.text('tenant_id'),
    Column.text('sku'),
    Column.text('name'),
    Column.text('description'),
    Column.text('category_id'),
    Column.text('unit_id'),
    Column.real('cost_price'),
    Column.real('sale_price'),
    Column.text('created_at'),
    Column.text('updated_at'),
  ]),
  Table('product_prices', [
    Column.text('tenant_id'),
    Column.text('product_id'),
    Column.text('price_list'),
    Column.real('price'),
    Column.text('currency'),
    Column.text('effective_at'),
  ]),
  Table('branches', [
    Column.text('tenant_id'),
    Column.text('name'),
    Column.text('address'),
    Column.text('phone'),
    Column.text('created_at'),
  ]),
  Table('warehouses', [
    Column.text('tenant_id'),
    Column.text('branch_id'),
    Column.text('name'),
    Column.text('code'),
    Column.text('created_at'),
  ]),
  Table('stock_movements', [
    Column.text('tenant_id'),
    Column.text('warehouse_id'),
    Column.text('product_id'),
    Column.text('movement_type'),
    Column.real('quantity'),
    Column.text('reference'),
    Column.text('created_at'),
  ]),
  Table('invoices', [
    Column.text('tenant_id'),
    Column.text('customer_id'),
    Column.text('number'),
    Column.text('status'),
    Column.text('currency'),
    Column.real('subtotal'),
    Column.real('tax'),
    Column.real('total'),
    Column.text('issued_at'),
    Column.text('created_at'),
    Column.text('updated_at'),
  ]),
  Table('invoice_items', [
    Column.text('tenant_id'),
    Column.text('invoice_id'),
    Column.text('product_id'),
    Column.text('description'),
    Column.real('quantity'),
    Column.real('unit_price'),
    Column.real('line_total'),
  ]),
  Table('invoice_payments', [
    Column.text('tenant_id'),
    Column.text('invoice_id'),
    Column.text('method'),
    Column.real('amount'),
    Column.text('currency'),
    Column.text('reference'),
    Column.text('paid_at'),
  ]),
  Table('tax_categories', [
    Column.text('tenant_id'),
    Column.text('name'),
    Column.real('rate'),
    Column.text('created_at'),
  ]),
  Table('tax_withholdings', [
    Column.text('tenant_id'),
    Column.text('invoice_id'),
    Column.text('type'),
    Column.real('rate'),
    Column.real('amount'),
    Column.text('created_at'),
  ]),
  Table('fiscal_periods', [
    Column.text('tenant_id'),
    Column.text('period'),
    Column.text('status'),
    Column.text('created_at'),
  ]),
  Table('fiscal_books', [
    Column.text('tenant_id'),
    Column.text('period'),
    Column.text('book_type'),
    Column.text('file_url'),
    Column.text('generated_at'),
  ]),
  Table('payment_links', [
    Column.text('tenant_id'),
    Column.text('invoice_id'),
    Column.text('provider'),
    Column.text('link_url'),
    Column.text('status'),
    Column.real('amount'),
    Column.text('currency'),
    Column.text('created_at'),
  ]),
  Table('payment_webhook_events', [
    Column.text('tenant_id'),
    Column.text('payment_link_id'),
    Column.text('event_type'),
    Column.text('payload'),
    Column.text('received_at'),
  ]),
  Table('purchase_orders', [
    Column.text('tenant_id'),
    Column.text('vendor_id'),
    Column.text('number'),
    Column.text('status'),
    Column.real('total'),
    Column.text('currency'),
    Column.text('created_at'),
    Column.text('updated_at'),
  ]),
  Table('purchase_order_items', [
    Column.text('tenant_id'),
    Column.text('purchase_order_id'),
    Column.text('product_id'),
    Column.text('description'),
    Column.real('quantity'),
    Column.real('unit_price'),
    Column.real('line_total'),
  ]),
  Table('quotations', [
    Column.text('tenant_id'),
    Column.text('customer_id'),
    Column.text('number'),
    Column.text('status'),
    Column.real('total'),
    Column.text('currency'),
    Column.text('valid_until'),
    Column.text('created_at'),
  ]),
  Table('quotation_items', [
    Column.text('tenant_id'),
    Column.text('quotation_id'),
    Column.text('product_id'),
    Column.text('description'),
    Column.real('quantity'),
    Column.real('unit_price'),
    Column.real('line_total'),
  ]),
  Table('sales_commissions', [
    Column.text('tenant_id'),
    Column.text('salesperson_id'),
    Column.text('invoice_id'),
    Column.real('rate'),
    Column.real('amount'),
    Column.text('status'),
    Column.text('created_at'),
  ]),
  Table('delivery_routes', [
    Column.text('tenant_id'),
    Column.text('driver_name'),
    Column.text('date'),
    Column.text('status'),
    Column.text('created_at'),
  ]),
  Table('delivery_stops', [
    Column.text('tenant_id'),
    Column.text('route_id'),
    Column.text('customer_id'),
    Column.text('address'),
    Column.integer('sequence'),
    Column.text('status'),
    Column.text('completed_at'),
  ]),
  Table('outgoing_webhooks', [
    Column.text('tenant_id'),
    Column.text('event_type'),
    Column.text('payload'),
    Column.text('status'),
    Column.text('sent_at'),
  ]),
  Table('subscriptions', [
    Column.text('tenant_id'),
    Column.text('plan'),
    Column.text('status'),
    Column.text('billing_cycle'),
    Column.text('current_period_end'),
    Column.text('created_at'),
  ]),
]);

late PowerSyncDatabase db;

Future<void> openDatabase() async {
  db = PowerSyncDatabase(
    schema: schema,
    path: 'chiguire.db',
  );
  await db.initialize();
}

/// Connect the local PowerSync database to the PowerSync service so that
/// remote changes are replicated into SQLite. Call this after [openDatabase]
/// during app initialization.
///
/// The connection runs as a long-lived background sync loop and is
/// auto-reopened by PowerSync on failure, so this does not block app startup.
/// Any immediate setup errors are caught and logged so the app keeps working
/// offline against the local database.
Future<void> connect() async {
  try {
    unawaited(
      db.connect(connector: _ChiguireConnector()).catchError((Object e) {
        debugPrint('PowerSync connect failed: $e');
      }),
    );
  } catch (e) {
    debugPrint('PowerSync connect setup failed: $e');
  }
}
