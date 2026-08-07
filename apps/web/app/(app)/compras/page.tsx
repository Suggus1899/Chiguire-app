'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type PurchaseOrder, type Vendor, type PurchaseOrderItem } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';
import { Modal } from '@/components/Modal';

export default function ComprasPage() {
  const [orders, setOrders] = useState<PurchaseOrder[]>([]);
  const [vendors, setVendors] = useState<Vendor[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ vendor_id: '', currency: 'VES' });
  const [selectedPO, setSelectedPO] = useState<PurchaseOrder | null>(null);
  const [poItems, setPoItems] = useState<PurchaseOrderItem[]>([]);
  const [receiveQty, setReceiveQty] = useState<Record<string, string>>({});

  useEffect(() => {
    Promise.all([
      api.purchases.list().catch(() => [] as PurchaseOrder[]),
      api.vendors.list().catch(() => [] as Vendor[]),
    ]).then(([po, vend]) => {
      setOrders(po);
      setVendors(vend);
      setLoading(false);
    });
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!form.vendor_id) return;
    const po = await api.purchases.create({ ...form, total: 0, status: 'draft' });
    setOrders((prev) => [po, ...prev]);
    setForm({ vendor_id: '', currency: 'VES' });
    setShowForm(false);
  }

  async function approve(id: string) {
    const po = await api.purchases.approve(id);
    setOrders((prev) => prev.map((p) => (p.id === id ? po : p)));
  }

  async function openDetail(po: PurchaseOrder) {
    setSelectedPO(po);
    try {
      const full = await api.purchases.get(po.id);
      setPoItems((full as unknown as { items?: PurchaseOrderItem[] }).items ?? []);
    } catch {
      setPoItems([]);
    }
  }

  async function receiveItem(itemId: string) {
    if (!selectedPO) return;
    const qty = parseFloat(receiveQty[itemId] ?? '0');
    if (qty <= 0) return;
    const updated = await api.purchases.receiveItem(selectedPO.id, itemId, { quantity: qty });
    setPoItems((prev) => prev.map((it) => (it.id === itemId ? updated : it)));
    setReceiveQty((prev) => ({ ...prev, [itemId]: '' }));
  }

  const vendorName = (id: string) => vendors.find((v) => v.id === id)?.name ?? id.slice(0, 8);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Órdenes de compra</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nueva orden'}</Button>
      </div>

      {showForm && (
        <Card title="Nueva orden de compra">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Select label="Proveedor" required value={form.vendor_id} onChange={(e) => setForm({ ...form, vendor_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {vendors.map((v) => <option key={v.id} value={v.id}>{v.name}</option>)}
              </Select>
              <Select label="Moneda" value={form.currency} onChange={(e) => setForm({ ...form, currency: e.target.value })}>
                <option value="VES">VES</option>
                <option value="USD">USD</option>
              </Select>
            </div>
            <Button type="submit">Crear</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<PurchaseOrder>
            columns={[
              { key: 'number', label: 'Número', render: (r) => r.number || r.id.slice(0, 8) },
              { key: 'vendor_id', label: 'Proveedor', render: (r) => vendorName(r.vendor_id) },
              { key: 'status', label: 'Estado' },
              { key: 'total', label: 'Total', render: (r) => (r.total ?? 0).toFixed(2) },
              {
                key: 'actions',
                label: '',
                render: (row) => (
                  <div className="flex gap-2">
                    <button onClick={() => openDetail(row)} className="text-blue-600 text-xs hover:underline">Detalle</button>
                    {row.status === 'draft' && (
                      <button onClick={() => approve(row.id)} className="text-green-600 text-xs hover:underline">Aprobar</button>
                    )}
                  </div>
                ),
              },
            ]}
            data={orders}
            empty="Sin órdenes de compra"
          />
        )}
      </Card>

      <Modal open={!!selectedPO} onClose={() => setSelectedPO(null)} title={`Orden ${selectedPO?.number ?? ''}`}>
        {poItems.length === 0 ? (
          <p className="text-sm text-gray-400">Sin items en esta orden.</p>
        ) : (
          <div className="space-y-3">
            {poItems.map((item) => (
              <div key={item.id} className="flex items-end gap-2 text-sm">
                <div className="flex-1">
                  <p className="text-gray-800">{item.product_id.slice(0, 8)}</p>
                  <p className="text-gray-400 text-xs">Recibido: {item.received_quantity ?? 0} / {item.quantity}</p>
                </div>
                <Input type="number" step="0.01" placeholder="Cant." value={receiveQty[item.id] ?? ''} onChange={(e) => setReceiveQty({ ...receiveQty, [item.id]: e.target.value })} className="w-24" />
                <Button variant="secondary" onClick={() => receiveItem(item.id)} className="!px-2 !py-1">Recibir</Button>
              </div>
            ))}
          </div>
        )}
      </Modal>
    </div>
  );
}
