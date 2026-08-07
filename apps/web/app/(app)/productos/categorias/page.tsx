'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type Category, type Unit } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Table } from '@/components/Table';

export default function CategoriasPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [units, setUnits] = useState<Unit[]>([]);
  const [loading, setLoading] = useState(true);
  const [catName, setCatName] = useState('');
  const [unitName, setUnitName] = useState('');
  const [unitSymbol, setUnitSymbol] = useState('');

  useEffect(() => {
    Promise.all([
      api.categories.list().catch((err) => { console.error('operation failed:', err); return [] as Category[]; }),
      api.units.list().catch((err) => { console.error('operation failed:', err); return [] as Unit[]; }),
    ]).then(([cat, uni]) => {
      setCategories(cat);
      setUnits(uni);
      setLoading(false);
    });
  }, []);

  async function addCategory(e: FormEvent) {
    e.preventDefault();
    if (!catName) return;
    const c = await api.categories.create({ name: catName });
    setCategories((prev) => [...prev, c]);
    setCatName('');
  }

  async function addUnit(e: FormEvent) {
    e.preventDefault();
    if (!unitName) return;
    const u = await api.units.create({ name: unitName, symbol: unitSymbol });
    setUnits((prev) => [...prev, u]);
    setUnitName('');
    setUnitSymbol('');
  }

  if (loading) return <p className="text-gray-500">Cargando...</p>;

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Categorías y Unidades</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card title="Categorías">
          <form onSubmit={addCategory} className="flex gap-2 mb-4">
            <Input placeholder="Nombre categoría" value={catName} onChange={(e) => setCatName(e.target.value)} />
            <Button type="submit">Agregar</Button>
          </form>
          <Table<Category>
            columns={[{ key: 'name', label: 'Nombre' }]}
            data={categories}
            empty="Sin categorías"
          />
        </Card>

        <Card title="Unidades de medida">
          <form onSubmit={addUnit} className="flex gap-2 mb-4">
            <Input placeholder="Nombre" value={unitName} onChange={(e) => setUnitName(e.target.value)} />
            <Input placeholder="Símbolo" value={unitSymbol} onChange={(e) => setUnitSymbol(e.target.value)} className="w-24" />
            <Button type="submit">Agregar</Button>
          </form>
          <Table<Unit>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'symbol', label: 'Símbolo' },
            ]}
            data={units}
            empty="Sin unidades"
          />
        </Card>
      </div>
    </div>
  );
}
