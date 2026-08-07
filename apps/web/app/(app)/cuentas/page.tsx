'use client';
import { useState, useEffect } from 'react';
import { api, type AccountItem, type PaymentMethod } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Modal } from '@/components/Modal';

export default function CuentasPage() {
  const [tab, setTab] = useState<'receivable' | 'payable'>('receivable');
  const [items, setItems] = useState<AccountItem[]>([]);
  const [methods, setMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState('all');
  const [payModal, setPayModal] = useState<AccountItem | null>(null);
  const [payForm, setPayForm] = useState({ amount: '', payment_method: '', reference: '', notes: '' });
  const [reminderModal, setReminderModal] = useState<AccountItem | null>(null);
  const [reminderChannel, setReminderChannel] = useState<'sms' | 'email'>('sms');

  useEffect(() => {
    api.paymentMethods.list().catch(() => [] as PaymentMethod[]).then(setMethods);
  }, []);

  useEffect(() => {
    setLoading(true);
    const fetcher = tab === 'receivable' ? api.accountsPayable.listReceivable : api.accountsPayable.listPayable;
    fetcher().then(setItems).catch(() => setItems([])).finally(() => setLoading(false));
  }, [tab]);

  const filtered = items.filter((i) => statusFilter === 'all' ? true : i.status === statusFilter);

  function openPay(item: AccountItem) {
    setPayModal(item);
    setPayForm({ amount: String(item.balance), payment_method: '', reference: '', notes: '' });
  }

  async function submitPay() {
    if (!payModal) return;
    await api.accountsPayable.createPayment({ account_id: payModal.id, amount: parseFloat(payForm.amount) || 0, payment_method: payForm.payment_method, reference: payForm.reference, notes: payForm.notes });
    setItems((prev) => prev.map((i) => (i.id === payModal.id ? { ...i, balance: i.balance - (parseFloat(payForm.amount) || 0) } : i)));
    setPayModal(null);
  }

  async function sendReminder() {
    if (!reminderModal) return;
    await api.accountsPayable.sendReminder(reminderModal.id, { channel: reminderChannel });
    setReminderModal(null);
    alert('Recordatorio enviado');
  }

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Cuentas por Cobrar/Pagar</h1>

      <div className="flex gap-2">
        <button onClick={() => setTab('receivable')} className={`px-4 py-2 text-sm font-medium rounded-lg ${tab === 'receivable' ? 'bg-blue-600 text-white' : 'bg-white border border-gray-300 text-gray-700'}`}>Por Cobrar</button>
        <button onClick={() => setTab('payable')} className={`px-4 py-2 text-sm font-medium rounded-lg ${tab === 'payable' ? 'bg-blue-600 text-white' : 'bg-white border border-gray-300 text-gray-700'}`}>Por Pagar</button>
      </div>

      <div className="flex gap-3 items-end">
        <Select label="Estado" value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="max-w-xs">
          <option value="all">Todos</option>
          <option value="open">Abierta</option>
          <option value="partial">Parcial</option>
          <option value="paid">Pagada</option>
          <option value="overdue">Vencida</option>
        </Select>
      </div>

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : filtered.length === 0 ? (
          <p className="text-sm text-gray-400 py-4 text-center">Sin cuentas registradas</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200 text-left text-gray-500">
                  <th className="py-2 px-3 font-medium">{tab === 'receivable' ? 'Cliente' : 'Proveedor'}</th>
                  <th className="py-2 px-3 font-medium">Monto</th>
                  <th className="py-2 px-3 font-medium">Balance</th>
                  <th className="py-2 px-3 font-medium">Estado</th>
                  <th className="py-2 px-3 font-medium">Vencimiento</th>
                  <th className="py-2 px-3 font-medium">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((i) => (
                  <tr key={i.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-2 px-3 text-gray-800">{i.party_name}</td>
                    <td className="py-2 px-3 text-gray-800">{i.amount.toFixed(2)}</td>
                    <td className="py-2 px-3 text-gray-800">{i.balance.toFixed(2)}</td>
                    <td className="py-2 px-3 text-gray-800">{i.status}</td>
                    <td className="py-2 px-3 text-gray-800">{i.due_date.slice(0, 10)}</td>
                    <td className="py-2 px-3">
                      <div className="flex gap-2">
                        <button onClick={() => openPay(i)} className="text-blue-600 text-xs hover:underline">Agregar Pago</button>
                        <button onClick={() => setReminderModal(i)} className="text-gray-600 text-xs hover:underline">Recordatorio</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      <Modal open={!!payModal} onClose={() => setPayModal(null)} title="Agregar pago">
        <div className="space-y-4">
          <Input label="Monto" type="number" step="0.01" value={payForm.amount} onChange={(e) => setPayForm({ ...payForm, amount: e.target.value })} />
          <Select label="Método de pago" value={payForm.payment_method} onChange={(e) => setPayForm({ ...payForm, payment_method: e.target.value })}>
            <option value="">Seleccionar...</option>
            {methods.map((m) => <option key={m.id} value={m.id}>{m.name}</option>)}
          </Select>
          <Input label="Referencia" value={payForm.reference} onChange={(e) => setPayForm({ ...payForm, reference: e.target.value })} />
          <Input label="Notas" value={payForm.notes} onChange={(e) => setPayForm({ ...payForm, notes: e.target.value })} />
          <Button onClick={submitPay}>Guardar pago</Button>
        </div>
      </Modal>

      <Modal open={!!reminderModal} onClose={() => setReminderModal(null)} title="Enviar recordatorio">
        <div className="space-y-4">
          <Select label="Canal" value={reminderChannel} onChange={(e) => setReminderChannel(e.target.value as 'sms' | 'email')}>
            <option value="sms">SMS</option>
            <option value="email">Email</option>
          </Select>
          <Button onClick={sendReminder}>Enviar</Button>
        </div>
      </Modal>
    </div>
  );
}
