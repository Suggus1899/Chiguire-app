'use client';
import { useState, useEffect } from 'react';
import { api, type Product } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';

type ReportKey = 'salesBook' | 'purchasesBook' | 'inventoryCurrent' | 'inventoryValued' | 'kardex' | 'art177' | 'igtfReport';

const reportConfig: { key: ReportKey; label: string; needsDates: boolean; needsProduct: boolean }[] = [
  { key: 'salesBook', label: 'Libro de Ventas', needsDates: true, needsProduct: false },
  { key: 'purchasesBook', label: 'Libro de Compras', needsDates: true, needsProduct: false },
  { key: 'inventoryCurrent', label: 'Inventario Actual', needsDates: false, needsProduct: false },
  { key: 'inventoryValued', label: 'Inventario Valorizado', needsDates: false, needsProduct: false },
  { key: 'kardex', label: 'Kardex', needsDates: true, needsProduct: true },
  { key: 'art177', label: 'Artículo 177 (ISLR)', needsDates: true, needsProduct: false },
  { key: 'igtfReport', label: 'Reporte IGTF', needsDates: true, needsProduct: false },
];

export default function InformesPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [results, setResults] = useState<Record<string, Record<string, unknown>[] | null>>({});
  const [loading, setLoading] = useState<string | null>(null);
  const [params, setParams] = useState<Record<string, { from: string; to: string; productId: string }>>({});

  useEffect(() => {
    api.products.list().catch(() => [] as Product[]).then(setProducts);
  }, []);

  async function generate(r: { key: ReportKey; needsDates: boolean; needsProduct: boolean }) {
    setLoading(r.key);
    const p = params[r.key] ?? { from: '', to: '', productId: '' };
    try {
      let data: Record<string, unknown>[];
      if (r.key === 'salesBook') data = await api.reports.salesBook(p.from, p.to);
      else if (r.key === 'purchasesBook') data = await api.reports.purchasesBook(p.from, p.to);
      else if (r.key === 'inventoryCurrent') data = await api.reports.inventoryCurrent();
      else if (r.key === 'inventoryValued') data = await api.reports.inventoryValued();
      else if (r.key === 'kardex') data = await api.reports.kardex(p.productId, p.from, p.to);
      else if (r.key === 'art177') data = await api.reports.art177(p.from, p.to);
      else data = await api.reports.igtfReport(p.from, p.to);
      setResults((prev) => ({ ...prev, [r.key]: data }));
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Error al generar');
    } finally {
      setLoading(null);
    }
  }

  function exportExcel(key: string) {
    const data = results[key];
    if (!data || data.length === 0) return;
    const headers = Object.keys(data[0]);
    const rows = data.map((r) => headers.map((h) => String(r[h] ?? '')).join(','));
    const csv = [headers.join(','), ...rows].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${key}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  }

  function setParam(key: string, field: 'from' | 'to' | 'productId', value: string) {
    setParams((prev) => ({ ...prev, [key]: { ...prev[key] ?? { from: '', to: '', productId: '' }, [field]: value } }));
  }

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Informes</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {reportConfig.map((r) => {
          const data = results[r.key];
          const p = params[r.key] ?? { from: '', to: '', productId: '' };
          return (
            <Card key={r.key} title={r.label}>
              <div className="space-y-3">
                <div className="grid grid-cols-2 gap-3">
                  {r.needsDates && (
                    <>
                      <Input label="Desde" type="date" value={p.from} onChange={(e) => setParam(r.key, 'from', e.target.value)} />
                      <Input label="Hasta" type="date" value={p.to} onChange={(e) => setParam(r.key, 'to', e.target.value)} />
                    </>
                  )}
                  {r.needsProduct && (
                    <div className="col-span-2">
                      <Select label="Producto" value={p.productId} onChange={(e) => setParam(r.key, 'productId', e.target.value)}>
                        <option value="">Seleccionar...</option>
                        {products.map((prod) => <option key={prod.id} value={prod.id}>{prod.name}</option>)}
                      </Select>
                    </div>
                  )}
                </div>
                <div className="flex gap-2">
                  <Button onClick={() => generate(r)} disabled={loading === r.key}>
                    {loading === r.key ? 'Generando...' : 'Generar'}
                  </Button>
                  {data && data.length > 0 && (
                    <Button variant="secondary" onClick={() => exportExcel(r.key)}>Exportar Excel</Button>
                  )}
                </div>
                {data && data.length > 0 && (
                  <div className="overflow-x-auto max-h-64">
                    <table className="w-full text-xs">
                      <thead>
                        <tr className="border-b border-gray-200 text-left text-gray-500">
                          {Object.keys(data[0]).map((h) => <th key={h} className="py-1 px-2 font-medium">{h}</th>)}
                        </tr>
                      </thead>
                      <tbody>
                        {data.slice(0, 20).map((row, i) => (
                          <tr key={i} className="border-b border-gray-100">
                            {Object.values(row).map((v, j) => <td key={j} className="py-1 px-2 text-gray-800">{String(v ?? '')}</td>)}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    {data.length > 20 && <p className="text-xs text-gray-400 mt-1">Mostrando 20 de {data.length} filas</p>}
                  </div>
                )}
              </div>
            </Card>
          );
        })}
      </div>
    </div>
  );
}
