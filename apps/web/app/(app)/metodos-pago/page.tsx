'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type PaymentMethod } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function MetodosPagoPage() {
  const [methods, setMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: '', type: 'cash', currency: 'VES', allow_invoices: true, allow_change: false, allow_refunds: false, is_active: true });

  useEffect(() => {
    api.paymentMethods.list()
      .then(setMethods)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const m = await api.paymentMethods.create(form);
    setMethods((prev) => [m, ...prev]);
    setForm({ name: '', type: 'cash', currency: 'VES', allow_invoices: true, allow_change: false, allow_refunds: false, is_active: true });
    setShowForm(false);
  }

  async function remove(id: string) {
    await api.paymentMethods.delete(id);
    setMethods((prev) => prev.filter((m) => m.id !== id));
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Métodos de Pago</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nuevo método'}</Button>
      </div>

      {showForm && (
        <Card title="Nuevo método de pago">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="Nombre" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              <Select label="Tipo" value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })}>
                <option value="cash">Efectivo</option>
                <option value="card">Tarjeta</option>
                <option value="transfer">Transferencia</option>
                <option value="mobile_payment">Pago móvil</option>
                <option value="credit">Crédito</option>
              </Select>
              <Select label="Moneda" value={form.currency} onChange={(e) => setForm({ ...form, currency: e.target.value })}>
                <option value="VES">Bolívares (VES)</option>
                <option value="USD">Dólares (USD)</option>
              </Select>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <label className="flex items-center gap-2 text-sm text-gray-700">
                <input type="checkbox" checked={form.allow_invoices} onChange={(e) => setForm({ ...form, allow_invoices: e.target.checked })} />
                Permite facturas
              </label>
              <label className="flex items-center gap-2 text-sm text-gray-700">
                <input type="checkbox" checked={form.allow_change} onChange={(e) => setForm({ ...form, allow_change: e.target.checked })} />
                Permite vueltos
              </label>
              <label className="flex items-center gap-2 text-sm text-gray-700">
                <input type="checkbox" checked={form.allow_refunds} onChange={(e) => setForm({ ...form, allow_refunds: e.target.checked })} />
                Permite reembolsos
              </label>
              <label className="flex items-center gap-2 text-sm text-gray-700">
                <input type="checkbox" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} />
                Activo
              </label>
            </div>
            <Button type="submit">Guardar</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<PaymentMethod>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'type', label: 'Tipo' },
              { key: 'currency', label: 'Moneda' },
              { key: 'invoices', label: 'Facturas', render: (r) => (r.allow_invoices ? 'Sí' : 'No') },
              { key: 'change', label: 'Vueltos', render: (r) => (r.allow_change ? 'Sí' : 'No') },
              { key: 'refunds', label: 'Reembolsos', render: (r) => (r.allow_refunds ? 'Sí' : 'No') },
              { key: 'status', label: 'Estado', render: (r) => (r.is_active ? 'Activo' : 'Inactivo') },
              { key: 'actions', label: '', render: (row) => <button onClick={() => remove(row.id)} className="text-red-600 text-xs hover:underline">Eliminar</button> },
            ]}
            data={methods}
            empty="Sin métodos de pago registrados"
          />
        )}
      </Card>
    </div>
  );
}
