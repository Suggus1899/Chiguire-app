'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type DeliveryRoute, type DeliveryStop, type Customer } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';
import { Modal } from '@/components/Modal';

export default function RutasPage() {
  const [routes, setRoutes] = useState<DeliveryRoute[]>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ driver_name: '', date: '' });
  const [selectedRoute, setSelectedRoute] = useState<DeliveryRoute | null>(null);
  const [stops, setStops] = useState<DeliveryStop[]>([]);
  const [stopForm, setStopForm] = useState({ customer_id: '', address: '' });

  useEffect(() => {
    Promise.all([
      api.delivery.listRoutes().catch(() => [] as DeliveryRoute[]),
      api.customers.list().catch(() => [] as Customer[]),
    ]).then(([rt, cust]) => {
      setRoutes(rt);
      setCustomers(cust);
      setLoading(false);
    });
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!form.driver_name) return;
    const r = await api.delivery.createRoute({ ...form, status: 'planned' });
    setRoutes((prev) => [r, ...prev]);
    setForm({ driver_name: '', date: '' });
    setShowForm(false);
  }

  async function openRoute(route: DeliveryRoute) {
    setSelectedRoute(route);
    try {
      const full = await api.delivery.getRoute(route.id);
      setStops((full as unknown as { stops?: DeliveryStop[] }).stops ?? []);
    } catch {
      setStops([]);
    }
  }

  async function addStop(e: FormEvent) {
    e.preventDefault();
    if (!selectedRoute || !stopForm.customer_id) return;
    const stop = await api.delivery.addStop(selectedRoute.id, {
      ...stopForm,
      sequence: stops.length + 1,
      status: 'pending',
    });
    setStops((prev) => [...prev, stop]);
    setStopForm({ customer_id: '', address: '' });
  }

  async function updateStatus(stopId: string, status: string) {
    if (!selectedRoute) return;
    const updated = await api.delivery.updateStopStatus(selectedRoute.id, stopId, status);
    setStops((prev) => prev.map((s) => (s.id === stopId ? updated : s)));
  }

  const customerName = (id: string) => customers.find((c) => c.id === id)?.name ?? id.slice(0, 8);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Rutas de entrega</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nueva ruta'}</Button>
      </div>

      {showForm && (
        <Card title="Nueva ruta">
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="Conductor" required value={form.driver_name} onChange={(e) => setForm({ ...form, driver_name: e.target.value })} />
              <Input label="Fecha" type="date" value={form.date} onChange={(e) => setForm({ ...form, date: e.target.value })} />
            </div>
            <Button type="submit">Crear</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<DeliveryRoute>
            columns={[
              { key: 'driver_name', label: 'Conductor' },
              { key: 'date', label: 'Fecha', render: (r) => r.date ? new Date(r.date).toLocaleDateString('es-VE') : '—' },
              { key: 'status', label: 'Estado' },
              {
                key: 'actions',
                label: '',
                render: (row) => (
                  <button onClick={() => openRoute(row)} className="text-blue-600 text-xs hover:underline">Ver paradas</button>
                ),
              },
            ]}
            data={routes}
            empty="Sin rutas"
          />
        )}
      </Card>

      <Modal open={!!selectedRoute} onClose={() => setSelectedRoute(null)} title={`Ruta - ${selectedRoute?.driver_name ?? ''}`}>
        <form onSubmit={addStop} className="space-y-3 mb-4 pb-4 border-b border-gray-200">
          <Select label="Cliente" required value={stopForm.customer_id} onChange={(e) => setStopForm({ ...stopForm, customer_id: e.target.value })}>
            <option value="">Seleccionar...</option>
            {customers.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </Select>
          <Input label="Dirección" value={stopForm.address} onChange={(e) => setStopForm({ ...stopForm, address: e.target.value })} />
          <Button type="submit" className="!py-1">Agregar parada</Button>
        </form>
        {stops.length === 0 ? (
          <p className="text-sm text-gray-400">Sin paradas.</p>
        ) : (
          <div className="space-y-2">
            {stops.map((stop) => (
              <div key={stop.id} className="flex justify-between items-center text-sm py-2 border-b border-gray-100">
                <div>
                  <p className="text-gray-800">{customerName(stop.customer_id)}</p>
                  <p className="text-gray-400 text-xs">{stop.address}</p>
                </div>
                <Select value={stop.status} onChange={(e) => updateStatus(stop.id, e.target.value)} className="!w-32 !py-1">
                  <option value="pending">Pendiente</option>
                  <option value="in_transit">En tránsito</option>
                  <option value="delivered">Entregado</option>
                  <option value="failed">Fallido</option>
                </Select>
              </div>
            ))}
          </div>
        )}
      </Modal>
    </div>
  );
}
