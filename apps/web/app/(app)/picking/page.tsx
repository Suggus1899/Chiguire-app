'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type PickingList, type PickingItem, type Warehouse, type Product, type DeliveryRoute } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';
import { Modal } from '@/components/Modal';

type PickFormItem = { product_id: string; qty_required: number };

export default function PickingPage() {
  const [lists, setLists] = useState<PickingList[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [routes, setRoutes] = useState<DeliveryRoute[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ warehouse_id: '', route_id: '' });
  const [items, setItems] = useState<PickFormItem[]>([]);
  const [detail, setDetail] = useState<{ list: PickingList; items: PickingItem[] } | null>(null);
  const [pickValues, setPickValues] = useState<Record<string, { qty: string; verified: boolean }>>({});

  useEffect(() => {
    Promise.all([
      api.picking.list().catch(() => [] as PickingList[]),
      api.warehouses.list().catch(() => [] as Warehouse[]),
      api.delivery.listRoutes().catch(() => [] as DeliveryRoute[]),
      api.products.list().catch(() => [] as Product[]),
    ]).then(([l, wh, r, prod]) => {
      setLists(l); setWarehouses(wh); setRoutes(r); setProducts(prod);
      setLoading(false);
    });
  }, []);

  function addItem() { setItems((prev) => [...prev, { product_id: '', qty_required: 1 }]); }
  function updateItem(i: number, patch: Partial<PickFormItem>) { setItems((prev) => prev.map((it, idx) => (idx === i ? { ...it, ...patch } : it))); }
  function removeItem(i: number) { setItems((prev) => prev.filter((_, idx) => idx !== i)); }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const pl = await api.picking.create({ warehouse_id: form.warehouse_id, route_id: form.route_id || undefined });
    setLists((prev) => [pl, ...prev]);
    setForm({ warehouse_id: '', route_id: '' }); setItems([]);
    setShowForm(false);
  }

  async function openDetail(list: PickingList) {
    const data = await api.picking.get(list.id);
    setDetail({ list, items: (data as unknown as { items: PickingItem[] }).items ?? [] });
    setPickValues({});
  }

  async function verifyItem(itemId: string) {
    if (!detail) return;
    const val = pickValues[itemId] ?? { qty: '0', verified: false };
    await api.picking.verifyItem(detail.list.id, itemId, { qty_picked: parseFloat(val.qty) || 0 });
    setDetail((prev) => prev ? { ...prev, items: prev.items.map((it) => (it.id === itemId ? { ...it, qty_picked: parseFloat(val.qty) || 0, verified: val.verified } : it)) } : prev);
  }

  async function complete(id: string) {
    const pl = await api.picking.complete(id);
    setLists((prev) => prev.map((x) => (x.id === id ? pl : x)));
    setDetail(null);
  }

  const whName = (id: string) => warehouses.find((w) => w.id === id)?.name ?? '—';
  const routeName = (id: string | null) => routes.find((r) => r.id === id)?.driver_name ?? '—';
  const prodName = (id: string) => products.find((p) => p.id === id)?.name ?? '—';

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Listas de picking</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nueva lista'}</Button>
      </div>

      {showForm && (
        <Card title="Nueva lista de picking">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Select label="Almacén" required value={form.warehouse_id} onChange={(e) => setForm({ ...form, warehouse_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
              </Select>
              <Select label="Ruta (opcional)" value={form.route_id} onChange={(e) => setForm({ ...form, route_id: e.target.value })}>
                <option value="">Sin ruta</option>
                {routes.map((r) => <option key={r.id} value={r.id}>{r.driver_name} — {r.date}</option>)}
              </Select>
            </div>
            <div>
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm font-medium text-gray-700">Items</span>
                <Button variant="secondary" onClick={addItem}>Agregar item</Button>
              </div>
              {items.map((item, i) => (
                <div key={i} className="grid grid-cols-12 gap-2 items-end mb-2">
                  <div className="col-span-9">
                    <Select value={item.product_id} onChange={(e) => updateItem(i, { product_id: e.target.value })}>
                      <option value="">Seleccionar...</option>
                      {products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                    </Select>
                  </div>
                  <div className="col-span-2"><Input type="number" value={item.qty_required} onChange={(e) => updateItem(i, { qty_required: parseFloat(e.target.value) || 0 })} /></div>
                  <div className="col-span-1"><Button variant="danger" onClick={() => removeItem(i)} className="!px-2">×</Button></div>
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
          <Table<PickingList>
            columns={[
              { key: 'warehouse', label: 'Almacén', render: (r) => whName(r.warehouse_id) },
              { key: 'route', label: 'Ruta', render: (r) => routeName(r.route_id) },
              { key: 'status', label: 'Estado' },
              { key: 'actions', label: '', render: (row) => (
                <button onClick={() => openDetail(row)} className="text-blue-600 text-xs hover:underline">Detalle</button>
              ) },
            ]}
            data={lists}
            empty="Sin listas registradas"
          />
        )}
      </Card>

      <Modal open={!!detail} onClose={() => setDetail(null)} title="Detalle de picking">
        {detail && (
          <div className="space-y-4">
            <div className="space-y-2 max-h-96 overflow-y-auto">
              {detail.items.length === 0 ? (
                <p className="text-sm text-gray-400">Sin items en esta lista.</p>
              ) : detail.items.map((it) => (
                <div key={it.id} className="border border-gray-200 rounded-lg p-3 space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="font-medium text-gray-800">{prodName(it.product_id)}</span>
                    <span className="text-gray-500">Requerido: {it.qty_required}</span>
                  </div>
                  <div className="flex gap-2 items-end">
                    <Input type="number" placeholder="Cant. pickeada" value={pickValues[it.id]?.qty ?? ''} onChange={(e) => setPickValues((prev) => ({ ...prev, [it.id]: { qty: e.target.value, verified: prev[it.id]?.verified ?? false } }))} />
                    <label className="flex items-center gap-1 text-sm text-gray-700">
                      <input type="checkbox" checked={pickValues[it.id]?.verified ?? it.verified} onChange={(e) => setPickValues((prev) => ({ ...prev, [it.id]: { qty: prev[it.id]?.qty ?? '', verified: e.target.checked } }))} />
                      Verificar
                    </label>
                    <Button variant="secondary" onClick={() => verifyItem(it.id)} className="!px-2">OK</Button>
                  </div>
                </div>
              ))}
            </div>
            <Button onClick={() => complete(detail.list.id)}>Completar picking</Button>
          </div>
        )}
      </Modal>
    </div>
  );
}
