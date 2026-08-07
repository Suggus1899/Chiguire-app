'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type CreditNote, type Invoice, type Product } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

type NoteItem = { product_id: string; quantity: number; unit_price: number };

export default function NotasCreditoPage() {
  const [notes, setNotes] = useState<CreditNote[]>([]);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ invoice_id: '', type: 'credit', reason: '' });
  const [items, setItems] = useState<NoteItem[]>([]);

  useEffect(() => {
    Promise.all([
      api.creditNotes.list().catch(() => [] as CreditNote[]),
      api.invoices.list().catch(() => [] as Invoice[]),
      api.products.list().catch(() => [] as Product[]),
    ]).then(([n, inv, prod]) => {
      setNotes(n);
      setInvoices(inv);
      setProducts(prod);
      setLoading(false);
    });
  }, []);

  function addItem() {
    setItems((prev) => [...prev, { product_id: '', quantity: 1, unit_price: 0 }]);
  }
  function updateItem(i: number, patch: Partial<NoteItem>) {
    setItems((prev) => prev.map((it, idx) => (idx === i ? { ...it, ...patch } : it)));
  }
  function removeItem(i: number) {
    setItems((prev) => prev.filter((_, idx) => idx !== i));
  }
  const total = items.reduce((s, it) => s + it.quantity * it.unit_price, 0);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const n = await api.creditNotes.create({ ...form, total });
    setNotes((prev) => [n, ...prev]);
    setForm({ invoice_id: '', type: 'credit', reason: '' });
    setItems([]);
    setShowForm(false);
  }

  async function voidNote(id: string) {
    await api.creditNotes.void(id);
    setNotes((prev) => prev.map((n) => (n.id === id ? { ...n, status: 'void' } : n)));
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Notas de Crédito/Débito</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nueva nota'}</Button>
      </div>

      {showForm && (
        <Card title="Nueva nota">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <Select label="Factura" required value={form.invoice_id} onChange={(e) => setForm({ ...form, invoice_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {invoices.map((inv) => <option key={inv.id} value={inv.id}>{inv.number || inv.id.slice(0, 8)}</option>)}
              </Select>
              <Select label="Tipo" value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })}>
                <option value="credit">Nota de crédito</option>
                <option value="debit">Nota de débito</option>
              </Select>
              <Input label="Motivo" required value={form.reason} onChange={(e) => setForm({ ...form, reason: e.target.value })} />
            </div>

            <div>
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm font-medium text-gray-700">Items</span>
                <Button variant="secondary" onClick={addItem}>Agregar item</Button>
              </div>
              {items.map((item, i) => (
                <div key={i} className="grid grid-cols-12 gap-2 items-end mb-2">
                  <div className="col-span-6">
                    <Select value={item.product_id} onChange={(e) => { const p = products.find((x) => x.id === e.target.value); updateItem(i, { product_id: e.target.value, unit_price: p?.cost ?? 0 }); }}>
                      <option value="">Seleccionar...</option>
                      {products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                    </Select>
                  </div>
                  <div className="col-span-2"><Input type="number" value={item.quantity} onChange={(e) => updateItem(i, { quantity: parseFloat(e.target.value) || 0 })} /></div>
                  <div className="col-span-3"><Input type="number" step="0.01" value={item.unit_price} onChange={(e) => updateItem(i, { unit_price: parseFloat(e.target.value) || 0 })} /></div>
                  <div className="col-span-1"><Button variant="danger" onClick={() => removeItem(i)} className="!px-2">×</Button></div>
                </div>
              ))}
              <p className="text-sm font-semibold mt-2">Total: {total.toFixed(2)}</p>
            </div>

            <Button type="submit">Guardar</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<CreditNote>
            columns={[
              { key: 'number', label: 'Número' },
              { key: 'type', label: 'Tipo' },
              { key: 'reason', label: 'Motivo' },
              { key: 'total', label: 'Total', render: (r) => r.total.toFixed(2) },
              { key: 'status', label: 'Estado' },
              { key: 'actions', label: '', render: (row) => row.status !== 'void' ? <button onClick={() => voidNote(row.id)} className="text-red-600 text-xs hover:underline">Anular</button> : null },
            ]}
            data={notes}
            empty="Sin notas registradas"
          />
        )}
      </Card>
    </div>
  );
}
