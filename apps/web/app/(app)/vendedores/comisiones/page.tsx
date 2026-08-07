'use client';
import { useState, useEffect } from 'react';
import { api, type Seller, type SellerCommission } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Modal } from '@/components/Modal';

export default function ComisionesPage() {
  const [sellers, setSellers] = useState<Seller[]>([]);
  const [commissions, setCommissions] = useState<SellerCommission[]>([]);
  const [loading, setLoading] = useState(true);
  const [sellerFilter, setSellerFilter] = useState('all');
  const [payStatusFilter, setPayStatusFilter] = useState('all');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const [selected, setSelected] = useState<string[]>([]);
  const [showPayModal, setShowPayModal] = useState(false);
  const [payForm, setPayForm] = useState({ reference: '', notes: '' });

  useEffect(() => {
    api.sellers.list()
      .then(setSellers)
      .catch((err) => console.error('operation failed:', err))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (sellerFilter !== 'all') {
      api.sellers.listCommissions(sellerFilter)
        .then(setCommissions)
        .catch((err) => { console.error('operation failed:', err); setCommissions([]); });
    } else {
      setCommissions([]);
    }
  }, [sellerFilter]);

  const filtered = commissions.filter((c) => {
    if (payStatusFilter !== 'all' && c.payment_status !== payStatusFilter) return false;
    if (dateFrom && c.created_at < dateFrom) return false;
    if (dateTo && c.created_at > dateTo + 'T23:59:59') return false;
    return true;
  });

  const selectedCommissions = filtered.filter((c) => selected.includes(c.id));
  const selectedTotal = selectedCommissions.reduce((s, c) => s + c.commission_amount, 0);
  const selectedSellerId = selectedCommissions.length > 0 ? selectedCommissions[0].seller_id : '';

  function toggle(id: string) {
    setSelected((prev) => {
      if (prev.includes(id)) return prev.filter((x) => x !== id);
      const comm = filtered.find((c) => c.id === id);
      if (!comm) return prev;
      if (prev.length > 0) {
        const firstComm = filtered.find((c) => c.id === prev[0]);
        if (firstComm && firstComm.seller_id !== comm.seller_id) {
          return [id];
        }
      }
      return [...prev, id];
    });
  }

  async function markPaid() {
    if (!selectedSellerId || selected.length === 0) return;
    await api.sellers.markCommissionsPaid(selectedSellerId, { commission_ids: selected, reference: payForm.reference, notes: payForm.notes });
    setCommissions((prev) => prev.map((c) => (selected.includes(c.id) ? { ...c, payment_status: 'paid' } : c)));
    setSelected([]);
    setShowPayModal(false);
    setPayForm({ reference: '', notes: '' });
  }

  function exportCsv() {
    const headers = ['Vendedor', 'Tipo doc', 'Monto base', 'Comisión', '%', 'Estado pago', 'CxC cliente', 'Fecha'];
    const rows = filtered.map((c) => {
      const seller = sellers.find((s) => s.id === c.seller_id);
      return [seller?.name ?? '', c.doc_type, c.base_amount, c.commission_amount, c.rate, c.payment_status, c.cxc_customer, c.created_at];
    });
    const csv = [headers, ...rows].map((r) => r.join(',')).join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'comisiones.csv';
    a.click();
    URL.revokeObjectURL(url);
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Comisiones de vendedores</h1>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={exportCsv}>Exportar</Button>
          <Button onClick={() => setShowPayModal(true)} disabled={selected.length === 0}>Marcar como Pagado</Button>
        </div>
      </div>

      <Card title="Filtros">
        <div className="grid grid-cols-4 gap-4">
          <Select label="Vendedor" value={sellerFilter} onChange={(e) => { setSellerFilter(e.target.value); setSelected([]); }}>
            <option value="all">Todos</option>
            {sellers.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
          </Select>
          <Select label="Estado pago" value={payStatusFilter} onChange={(e) => setPayStatusFilter(e.target.value)}>
            <option value="all">Todos</option>
            <option value="pending">Pendiente</option>
            <option value="paid">Pagado</option>
          </Select>
          <Input label="Desde" type="date" value={dateFrom} onChange={(e) => setDateFrom(e.target.value)} />
          <Input label="Hasta" type="date" value={dateTo} onChange={(e) => setDateTo(e.target.value)} />
        </div>
      </Card>

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : filtered.length === 0 ? (
          <p className="text-sm text-gray-400 py-4 text-center">Sin comisiones. Selecciona un vendedor.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200 text-left text-gray-500">
                  <th className="py-2 px-3 font-medium"><input type="checkbox" /></th>
                  <th className="py-2 px-3 font-medium">Vendedor</th>
                  <th className="py-2 px-3 font-medium">Tipo doc</th>
                  <th className="py-2 px-3 font-medium">Monto base</th>
                  <th className="py-2 px-3 font-medium">Comisión</th>
                  <th className="py-2 px-3 font-medium">%</th>
                  <th className="py-2 px-3 font-medium">Estado pago</th>
                  <th className="py-2 px-3 font-medium">CxC cliente</th>
                  <th className="py-2 px-3 font-medium">Fecha</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((c) => {
                  const seller = sellers.find((s) => s.id === c.seller_id);
                  return (
                    <tr key={c.id} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-2 px-3"><input type="checkbox" checked={selected.includes(c.id)} onChange={() => toggle(c.id)} /></td>
                      <td className="py-2 px-3 text-gray-800">{seller?.name ?? '—'}</td>
                      <td className="py-2 px-3 text-gray-800">{c.doc_type}</td>
                      <td className="py-2 px-3 text-gray-800">{c.base_amount.toFixed(2)}</td>
                      <td className="py-2 px-3 text-gray-800">{c.commission_amount.toFixed(2)}</td>
                      <td className="py-2 px-3 text-gray-800">{c.rate}%</td>
                      <td className="py-2 px-3 text-gray-800">{c.payment_status}</td>
                      <td className="py-2 px-3 text-gray-800">{c.cxc_customer}</td>
                      <td className="py-2 px-3 text-gray-800">{c.created_at.slice(0, 10)}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      <Modal open={showPayModal} onClose={() => setShowPayModal(false)} title="Marcar comisiones como pagadas">
        <div className="space-y-4">
          <p className="text-sm text-gray-600">Total seleccionado: <strong>{selectedTotal.toFixed(2)}</strong></p>
          <Input label="Referencia" value={payForm.reference} onChange={(e) => setPayForm({ ...payForm, reference: e.target.value })} />
          <Input label="Notas" value={payForm.notes} onChange={(e) => setPayForm({ ...payForm, notes: e.target.value })} />
          <Button onClick={markPaid}>Confirmar pago</Button>
        </div>
      </Modal>
    </div>
  );
}
