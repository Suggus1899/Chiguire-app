'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type Transfer, type Warehouse, type Product } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

type TransferItem = { product_id: string; quantity: number };

export default function TransferenciasPage() {
  const [transfers, setTransfers] = useState<Transfer[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [fromWh, setFromWh] = useState('');
  const [toWh, setToWh] = useState('');
  const [items, setItems] = useState<TransferItem[]>([]);

  useEffect(() => {
    Promise.all([
      api.transfers.list().catch((err) => { console.error('operation failed:', err); return [] as Transfer[]; }),
      api.warehouses.list().catch((err) => { console.error('operation failed:', err); return [] as Warehouse[]; }),
      api.products.list().catch((err) => { console.error('operation failed:', err); return [] as Product[]; }),
    ]).then(([t, wh, prod]) => {
      setTransfers(t);
      setWarehouses(wh);
      setProducts(prod);
      setLoading(false);
    });
  }, []);

  function addItem() { setItems((prev) => [...prev, { product_id: '', quantity: 1 }]); }
  function updateItem(i: number, patch: Partial<TransferItem>) { setItems((prev) => prev.map((it, idx) => (idx === i ? { ...it, ...patch } : it))); }
  function removeItem(i: number) { setItems((prev) => prev.filter((_, idx) => idx !== i)); }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const t = await api.transfers.create({ from_warehouse_id: fromWh, to_warehouse_id: toWh });
    setTransfers((prev) => [t, ...prev]);
    setFromWh(''); setToWh(''); setItems([]);
    setShowForm(false);
  }

  async function ship(id: string) {
    const t = await api.transfers.ship(id);
    setTransfers((prev) => prev.map((x) => (x.id === id ? t : x)));
  }
  async function receive(id: string) {
    const t = await api.transfers.receive(id);
    setTransfers((prev) => prev.map((x) => (x.id === id ? t : x)));
  }

  const whName = (id: string) => warehouses.find((w) => w.id === id)?.name ?? '—';

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Transferencias de inventario</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nueva transferencia'}</Button>
      </div>

      {showForm && (
        <Card title="Nueva transferencia">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Select label="Almacén origen" required value={fromWh} onChange={(e) => setFromWh(e.target.value)}>
                <option value="">Seleccionar...</option>
                {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
              </Select>
              <Select label="Almacén destino" required value={toWh} onChange={(e) => setToWh(e.target.value)}>
                <option value="">Seleccionar...</option>
                {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
              </Select>
            </div>
            <div>
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm font-medium text-gray-700">Items</span>
                <Button variant="secondary" onClick={addItem}>Agregar item</Button>
              </div>
              {items.map((item, i) => (
                <div key={i} className="grid grid-cols-12 gap-2 items-end mb-2">
                  <div className="col-span-8">
                    <Select value={item.product_id} onChange={(e) => updateItem(i, { product_id: e.target.value })}>
                      <option value="">Seleccionar...</option>
                      {products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                    </Select>
                  </div>
                  <div className="col-span-3"><Input type="number" value={item.quantity} onChange={(e) => updateItem(i, { quantity: parseFloat(e.target.value) || 0 })} /></div>
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
          <Table<Transfer>
            columns={[
              { key: 'from', label: 'Origen', render: (r) => whName(r.from_warehouse_id) },
              { key: 'to', label: 'Destino', render: (r) => whName(r.to_warehouse_id) },
              { key: 'status', label: 'Estado' },
              { key: 'date', label: 'Fecha', render: (r) => r.created_at.slice(0, 10) },
              { key: 'actions', label: '', render: (row) => (
                <div className="flex gap-2">
                  {row.status === 'draft' && <button onClick={() => ship(row.id)} className="text-blue-600 text-xs hover:underline">Despachar</button>}
                  {row.status === 'shipped' && <button onClick={() => receive(row.id)} className="text-green-600 text-xs hover:underline">Recibir</button>}
                </div>
              ) },
            ]}
            data={transfers}
            empty="Sin transferencias registradas"
          />
        )}
      </Card>
    </div>
  );
}
