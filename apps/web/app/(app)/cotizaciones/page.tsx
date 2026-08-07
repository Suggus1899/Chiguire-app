'use client';
import { useState, useEffect, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { api, type Quotation, type Customer } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function CotizacionesPage() {
  const router = useRouter();
  const [quotations, setQuotations] = useState<Quotation[]>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ customer_id: '', currency: 'VES', valid_until: '' });

  useEffect(() => {
    Promise.all([
      api.quotations.list().catch((err) => { console.error('operation failed:', err); return [] as Quotation[]; }),
      api.customers.list().catch((err) => { console.error('operation failed:', err); return [] as Customer[]; }),
    ]).then(([quotes, cust]) => {
      setQuotations(quotes);
      setCustomers(cust);
      setLoading(false);
    });
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!form.customer_id) return;
    const q = await api.quotations.create({
      ...form,
      total: 0,
      status: 'draft',
    });
    setQuotations((prev) => [q, ...prev]);
    setForm({ customer_id: '', currency: 'VES', valid_until: '' });
    setShowForm(false);
  }

  async function convert(id: string) {
    const inv = await api.quotations.convertToInvoice(id);
    router.push(`/facturas/${inv.id}`);
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Cotizaciones</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nueva cotización'}</Button>
      </div>

      {showForm && (
        <Card title="Nueva cotización">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Select label="Cliente" required value={form.customer_id} onChange={(e) => setForm({ ...form, customer_id: e.target.value })}>
                <option value="">Seleccionar...</option>
                {customers.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
              </Select>
              <Select label="Moneda" value={form.currency} onChange={(e) => setForm({ ...form, currency: e.target.value })}>
                <option value="VES">Bolívares (VES)</option>
                <option value="USD">Dólares (USD)</option>
              </Select>
              <Input label="Válida hasta" type="date" value={form.valid_until} onChange={(e) => setForm({ ...form, valid_until: e.target.value })} />
            </div>
            <Button type="submit">Crear</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<Quotation>
            columns={[
              { key: 'number', label: 'Número', render: (r) => r.number || r.id.slice(0, 8) },
              { key: 'status', label: 'Estado' },
              { key: 'currency', label: 'Moneda' },
              { key: 'total', label: 'Total', render: (r) => (r.total ?? 0).toFixed(2) },
              {
                key: 'actions',
                label: '',
                render: (row) => (
                  <button onClick={() => convert(row.id)} className="text-blue-600 text-xs hover:underline">Convertir a factura</button>
                ),
              },
            ]}
            data={quotations}
            empty="Sin cotizaciones"
          />
        )}
      </Card>
    </div>
  );
}
