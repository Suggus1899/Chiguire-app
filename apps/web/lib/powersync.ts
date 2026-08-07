import { PowerSyncDatabase, Schema, Table, column } from '@powersync/web';

// Phase 0: minimal schema — proof-of-concept sync table
const testItems = new Table({
  tenant_id: column.text,
  label: column.text,
  created_at: column.text,
});

export const AppSchema = new Schema({ test_items: testItems });

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
