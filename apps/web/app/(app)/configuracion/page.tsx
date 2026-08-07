'use client';
import { useState, useEffect, FormEvent, ChangeEvent } from 'react';
import { api, type Subscription, type OutgoingWebhook, type WebhookDelivery } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function ConfiguracionPage() {
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [webhooks, setWebhooks] = useState<OutgoingWebhook[]>([]);
  const [deliveries, setDeliveries] = useState<WebhookDelivery[]>([]);
  const [loading, setLoading] = useState(true);
  const [whForm, setWhForm] = useState({ url: '', event: 'invoice.created', secret: '' });
  const [importMsg, setImportMsg] = useState('');

  useEffect(() => {
    Promise.all([
      api.saas.getSubscription().catch(() => null),
      api.webhooks.list().catch(() => [] as OutgoingWebhook[]),
      api.webhooks.listDeliveries().catch(() => [] as WebhookDelivery[]),
    ]).then(([sub, wh, del]) => {
      setSubscription(sub);
      setWebhooks(wh);
      setDeliveries(del);
      setLoading(false);
    });
  }, []);

  async function updatePlan(plan: string) {
    const sub = await api.saas.updatePlan(plan);
    setSubscription(sub);
  }

  async function addWebhook(e: FormEvent) {
    e.preventDefault();
    if (!whForm.url) return;
    const wh = await api.webhooks.create({ ...whForm, is_active: true });
    setWebhooks((prev) => [...prev, wh]);
    setWhForm({ url: '', event: 'invoice.created', secret: '' });
  }

  async function removeWebhook(id: string) {
    await api.webhooks.delete(id);
    setWebhooks((prev) => prev.filter((w) => w.id !== id));
  }

  async function testWebhook(id: string) {
    try {
      const result = await api.webhooks.test(id);
      alert(`Test: ${result.success ? 'OK' : 'Falló'} (HTTP ${result.status})`);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Error');
    }
  }

  async function handleImport(e: ChangeEvent<HTMLInputElement>, type: 'customers' | 'products') {
    const file = e.target.files?.[0];
    if (!file) return;
    try {
      const result = type === 'customers'
        ? await api.importExport.importCustomers(file)
        : await api.importExport.importProducts(file);
      setImportMsg(`Importados: ${result.imported}`);
    } catch (err) {
      setImportMsg(err instanceof Error ? err.message : 'Error');
    }
  }

  async function handleExport(type: 'invoices' | 'customers') {
    try {
      const blob = type === 'invoices'
        ? await api.importExport.exportInvoices()
        : await api.importExport.exportCustomers();
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${type}.csv`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Error');
    }
  }

  if (loading) return <p className="text-gray-500">Cargando...</p>;

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Configuración</h1>

      <Card title="Suscripción y plan">
        {subscription ? (
          <div className="space-y-3">
            <dl className="text-sm space-y-1">
              <div className="flex justify-between"><dt className="text-gray-500">Plan actual:</dt><dd className="font-semibold">{subscription.plan}</dd></div>
              <div className="flex justify-between"><dt className="text-gray-500">Estado:</dt><dd>{subscription.status}</dd></div>
              <div className="flex justify-between"><dt className="text-gray-500">Usuarios:</dt><dd>{subscription.seats}</dd></div>
              <div className="flex justify-between"><dt className="text-gray-500">Renovación:</dt><dd>{subscription.current_period_end ? new Date(subscription.current_period_end).toLocaleDateString('es-VE') : '—'}</dd></div>
            </dl>
            <div className="flex gap-2">
              <Button variant="secondary" onClick={() => updatePlan('starter')}>Starter</Button>
              <Button variant="secondary" onClick={() => updatePlan('pro')}>Pro</Button>
              <Button variant="secondary" onClick={() => updatePlan('enterprise')}>Enterprise</Button>
            </div>
          </div>
        ) : (
          <p className="text-sm text-gray-400">Sin suscripción activa.</p>
        )}
      </Card>

      <Card title="Webhooks">
        <form onSubmit={addWebhook} className="grid grid-cols-3 gap-3 mb-4">
          <Input label="URL" required value={whForm.url} onChange={(e) => setWhForm({ ...whForm, url: e.target.value })} />
          <Select label="Evento" value={whForm.event} onChange={(e) => setWhForm({ ...whForm, event: e.target.value })}>
            <option value="invoice.created">Factura creada</option>
            <option value="invoice.emitted">Factura emitida</option>
            <option value="payment.received">Pago recibido</option>
            <option value="stock.movement">Movimiento de stock</option>
          </Select>
          <Input label="Secret" value={whForm.secret} onChange={(e) => setWhForm({ ...whForm, secret: e.target.value })} />
          <div className="col-span-3">
            <Button type="submit">Agregar webhook</Button>
          </div>
        </form>
        <Table<OutgoingWebhook>
          columns={[
            { key: 'url', label: 'URL' },
            { key: 'event', label: 'Evento' },
            { key: 'is_active', label: 'Activo', render: (r) => (r.is_active ? 'Sí' : 'No') },
            {
              key: 'actions',
              label: '',
              render: (row) => (
                <div className="flex gap-2">
                  <button onClick={() => testWebhook(row.id)} className="text-blue-600 text-xs hover:underline">Probar</button>
                  <button onClick={() => removeWebhook(row.id)} className="text-red-600 text-xs hover:underline">Eliminar</button>
                </div>
              ),
            },
          ]}
          data={webhooks}
          empty="Sin webhooks"
        />
      </Card>

      {deliveries.length > 0 && (
        <Card title="Entregas de webhooks">
          <Table<WebhookDelivery>
            columns={[
              { key: 'event', label: 'Evento' },
              { key: 'status', label: 'HTTP' },
              { key: 'attempt', label: 'Intento' },
              { key: 'created_at', label: 'Fecha', render: (r) => new Date(r.created_at).toLocaleString('es-VE') },
            ]}
            data={deliveries}
            empty=""
          />
        </Card>
      )}

      <Card title="Importar / Exportar">
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-3">
            <h3 className="text-sm font-medium text-gray-700">Importar</h3>
            <div>
              <label className="block text-sm text-gray-500 mb-1">Clientes (CSV)</label>
              <input type="file" accept=".csv" onChange={(e) => handleImport(e, 'customers')} className="text-sm" />
            </div>
            <div>
              <label className="block text-sm text-gray-500 mb-1">Productos (CSV)</label>
              <input type="file" accept=".csv" onChange={(e) => handleImport(e, 'products')} className="text-sm" />
            </div>
            {importMsg && <p className="text-sm text-blue-600">{importMsg}</p>}
          </div>
          <div className="space-y-3">
            <h3 className="text-sm font-medium text-gray-700">Exportar</h3>
            <div className="flex gap-2">
              <Button variant="secondary" onClick={() => handleExport('invoices')}>Exportar facturas</Button>
              <Button variant="secondary" onClick={() => handleExport('customers')}>Exportar clientes</Button>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}
