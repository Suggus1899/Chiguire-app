'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';

const TOTAL_STEPS = 5;

export default function OnboardingPage() {
  const router = useRouter();
  const [step, setStep] = useState(0);
  const [saving, setSaving] = useState(false);
  const [company, setCompany] = useState({ name: '', rif: '', address: '' });
  const [branch, setBranch] = useState({ branchName: '', branchAddress: '', warehouseName: '', warehouseCode: '' });
  const [tax, setTax] = useState({ iva: '16', igtf: '0' });
  const [product, setProduct] = useState({ sku: '', name: '', cost: '0' });
  const [fiscal, setFiscal] = useState({ name: '', type: 'forma-libre', model: '', serial: '' });

  const progress = ((step + 1) / TOTAL_STEPS) * 100;

  async function finish() {
    setSaving(true);
    try {
      await api.branches.create({ name: company.name || 'Matriz', address: company.address });
      await api.warehouses.create({ name: branch.warehouseName || 'Almacén principal', code: branch.warehouseCode || 'ALM01' });
      await api.products.create({ sku: product.sku, name: product.name, cost: parseFloat(product.cost) || 0 });
      await api.fiscalDevices.create({ name: fiscal.name, type: fiscal.type, model: fiscal.model, serial: fiscal.serial });
      router.push('/dashboard');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Error al guardar');
    } finally {
      setSaving(false);
    }
  }

  function next() {
    if (step < TOTAL_STEPS - 1) setStep(step + 1);
    else finish();
  }
  function back() {
    if (step > 0) setStep(step - 1);
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Configuración inicial</h1>

      <div className="w-full bg-gray-200 rounded-full h-2">
        <div className="bg-blue-600 h-2 rounded-full transition-all" style={{ width: `${progress}%` }} />
      </div>
      <p className="text-sm text-gray-500">Paso {step + 1} de {TOTAL_STEPS}</p>

      {step === 0 && (
        <Card title="Datos de la empresa">
          <div className="space-y-4">
            <Input label="Nombre de la empresa" required value={company.name} onChange={(e) => setCompany({ ...company, name: e.target.value })} />
            <Input label="RIF" required value={company.rif} onChange={(e) => setCompany({ ...company, rif: e.target.value })} />
            <Input label="Dirección" value={company.address} onChange={(e) => setCompany({ ...company, address: e.target.value })} />
          </div>
        </Card>
      )}

      {step === 1 && (
        <Card title="Sucursal y almacén">
          <div className="space-y-4">
            <Input label="Nombre de la sucursal" required value={branch.branchName} onChange={(e) => setBranch({ ...branch, branchName: e.target.value })} />
            <Input label="Dirección de la sucursal" value={branch.branchAddress} onChange={(e) => setBranch({ ...branch, branchAddress: e.target.value })} />
            <Input label="Nombre del almacén" required value={branch.warehouseName} onChange={(e) => setBranch({ ...branch, warehouseName: e.target.value })} />
            <Input label="Código del almacén" value={branch.warehouseCode} onChange={(e) => setBranch({ ...branch, warehouseCode: e.target.value })} />
          </div>
        </Card>
      )}

      {step === 2 && (
        <Card title="Configuración de impuestos">
          <div className="grid grid-cols-2 gap-4">
            <Input label="IVA (%)" type="number" value={tax.iva} onChange={(e) => setTax({ ...tax, iva: e.target.value })} />
            <Input label="IGTF (%)" type="number" value={tax.igtf} onChange={(e) => setTax({ ...tax, igtf: e.target.value })} />
          </div>
        </Card>
      )}

      {step === 3 && (
        <Card title="Primer producto">
          <div className="space-y-4">
            <Input label="SKU" value={product.sku} onChange={(e) => setProduct({ ...product, sku: e.target.value })} />
            <Input label="Nombre" required value={product.name} onChange={(e) => setProduct({ ...product, name: e.target.value })} />
            <Input label="Costo" type="number" step="0.01" value={product.cost} onChange={(e) => setProduct({ ...product, cost: e.target.value })} />
          </div>
        </Card>
      )}

      {step === 4 && (
        <Card title="Dispositivo fiscal">
          <div className="space-y-4">
            <Input label="Nombre" required value={fiscal.name} onChange={(e) => setFiscal({ ...fiscal, name: e.target.value })} />
            <Select label="Tipo" value={fiscal.type} onChange={(e) => setFiscal({ ...fiscal, type: e.target.value })}>
              <option value="forma-libre">Forma libre</option>
              <option value="maquina-fiscal">Máquina fiscal</option>
              <option value="imprenta-digital">Imprenta digital</option>
            </Select>
            <Input label="Modelo" value={fiscal.model} onChange={(e) => setFiscal({ ...fiscal, model: e.target.value })} />
            <Input label="Serial" value={fiscal.serial} onChange={(e) => setFiscal({ ...fiscal, serial: e.target.value })} />
          </div>
        </Card>
      )}

      <div className="flex justify-between">
        <Button variant="secondary" onClick={back} disabled={step === 0}>Atrás</Button>
        <Button onClick={next} disabled={saving}>
          {saving ? 'Guardando...' : step === TOTAL_STEPS - 1 ? 'Finalizar' : 'Siguiente'}
        </Button>
      </div>
    </div>
  );
}
