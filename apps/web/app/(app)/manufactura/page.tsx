'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type ManufacturingOrder, type Warehouse, type Product } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

type BomItem = { component_id: string; qty_per_unit: number };

export default function ManufacturaPage() {
  const [orders, setOrders] = useState<ManufacturingOrder[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ product_id: '', warehouse_id: '', quantity: '1' });
  const [bom, setBom] = useState<BomItem[]>([]);

  useEffect(() => {
    Promise.all([
      api.manufacturing.list().catch(() => [] as ManufacturingOrder[]),
      api.warehouses.list().catch(() => [] as Warehouse[]),
      api.products.list().catch(() => [] as Product[]),
    ]).then(([o, wh, prod]) => {
      setOrders(o); setWarehouses(wh); setProducts(prod);
      setLoading(false);
    });
  }, []);

  function addBom() { setBom((prev) => [...prev, { component_id: '', qty_per_unit: 1 }]); }
  function updateBom(i: number, patch: Partial<BomItem>) { setBom((prev) => prev.map((it, idx) => (idx === i ? { ...it, ...patch } : it))); }
  function removeBom(i: number) { setBom((prev) => prev.filter((_, idx) => idx !== i)); }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const o = await api.manufacturing.create({ ...form, quantity: parseFloat(form.quantity) || 0 });
    setOrders((prev) => [o, ...prev]);
    setForm({ product_id: '', warehouse_id: '', quantity: '1' });
    setBom([]);
    setShowForm(false);
  }

  async function start(id: string) {
    const o = await api.manufacturing.start(id);
    setOrders((prev) => prev.map((x) => (x.id === id ? o : x)));
  }
  async function complete(id: string) {
    const o = await api.manufacturing.complete(id);
    setOrders((prev) => prev.map((x) => (x.id === id ? o : x)));
  }

  const prodName = (id: string) => products.find((p) => p.id === id)?.name ?? '—';
  const whName = (id: string) => warehouses.find((w) => w.id === id)?.name ?? '—';

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Órdenes de manufactura</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nueva orden'}</Button>
      </div>

      {showForm && (
        <Card title="Nueva orden de manufactura">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <Select label="Producto" required value={form.product_id} onChange={(e) => setForm({ ...form, product_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
              </Select>
              <Select label="Almacén" required value={form.warehouse_id} onChange={(e) => setForm({ ...form, warehouse_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
              </Select>
              <Input label="Cantidad" type="number" value={form.quantity} onChange={(e) => setForm({ ...form, quantity: e.target.value })} />
            </div>
            <div>
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm font-medium text-gray-700">Lista de materiales (BOM)</span>
                <Button variant="secondary" onClick={addBom}>Agregar componente</Button>
              </div>
              {bom.map((item, i) => (
                <div key={i} className="grid grid-cols-12 gap-2 items-end mb-2">
                  <div className="col-span-9">
                    <Select value={item.component_id} onChange={(e) => updateBom(i, { component_id: e.target.value })}>
                      <option value="">Seleccionar...</option>
                      {products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                    </Select>
                  </div>
                  <div className="col-span-2"><Input type="number" step="0.01" value={item.qty_per_unit} onChange={(e) => updateBom(i, { qty_per_unit: parseFloat(e.target.value) || 0 })} /></div>
                  <div className="col-span-1"><Button variant="danger" onClick={() => removeBom(i)} className="!px-2">×</Button></div>
                </div>
              ))}
            </div>
            <Button type="submit">Guardar</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<ManufacturingOrder>
            columns={[
              { key: 'product', label: 'Producto', render: (r) => prodName(r.product_id) },
              { key: 'warehouse', label: 'Almacén', render: (r) => whName(r.warehouse_id) },
              { key: 'quantity', label: 'Cantidad' },
              { key: 'status', label: 'Estado' },
              { key: 'started', label: 'Inicio', render: (r) => r.started_at?.slice(0, 10) ?? '—' },
              { key: 'completed', label: 'Fin', render: (r) => r.completed_at?.slice(0, 10) ?? '—' },
              { key: 'actions', label: '', render: (row) => (
                <div className="flex gap-2">
                  {row.status === 'draft' && <button onClick={() => start(row.id)} className="text-blue-600 text-xs hover:underline">Iniciar</button>}
                  {row.status === 'in_progress' && <button onClick={() => complete(row.id)} className="text-green-600 text-xs hover:underline">Completar</button>}
                </div>
              ) },
            ]}
            data={orders}
            empty="Sin órdenes registradas"
          />
        )}
      </Card>
    </div>
  );
}
