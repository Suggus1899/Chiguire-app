import 'package:powersync/powersync.dart';

// Phase 0: minimal schema for sync proof-of-concept
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
]);

late PowerSyncDatabase db;

Future<void> openDatabase() async {
  db = PowerSyncDatabase(
    schema: schema,
    path: 'chiguire.db',
  );
  await db.initialize();
}
