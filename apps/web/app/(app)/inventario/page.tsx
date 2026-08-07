'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type StockMovement, type Product, type Warehouse } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function InventarioPage() {
  const [movements, setMovements] = useState<StockMovement[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ product_id: '', warehouse_id: '', movement_type: 'in', quantity: '0', reference: '' });

  useEffect(() => {
    Promise.all([
      api.stock.listMovements().catch(() => [] as StockMovement[]),
      api.products.list().catch(() => [] as Product[]),
      api.warehouses.list().catch(() => [] as Warehouse[]),
    ]).then(([mov, prod, wh]) => {
      setMovements(mov);
      setProducts(prod);
      setWarehouses(wh);
      setLoading(false);
    });
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const m = await api.stock.createMovement({
      ...form,
      quantity: parseFloat(form.quantity) || 0,
    });
    setMovements((prev) => [m, ...prev]);
    setForm({ product_id: '', warehouse_id: '', movement_type: 'in', quantity: '0', reference: '' });
    setShowForm(false);
  }

  const productName = (id: string) => products.find((p) => p.id === id)?.name ?? id.slice(0, 8);
  const whName = (id: string) => warehouses.find((w) => w.id === id)?.name ?? id.slice(0, 8);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Inventario</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nuevo movimiento'}</Button>
      </div>

      {showForm && (
        <Card title="Nuevo movimiento de stock">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Select label="Producto" required value={form.product_id} onChange={(e) => setForm({ ...form, product_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
              </Select>
              <Select label="Almacén" required value={form.warehouse_id} onChange={(e) => setForm({ ...form, warehouse_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
              </Select>
              <Select label="Tipo" value={form.movement_type} onChange={(e) => setForm({ ...form, movement_type: e.target.value })}>
                <option value="in">Entrada</option>
                <option value="out">Salida</option>
                <option value="adjust">Ajuste</option>
              </Select>
              <Input label="Cantidad" type="number" step="0.01" required value={form.quantity} onChange={(e) => setForm({ ...form, quantity: e.target.value })} />
            </div>
            <Input label="Referencia" value={form.reference} onChange={(e) => setForm({ ...form, reference: e.target.value })} />
            <Button type="submit">Registrar</Button>
          </form>
        </Card>
      )}

      <Card title="Movimientos de stock">
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<StockMovement>
            columns={[
              { key: 'product_id', label: 'Producto', render: (r) => productName(r.product_id) },
              { key: 'warehouse_id', label: 'Almacén', render: (r) => whName(r.warehouse_id) },
              { key: 'movement_type', label: 'Tipo' },
              { key: 'quantity', label: 'Cantidad' },
              { key: 'reference', label: 'Referencia' },
              { key: 'created_at', label: 'Fecha', render: (r) => new Date(r.created_at).toLocaleDateString('es-VE') },
            ]}
            data={movements}
            empty="Sin movimientos"
          />
        )}
      </Card>
    </div>
  );
}
