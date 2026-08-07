'use client';
import { useState, useEffect, FormEvent } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { api, type Invoice, type InvoiceItem, type InvoicePayment, type Customer } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function FacturaDetallePage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const id = params.id;
  const [invoice, setInvoice] = useState<Invoice | null>(null);
  const [items, setItems] = useState<InvoiceItem[]>([]);
  const [payments, setPayments] = useState<InvoicePayment[]>([]);
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [loading, setLoading] = useState(true);
  const [showPayment, setShowPayment] = useState(false);
  const [payForm, setPayForm] = useState({ amount: '0', currency: 'VES', payment_method: 'cash', reference: '' });

  useEffect(() => {
    api.invoices.get(id)
      .then((inv) => {
        setInvoice(inv);
        if (inv.customer_id) {
          api.customers.get(inv.customer_id).then(setCustomer).catch((err) => console.error('operation failed:', err));
        }
        return inv;
      })
      .catch((err) => console.error('operation failed:', err))
      .finally(() => setLoading(false));
  }, [id]);

  async function emit() {
    const inv = await api.invoices.emit(id);
    setInvoice(inv);
  }

  async function addPayment(e: FormEvent) {
    e.preventDefault();
    const p = await api.invoices.addPayment(id, {
      ...payForm,
      amount: parseFloat(payForm.amount) || 0,
    });
    setPayments((prev) => [...prev, p]);
    setPayForm({ amount: '0', currency: 'VES', payment_method: 'cash', reference: '' });
    setShowPayment(false);
  }

  if (loading) return <p className="text-gray-500">Cargando...</p>;
  if (!invoice) return <p className="text-gray-500">Factura no encontrada</p>;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Factura {invoice.number || invoice.id.slice(0, 8)}</h1>
        <Button variant="secondary" onClick={() => router.push('/facturas')}>Volver</Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card title="Información">
          <dl className="text-sm space-y-1">
            <div className="flex justify-between"><dt className="text-gray-500">Cliente:</dt><dd>{customer?.name ?? invoice.customer_id}</dd></div>
            <div className="flex justify-between"><dt className="text-gray-500">Estado:</dt><dd>{invoice.status}</dd></div>
            <div className="flex justify-between"><dt className="text-gray-500">Moneda:</dt><dd>{invoice.currency}</dd></div>
            <div className="flex justify-between"><dt className="text-gray-500">Subtotal:</dt><dd>{(invoice.subtotal ?? 0).toFixed(2)}</dd></div>
            <div className="flex justify-between"><dt className="text-gray-500">IVA:</dt><dd>{(invoice.tax ?? 0).toFixed(2)}</dd></div>
            <div className="flex justify-between font-semibold"><dt>Total:</dt><dd>{(invoice.total ?? 0).toFixed(2)} {invoice.currency}</dd></div>
          </dl>
          {invoice.status === 'draft' && (
            <div className="mt-4">
              <Button onClick={emit}>Emitir factura</Button>
            </div>
          )}
        </Card>

        <Card title="Pagos" actions={<Button variant="secondary" onClick={() => setShowPayment(!showPayment)}>Agregar pago</Button>}>
          {showPayment && (
            <form onSubmit={addPayment} className="space-y-3 mb-4 pb-4 border-b border-gray-200">
              <div className="grid grid-cols-2 gap-3">
                <Input label="Monto" type="number" step="0.01" required value={payForm.amount} onChange={(e) => setPayForm({ ...payForm, amount: e.target.value })} />
                <Select label="Moneda" value={payForm.currency} onChange={(e) => setPayForm({ ...payForm, currency: e.target.value })}>
                  <option value="VES">VES</option>
                  <option value="USD">USD</option>
                </Select>
                <Select label="Método" value={payForm.payment_method} onChange={(e) => setPayForm({ ...payForm, payment_method: e.target.value })}>
                  <option value="cash">Efectivo</option>
                  <option value="transfer">Transferencia</option>
                  <option value="card">Tarjeta</option>
                  <option value="mobile">Pago móvil</option>
                </Select>
                <Input label="Referencia" value={payForm.reference} onChange={(e) => setPayForm({ ...payForm, reference: e.target.value })} />
              </div>
              <Button type="submit">Registrar pago</Button>
            </form>
          )}
          {payments.length === 0 ? (
            <p className="text-sm text-gray-400">Sin pagos registrados.</p>
          ) : (
            <Table<InvoicePayment>
              columns={[
                { key: 'amount', label: 'Monto', render: (r) => (r.amount ?? 0).toFixed(2) },
                { key: 'currency', label: 'Moneda' },
                { key: 'payment_method', label: 'Método' },
              ]}
              data={payments}
              empty=""
            />
          )}
        </Card>
      </div>

      <Card title="Items">
        {items.length === 0 ? (
          <p className="text-sm text-gray-400">Sin items.</p>
        ) : (
          <Table<InvoiceItem>
            columns={[
              { key: 'description', label: 'Descripción' },
              { key: 'quantity', label: 'Cantidad' },
              { key: 'unit_price', label: 'Precio', render: (r) => (r.unit_price ?? 0).toFixed(2) },
              { key: 'line_total', label: 'Total', render: (r) => (r.line_total ?? 0).toFixed(2) },
            ]}
            data={items}
            empty=""
          />
        )}
      </Card>
    </div>
  );
}
