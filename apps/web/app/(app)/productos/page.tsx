'use client';
import { useState, useEffect, FormEvent } from 'react';
import Link from 'next/link';
import { api, type Product, type Category, type Unit } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function ProductosPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [units, setUnits] = useState<Unit[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ sku: '', name: '', description: '', category_id: '', unit_id: '', cost: '0' });

  useEffect(() => {
    Promise.all([
      api.products.list().catch((err) => { console.error('operation failed:', err); return [] as Product[]; }),
      api.categories.list().catch((err) => { console.error('operation failed:', err); return [] as Category[]; }),
      api.units.list().catch((err) => { console.error('operation failed:', err); return [] as Unit[]; }),
    ]).then(([prod, cat, uni]) => {
      setProducts(prod);
      setCategories(cat);
      setUnits(uni);
      setLoading(false);
    });
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const p = await api.products.create({
      ...form,
      cost: parseFloat(form.cost) || 0,
      is_active: true,
    });
    setProducts((prev) => [p, ...prev]);
    setForm({ sku: '', name: '', description: '', category_id: '', unit_id: '', cost: '0' });
    setShowForm(false);
  }

  async function remove(id: string) {
    await api.products.delete(id);
    setProducts((prev) => prev.filter((p) => p.id !== id));
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Productos</h1>
        <div className="flex gap-2">
          <Link href="/productos/categorias"><Button variant="secondary">Categorías</Button></Link>
          <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nuevo producto'}</Button>
        </div>
      </div>

      {showForm && (
        <Card title="Nuevo producto">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="SKU" required value={form.sku} onChange={(e) => setForm({ ...form, sku: e.target.value })} />
              <Input label="Nombre" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              <Select label="Categoría" value={form.category_id} onChange={(e) => setForm({ ...form, category_id: e.target.value })}>
                <option value="">Sin categoría</option>
                {categories.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
              </Select>
              <Select label="Unidad" value={form.unit_id} onChange={(e) => setForm({ ...form, unit_id: e.target.value })}>
                <option value="">Sin unidad</option>
                {units.map((u) => <option key={u.id} value={u.id}>{u.name}</option>)}
              </Select>
              <Input label="Costo" type="number" step="0.01" value={form.cost} onChange={(e) => setForm({ ...form, cost: e.target.value })} />
            </div>
            <Input label="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
            <Button type="submit">Guardar</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<Product>
            columns={[
              { key: 'sku', label: 'SKU' },
              { key: 'name', label: 'Nombre' },
              { key: 'cost', label: 'Costo', render: (r) => (r.cost ?? 0).toFixed(2) },
              {
                key: 'actions',
                label: '',
                render: (row) => (
                  <button onClick={() => remove(row.id)} className="text-red-600 text-xs hover:underline">Eliminar</button>
                ),
              },
            ]}
            data={products}
            empty="Sin productos registrados"
          />
        )}
      </Card>
    </div>
  );
}
