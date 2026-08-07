'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type Customer } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Table } from '@/components/Table';

export default function ClientesPage() {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: '', tax_id: '', email: '', phone: '', address: '' });

  useEffect(() => {
    api.customers.list()
      .then(setCustomers)
      .catch((err) => console.error('operation failed:', err))
      .finally(() => setLoading(false));
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const c = await api.customers.create(form);
    setCustomers((prev) => [c, ...prev]);
    setForm({ name: '', tax_id: '', email: '', phone: '', address: '' });
    setShowForm(false);
  }

  async function remove(id: string) {
    await api.customers.delete(id);
    setCustomers((prev) => prev.filter((c) => c.id !== id));
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Clientes</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nuevo cliente'}</Button>
      </div>

      {showForm && (
        <Card title="Nuevo cliente">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="Nombre" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              <Input label="RIF / Cédula" value={form.tax_id} onChange={(e) => setForm({ ...form, tax_id: e.target.value })} />
              <Input label="Email" type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
              <Input label="Teléfono" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
            </div>
            <Input label="Dirección" value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} />
            <Button type="submit">Guardar</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<Customer>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'tax_id', label: 'RIF' },
              { key: 'email', label: 'Email' },
              { key: 'phone', label: 'Teléfono' },
              {
                key: 'actions',
                label: '',
                render: (row) => (
                  <button onClick={() => remove(row.id)} className="text-red-600 text-xs hover:underline">Eliminar</button>
                ),
              },
            ]}
            data={customers}
            empty="Sin clientes registrados"
          />
        )}
      </Card>
    </div>
  );
}
