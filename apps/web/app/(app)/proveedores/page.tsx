'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type Vendor } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Table } from '@/components/Table';

export default function ProveedoresPage() {
  const [vendors, setVendors] = useState<Vendor[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: '', tax_id: '', email: '', phone: '', contact_name: '' });

  useEffect(() => {
    api.vendors.list()
      .then(setVendors)
      .catch((err) => console.error('operation failed:', err))
      .finally(() => setLoading(false));
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const v = await api.vendors.create(form);
    setVendors((prev) => [v, ...prev]);
    setForm({ name: '', tax_id: '', email: '', phone: '', contact_name: '' });
    setShowForm(false);
  }

  async function remove(id: string) {
    await api.vendors.delete(id);
    setVendors((prev) => prev.filter((v) => v.id !== id));
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Proveedores</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nuevo proveedor'}</Button>
      </div>

      {showForm && (
        <Card title="Nuevo proveedor">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="Nombre" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              <Input label="RIF" value={form.tax_id} onChange={(e) => setForm({ ...form, tax_id: e.target.value })} />
              <Input label="Email" type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
              <Input label="Teléfono" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
              <Input label="Contacto" value={form.contact_name} onChange={(e) => setForm({ ...form, contact_name: e.target.value })} />
            </div>
            <Button type="submit">Guardar</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<Vendor>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'tax_id', label: 'RIF' },
              { key: 'contact_name', label: 'Contacto' },
              { key: 'phone', label: 'Teléfono' },
              {
                key: 'actions',
                label: '',
                render: (row) => (
                  <button onClick={() => remove(row.id)} className="text-red-600 text-xs hover:underline">Eliminar</button>
                ),
              },
            ]}
            data={vendors}
            empty="Sin proveedores registrados"
          />
        )}
      </Card>
    </div>
  );
}
