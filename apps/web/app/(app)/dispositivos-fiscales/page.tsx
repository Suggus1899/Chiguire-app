'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type FiscalDevice, type DocumentSequence, type ContingencyBook, type Branch } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';
import { Table } from '@/components/Table';

export default function DispositivosFiscalesPage() {
  const [devices, setDevices] = useState<FiscalDevice[]>([]);
  const [branches, setBranches] = useState<Branch[]>([]);
  const [sequences, setSequences] = useState<DocumentSequence[]>([]);
  const [contingencies, setContingencies] = useState<ContingencyBook[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedDevice, setSelectedDevice] = useState('');
  const [showDeviceForm, setShowDeviceForm] = useState(false);
  const [showSeqForm, setShowSeqForm] = useState(false);
  const [showContForm, setShowContForm] = useState(false);
  const [deviceForm, setDeviceForm] = useState({ name: '', type: 'forma-libre', model: '', serial: '', branch_id: '' });
  const [seqForm, setSeqForm] = useState({ doc_type: 'invoice', prefix: '', suffix: '', last_seq: '0' });
  const [contForm, setContForm] = useState({ doc_type: 'invoice', prefix: '', start_seq: '1', end_seq: '100' });

  useEffect(() => {
    Promise.all([
      api.fiscalDevices.list().catch((err) => { console.error('operation failed:', err); return [] as FiscalDevice[]; }),
      api.branches.list().catch((err) => { console.error('operation failed:', err); return [] as Branch[]; }),
    ]).then(([devs, brs]) => {
      setDevices(devs);
      setBranches(brs);
      setLoading(false);
    });
  }, []);

  useEffect(() => {
    if (!selectedDevice) return;
    Promise.all([
      api.fiscalDevices.listSequences(selectedDevice).catch((err) => { console.error('operation failed:', err); return [] as DocumentSequence[]; }),
      api.fiscalDevices.listContingency(selectedDevice).catch((err) => { console.error('operation failed:', err); return [] as ContingencyBook[]; }),
    ]).then(([seqs, conts]) => {
      setSequences(seqs);
      setContingencies(conts);
    });
  }, [selectedDevice]);

  async function onDeviceSubmit(e: FormEvent) {
    e.preventDefault();
    const d = await api.fiscalDevices.create(deviceForm);
    setDevices((prev) => [d, ...prev]);
    setDeviceForm({ name: '', type: 'forma-libre', model: '', serial: '', branch_id: '' });
    setShowDeviceForm(false);
  }

  async function onSeqSubmit(e: FormEvent) {
    e.preventDefault();
    if (!selectedDevice) return;
    const s = await api.fiscalDevices.createSequence(selectedDevice, { ...seqForm, last_seq: parseInt(seqForm.last_seq) || 0 });
    setSequences((prev) => [s, ...prev]);
    setSeqForm({ doc_type: 'invoice', prefix: '', suffix: '', last_seq: '0' });
    setShowSeqForm(false);
  }

  async function onContSubmit(e: FormEvent) {
    e.preventDefault();
    if (!selectedDevice) return;
    const c = await api.fiscalDevices.createContingency(selectedDevice, { ...contForm, start_seq: parseInt(contForm.start_seq) || 0, end_seq: parseInt(contForm.end_seq) || 0 });
    setContingencies((prev) => [c, ...prev]);
    setContForm({ doc_type: 'invoice', prefix: '', start_seq: '1', end_seq: '100' });
    setShowContForm(false);
  }

  async function removeDevice(id: string) {
    await api.fiscalDevices.delete(id);
    setDevices((prev) => prev.filter((d) => d.id !== id));
    if (selectedDevice === id) setSelectedDevice('');
  }

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Dispositivos Fiscales</h1>

      <Card title="Dispositivos" actions={<Button variant="secondary" onClick={() => setShowDeviceForm(!showDeviceForm)}>{showDeviceForm ? 'Cancelar' : 'Nuevo'}</Button>}>
        {showDeviceForm && (
          <form onSubmit={onDeviceSubmit} className="space-y-4 mb-4">
            <div className="grid grid-cols-2 gap-4">
              <Input label="Nombre" required value={deviceForm.name} onChange={(e) => setDeviceForm({ ...deviceForm, name: e.target.value })} />
              <Select label="Tipo" value={deviceForm.type} onChange={(e) => setDeviceForm({ ...deviceForm, type: e.target.value })}>
                <option value="forma-libre">Forma libre</option>
                <option value="maquina-fiscal">Máquina fiscal</option>
                <option value="imprenta-digital">Imprenta digital</option>
              </Select>
              <Input label="Modelo" value={deviceForm.model} onChange={(e) => setDeviceForm({ ...deviceForm, model: e.target.value })} />
              <Input label="Serial" value={deviceForm.serial} onChange={(e) => setDeviceForm({ ...deviceForm, serial: e.target.value })} />
              <Select label="Sucursal" value={deviceForm.branch_id} onChange={(e) => setDeviceForm({ ...deviceForm, branch_id: e.target.value })}>
                <option value="">Sin sucursal</option>
                {branches.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
              </Select>
            </div>
            <Button type="submit">Guardar</Button>
          </form>
        )}
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<FiscalDevice>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'type', label: 'Tipo' },
              { key: 'model', label: 'Modelo' },
              { key: 'serial', label: 'Serial' },
              { key: 'status', label: 'Estado' },
              { key: 'actions', label: '', render: (row) => (
                <div className="flex gap-2">
                  <button onClick={() => setSelectedDevice(row.id)} className="text-blue-600 text-xs hover:underline">Secuencias</button>
                  <button onClick={() => removeDevice(row.id)} className="text-red-600 text-xs hover:underline">Eliminar</button>
                </div>
              ) },
            ]}
            data={devices}
            empty="Sin dispositivos registrados"
          />
        )}
      </Card>

      {selectedDevice && (
        <>
          <Card title="Secuencias de documentos" actions={<Button variant="secondary" onClick={() => setShowSeqForm(!showSeqForm)}>{showSeqForm ? 'Cancelar' : 'Nueva'}</Button>}>
            {showSeqForm && (
              <form onSubmit={onSeqSubmit} className="space-y-4 mb-4">
                <div className="grid grid-cols-4 gap-4">
                  <Select label="Tipo doc" value={seqForm.doc_type} onChange={(e) => setSeqForm({ ...seqForm, doc_type: e.target.value })}>
                    <option value="invoice">Factura</option>
                    <option value="credit_note">Nota de crédito</option>
                    <option value="debit_note">Nota de débito</option>
                  </Select>
                  <Input label="Prefijo" value={seqForm.prefix} onChange={(e) => setSeqForm({ ...seqForm, prefix: e.target.value })} />
                  <Input label="Sufijo" value={seqForm.suffix} onChange={(e) => setSeqForm({ ...seqForm, suffix: e.target.value })} />
                  <Input label="Último n°" type="number" value={seqForm.last_seq} onChange={(e) => setSeqForm({ ...seqForm, last_seq: e.target.value })} />
                </div>
                <Button type="submit">Guardar</Button>
              </form>
            )}
            {sequences.length === 0 ? (
              <p className="text-sm text-gray-400 py-4 text-center">Sin secuencias</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead><tr className="border-b border-gray-200 text-left text-gray-500">
                    <th className="py-2 px-3 font-medium">Tipo doc</th>
                    <th className="py-2 px-3 font-medium">Prefijo</th>
                    <th className="py-2 px-3 font-medium">Sufijo</th>
                    <th className="py-2 px-3 font-medium">Último n°</th>
                  </tr></thead>
                  <tbody>
                    {sequences.map((s) => (
                      <tr key={s.id} className="border-b border-gray-100">
                        <td className="py-2 px-3 text-gray-800">{s.doc_type}</td>
                        <td className="py-2 px-3 text-gray-800">{s.prefix}</td>
                        <td className="py-2 px-3 text-gray-800">{s.suffix}</td>
                        <td className="py-2 px-3 text-gray-800">{s.last_seq}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </Card>

          <Card title="Libros de contingencia" actions={<Button variant="secondary" onClick={() => setShowContForm(!showContForm)}>{showContForm ? 'Cancelar' : 'Nuevo'}</Button>}>
            {showContForm && (
              <form onSubmit={onContSubmit} className="space-y-4 mb-4">
                <div className="grid grid-cols-4 gap-4">
                  <Select label="Tipo doc" value={contForm.doc_type} onChange={(e) => setContForm({ ...contForm, doc_type: e.target.value })}>
                    <option value="invoice">Factura</option>
                    <option value="credit_note">Nota de crédito</option>
                  </Select>
                  <Input label="Prefijo" value={contForm.prefix} onChange={(e) => setContForm({ ...contForm, prefix: e.target.value })} />
                  <Input label="Desde" type="number" value={contForm.start_seq} onChange={(e) => setContForm({ ...contForm, start_seq: e.target.value })} />
                  <Input label="Hasta" type="number" value={contForm.end_seq} onChange={(e) => setContForm({ ...contForm, end_seq: e.target.value })} />
                </div>
                <Button type="submit">Guardar</Button>
              </form>
            )}
            {contingencies.length === 0 ? (
              <p className="text-sm text-gray-400 py-4 text-center">Sin libros de contingencia</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead><tr className="border-b border-gray-200 text-left text-gray-500">
                    <th className="py-2 px-3 font-medium">Tipo doc</th>
                    <th className="py-2 px-3 font-medium">Prefijo</th>
                    <th className="py-2 px-3 font-medium">Desde</th>
                    <th className="py-2 px-3 font-medium">Hasta</th>
                    <th className="py-2 px-3 font-medium">Actual</th>
                  </tr></thead>
                  <tbody>
                    {contingencies.map((c) => (
                      <tr key={c.id} className="border-b border-gray-100">
                        <td className="py-2 px-3 text-gray-800">{c.doc_type}</td>
                        <td className="py-2 px-3 text-gray-800">{c.prefix}</td>
                        <td className="py-2 px-3 text-gray-800">{c.start_seq}</td>
                        <td className="py-2 px-3 text-gray-800">{c.end_seq}</td>
                        <td className="py-2 px-3 text-gray-800">{c.current_seq}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </Card>
        </>
      )}
    </div>
  );
}
