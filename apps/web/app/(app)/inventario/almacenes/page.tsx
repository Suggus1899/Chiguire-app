'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type Branch, type Warehouse } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function AlmacenesPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [loading, setLoading] = useState(true);
  const [branchForm, setBranchForm] = useState({ name: '', address: '', phone: '' });
  const [whForm, setWhForm] = useState({ branch_id: '', name: '', code: '' });

  useEffect(() => {
    Promise.all([
      api.branches.list().catch((err) => { console.error('operation failed:', err); return [] as Branch[]; }),
      api.warehouses.list().catch((err) => { console.error('operation failed:', err); return [] as Warehouse[]; }),
    ]).then(([br, wh]) => {
      setBranches(br);
      setWarehouses(wh);
      setLoading(false);
    });
  }, []);

  async function addBranch(e: FormEvent) {
    e.preventDefault();
    if (!branchForm.name) return;
    const b = await api.branches.create({ ...branchForm, is_active: true });
    setBranches((prev) => [...prev, b]);
    setBranchForm({ name: '', address: '', phone: '' });
  }

  async function addWarehouse(e: FormEvent) {
    e.preventDefault();
    if (!whForm.name) return;
    const w = await api.warehouses.create(whForm);
    setWarehouses((prev) => [...prev, w]);
    setWhForm({ branch_id: '', name: '', code: '' });
  }

  const branchName = (id: string) => branches.find((b) => b.id === id)?.name ?? '—';

  if (loading) return <p className="text-gray-500">Cargando...</p>;

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Sucursales y Almacenes</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card title="Sucursales">
          <form onSubmit={addBranch} className="space-y-3 mb-4">
            <Input label="Nombre" required value={branchForm.name} onChange={(e) => setBranchForm({ ...branchForm, name: e.target.value })} />
            <Input label="Dirección" value={branchForm.address} onChange={(e) => setBranchForm({ ...branchForm, address: e.target.value })} />
            <Input label="Teléfono" value={branchForm.phone} onChange={(e) => setBranchForm({ ...branchForm, phone: e.target.value })} />
            <Button type="submit">Agregar sucursal</Button>
          </form>
          <Table<Branch>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'phone', label: 'Teléfono' },
            ]}
            data={branches}
            empty="Sin sucursales"
          />
        </Card>

        <Card title="Almacenes">
          <form onSubmit={addWarehouse} className="space-y-3 mb-4">
            <Select label="Sucursal" value={whForm.branch_id} onChange={(e) => setWhForm({ ...whForm, branch_id: e.target.value })}>
              <option value="">Sin sucursal</option>
              {branches.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
            </Select>
            <Input label="Nombre" required value={whForm.name} onChange={(e) => setWhForm({ ...whForm, name: e.target.value })} />
            <Input label="Código" value={whForm.code} onChange={(e) => setWhForm({ ...whForm, code: e.target.value })} />
            <Button type="submit">Agregar almacén</Button>
          </form>
          <Table<Warehouse>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'code', label: 'Código' },
              { key: 'branch_id', label: 'Sucursal', render: (r) => branchName(r.branch_id) },
            ]}
            data={warehouses}
            empty="Sin almacenes"
          />
        </Card>
      </div>
    </div>
  );
}
