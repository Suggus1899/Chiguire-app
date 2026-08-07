'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type TaxCategory, type Withholding, type FiscalBook } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function FiscalPage() {
  const [taxCategories, setTaxCategories] = useState<TaxCategory[]>([]);
  const [withholdings, setWithholdings] = useState<Withholding[]>([]);
  const [books, setBooks] = useState<FiscalBook[]>([]);
  const [rate, setRate] = useState<{ currency: string; rate: number; date: string } | null>(null);
  const [newRate, setNewRate] = useState('');
  const [whForm, setWhForm] = useState({ invoice_id: '', type: 'IVA', base_amount: '0', rate: '75' });
  const [bookPeriod, setBookPeriod] = useState('');

  useEffect(() => {
    Promise.all([
      api.fiscal.listTaxCategories().catch(() => [] as TaxCategory[]),
      api.fiscal.listWithholdings().catch(() => [] as Withholding[]),
      api.fiscal.getExchangeRate('USD').catch(() => null),
    ]).then(([tc, wh, r]) => {
      setTaxCategories(tc);
      setWithholdings(wh);
      if (r) setRate(r);
    });
  }, []);

  async function setExchangeRate(e: FormEvent) {
    e.preventDefault();
    const r = await api.fiscal.setExchangeRate('USD', parseFloat(newRate) || 0);
    setRate({ currency: 'USD', rate: r.rate, date: new Date().toISOString() });
    setNewRate('');
  }

  async function createWithholding(e: FormEvent) {
    e.preventDefault();
    if (!whForm.invoice_id) return;
    const wh = await api.fiscal.createWithholding({
      ...whForm,
      base_amount: parseFloat(whForm.base_amount) || 0,
      rate: parseFloat(whForm.rate) || 0,
      amount: (parseFloat(whForm.base_amount) || 0) * (parseFloat(whForm.rate) || 0) / 100,
    });
    setWithholdings((prev) => [wh, ...prev]);
    setWhForm({ invoice_id: '', type: 'IVA', base_amount: '0', rate: '75' });
  }

  async function generateBook(e: FormEvent) {
    e.preventDefault();
    if (!bookPeriod) return;
    const result = await api.fiscal.generateFiscalBook(bookPeriod);
    const book = await api.fiscal.getFiscalBook(result.id);
    setBooks((prev) => [book, ...prev]);
    setBookPeriod('');
  }

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Fiscal</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card title="Tasa de cambio (USD)">
          {rate && (
            <dl className="text-sm space-y-1 mb-4">
              <div className="flex justify-between"><dt className="text-gray-500">Moneda:</dt><dd>{rate.currency}</dd></div>
              <div className="flex justify-between"><dt className="text-gray-500">Tasa BCV:</dt><dd className="font-semibold">{rate.rate.toFixed(2)}</dd></div>
              <div className="flex justify-between"><dt className="text-gray-500">Fecha:</dt><dd>{new Date(rate.date).toLocaleDateString('es-VE')}</dd></div>
            </dl>
          )}
          <form onSubmit={setExchangeRate} className="flex gap-2">
            <Input type="number" step="0.01" placeholder="Nueva tasa" value={newRate} onChange={(e) => setNewRate(e.target.value)} />
            <Button type="submit">Actualizar</Button>
          </form>
        </Card>

        <Card title="Categorías de impuesto">
          {taxCategories.length === 0 ? (
            <p className="text-sm text-gray-400">Sin categorías configuradas.</p>
          ) : (
            <Table<TaxCategory>
              columns={[
                { key: 'name', label: 'Nombre' },
                { key: 'rate', label: 'Tasa', render: (r) => `${r.rate}%` },
                { key: 'is_retention', label: 'Retención', render: (r) => (r.is_retention ? 'Sí' : 'No') },
              ]}
              data={taxCategories}
              empty=""
            />
          )}
        </Card>
      </div>

      <Card title="Retenciones">
        <form onSubmit={createWithholding} className="grid grid-cols-4 gap-3 mb-4">
          <Input label="ID Factura" required value={whForm.invoice_id} onChange={(e) => setWhForm({ ...whForm, invoice_id: e.target.value })} />
          <Select label="Tipo" value={whForm.type} onChange={(e) => setWhForm({ ...whForm, type: e.target.value })}>
            <option value="IVA">IVA</option>
            <option value="ISLR">ISLR</option>
          </Select>
          <Input label="Base" type="number" step="0.01" value={whForm.base_amount} onChange={(e) => setWhForm({ ...whForm, base_amount: e.target.value })} />
          <Input label="Tasa %" type="number" value={whForm.rate} onChange={(e) => setWhForm({ ...whForm, rate: e.target.value })} />
          <div className="col-span-4">
            <Button type="submit">Crear retención</Button>
          </div>
        </form>
        <Table<Withholding>
          columns={[
            { key: 'type', label: 'Tipo' },
            { key: 'base_amount', label: 'Base', render: (r) => (r.base_amount ?? 0).toFixed(2) },
            { key: 'rate', label: 'Tasa', render: (r) => `${r.rate}%` },
            { key: 'amount', label: 'Monto', render: (r) => (r.amount ?? 0).toFixed(2) },
          ]}
          data={withholdings}
          empty="Sin retenciones"
        />
      </Card>

      <Card title="Libros fiscales">
        <form onSubmit={generateBook} className="flex gap-2 mb-4">
          <Input label="Período (ej. 2025-01)" value={bookPeriod} onChange={(e) => setBookPeriod(e.target.value)} />
          <Button type="submit">Generar</Button>
        </form>
        {books.length > 0 && (
          <Table<FiscalBook>
            columns={[
              { key: 'period', label: 'Período' },
              { key: 'book_type', label: 'Tipo' },
              { key: 'status', label: 'Estado' },
            ]}
            data={books}
            empty=""
          />
        )}
      </Card>
    </div>
  );
}
