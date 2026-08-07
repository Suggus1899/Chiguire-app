'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type Seller } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';
import { Modal } from '@/components/Modal';

export default function VendedoresPage() {
  const [sellers, setSellers] = useState<Seller[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [statusFilter, setStatusFilter] = useState('all');
  const [menuId, setMenuId] = useState<string | null>(null);
  const [editSeller, setEditSeller] = useState<Seller | null>(null);
  const [form, setForm] = useState({ name: '', email: '', phone: '', commission_pct: '0', is_active: true });

  useEffect(() => {
    api.sellers.list()
      .then(setSellers)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const data = { ...form, commission_pct: parseFloat(form.commission_pct) || 0 };
    if (editSeller) {
      const updated = await api.sellers.update(editSeller.id, data);
      setSellers((prev) => prev.map((s) => (s.id === updated.id ? updated : s)));
    } else {
      const s = await api.sellers.create(data);
      setSellers((prev) => [s, ...prev]);
    }
    resetForm();
  }

  function resetForm() {
    setForm({ name: '', email: '', phone: '', commission_pct: '0', is_active: true });
    setEditSeller(null);
    setShowForm(false);
  }

  function startEdit(s: Seller) {
    setEditSeller(s);
    setForm({ name: s.name, email: s.email, phone: s.phone, commission_pct: String(s.commission_pct), is_active: s.is_active });
    setShowForm(true);
    setMenuId(null);
  }

  async function remove(id: string) {
    await api.sellers.delete(id);
    setSellers((prev) => prev.filter((s) => s.id !== id));
    setMenuId(null);
  }

  const filtered = sellers.filter((s) =>
    statusFilter === 'all' ? true : statusFilter === 'activo' ? s.is_active : !s.is_active,
  );

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Vendedores</h1>
        <Button onClick={() => { resetForm(); setShowForm(!showForm); }}>{showForm ? 'Cancelar' : 'Nuevo vendedor'}</Button>
      </div>

      <div className="flex gap-3 items-end">
        <Select label="Estado" value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="max-w-xs">
          <option value="all">Todos</option>
          <option value="activo">Activo</option>
          <option value="inactivo">Inactivo</option>
        </Select>
      </div>

      {showForm && (
        <Card title={editSeller ? 'Editar vendedor' : 'Nuevo vendedor'}>
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="Nombre" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              <Input label="Email" type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
              <Input label="Teléfono" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
              <Input label="Comisión (%)" type="number" step="0.01" value={form.commission_pct} onChange={(e) => setForm({ ...form, commission_pct: e.target.value })} />
            </div>
            <label className="flex items-center gap-2 text-sm text-gray-700">
              <input type="checkbox" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} />
              Activo
            </label>
            <Button type="submit">Guardar</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<Seller>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'email', label: 'Email' },
              { key: 'phone', label: 'Teléfono' },
              { key: 'commission_pct', label: 'Comisión %', render: (r) => `${r.commission_pct}%` },
              { key: 'status', label: 'Estado', render: (r) => (r.is_active ? 'Activo' : 'Inactivo') },
              {
                key: 'actions',
                label: '',
                render: (row) => (
                  <div className="relative">
                    <button onClick={() => setMenuId(menuId === row.id ? null : row.id)} className="text-gray-400 hover:text-gray-700 px-2">⋮</button>
                    {menuId === row.id && (
                      <div className="absolute right-0 top-6 bg-white border border-gray-200 rounded-lg shadow-lg z-10 w-32 text-sm">
                        <button onClick={() => startEdit(row)} className="block w-full text-left px-3 py-2 hover:bg-gray-50">Editar</button>
                        <button onClick={() => remove(row.id)} className="block w-full text-left px-3 py-2 text-red-600 hover:bg-gray-50">Eliminar</button>
                      </div>
                    )}
                  </div>
                ),
              },
            ]}
            data={filtered}
            empty="Sin vendedores registrados"
          />
        )}
      </Card>
    </div>
  );
}
