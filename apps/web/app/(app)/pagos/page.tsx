'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type PaymentLink, type Invoice } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function PagosPage() {
  const [links, setLinks] = useState<PaymentLink[]>([]);
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ invoice_id: '', provider: 'cashea', amount: '0', currency: 'VES' });

  useEffect(() => {
    Promise.all([
      api.payments.listLinks().catch((err) => { console.error('operation failed:', err); return [] as PaymentLink[]; }),
      api.invoices.list().catch((err) => { console.error('operation failed:', err); return [] as Invoice[]; }),
    ]).then(([lnk, inv]) => {
      setLinks(lnk);
      setInvoices(inv);
      setLoading(false);
    });
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!form.invoice_id) return;
    const link = await api.payments.createLink({
      ...form,
      amount: parseFloat(form.amount) || 0,
    });
    setLinks((prev) => [link, ...prev]);
    setForm({ invoice_id: '', provider: 'cashea', amount: '0', currency: 'VES' });
    setShowForm(false);
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Pagos</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nuevo link de pago'}</Button>
      </div>

      {showForm && (
        <Card title="Nuevo link de pago">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Select label="Factura" required value={form.invoice_id} onChange={(e) => setForm({ ...form, invoice_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {invoices.map((inv) => <option key={inv.id} value={inv.id}>{inv.number || inv.id.slice(0, 8)}</option>)}
              </Select>
              <Select label="Proveedor" value={form.provider} onChange={(e) => setForm({ ...form, provider: e.target.value })}>
                <option value="cashea">Cashea</option>
                <option value="spidi">Spidi</option>
                <option value="wayupay">WayuPay</option>
                <option value="biopago">Biopago</option>
              </Select>
              <Input label="Monto" type="number" step="0.01" required value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} />
              <Select label="Moneda" value={form.currency} onChange={(e) => setForm({ ...form, currency: e.target.value })}>
                <option value="VES">VES</option>
                <option value="USD">USD</option>
              </Select>
            </div>
            <Button type="submit">Crear link</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<PaymentLink>
            columns={[
              { key: 'provider', label: 'Proveedor' },
              { key: 'amount', label: 'Monto', render: (r) => (r.amount ?? 0).toFixed(2) },
              { key: 'currency', label: 'Moneda' },
              { key: 'status', label: 'Estado' },
              {
                key: 'url',
                label: 'Link',
                render: (r) => r.url ? <a href={r.url} target="_blank" rel="noopener noreferrer" className="text-blue-600 text-xs hover:underline">Abrir</a> : '—',
              },
            ]}
            data={links}
            empty="Sin links de pago"
          />
        )}
      </Card>
    </div>
  );
}
