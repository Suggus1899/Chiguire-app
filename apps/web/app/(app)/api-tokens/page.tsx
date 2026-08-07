'use client';
import { useState, useEffect, FormEvent } from 'react';
import { api, type ApiToken } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Table } from '@/components/Table';
import { Modal } from '@/components/Modal';

export default function ApiTokensPage() {
  const [tokens, setTokens] = useState<ApiToken[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: '', expires_at: '' });
  const [newToken, setNewToken] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    api.apiTokens.list()
      .then(setTokens)
      .catch((err) => console.error('operation failed:', err))
      .finally(() => setLoading(false));
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const res = await api.apiTokens.create({ name: form.name, expires_at: form.expires_at || undefined });
    setNewToken(res.token);
    setForm({ name: '', expires_at: '' });
    setShowForm(false);
    api.apiTokens.list().then(setTokens).catch((err) => console.error('operation failed:', err));
  }

  async function revoke(id: string) {
    await api.apiTokens.revoke(id);
    setTokens((prev) => prev.filter((t) => t.id !== id));
  }

  function copyToken() {
    if (!newToken) return;
    navigator.clipboard.writeText(newToken);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">API Tokens</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancelar' : 'Nuevo token'}</Button>
      </div>

      {showForm && (
        <Card title="Nuevo token">
          <form onSubmit={onSubmit} className="space-y-4">
            <Input label="Nombre" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
            <Input label="Expira (opcional)" type="date" value={form.expires_at} onChange={(e) => setForm({ ...form, expires_at: e.target.value })} />
            <Button type="submit">Crear</Button>
          </form>
        </Card>
      )}

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<ApiToken>
            columns={[
              { key: 'name', label: 'Nombre' },
              { key: 'token_preview', label: 'Token' },
              { key: 'last_used_at', label: 'Último uso', render: (r) => r.last_used_at?.slice(0, 10) ?? 'Nunca' },
              { key: 'expires_at', label: 'Expira', render: (r) => r.expires_at.slice(0, 10) },
              { key: 'status', label: 'Estado', render: (r) => (r.is_active ? 'Activo' : 'Revocado') },
              { key: 'actions', label: '', render: (row) => row.is_active ? <button onClick={() => revoke(row.id)} className="text-red-600 text-xs hover:underline">Revocar</button> : null },
            ]}
            data={tokens}
            empty="Sin tokens registrados"
          />
        )}
      </Card>

      <Modal open={!!newToken} onClose={() => setNewToken(null)} title="Token creado">
        <div className="space-y-4">
          <p className="text-sm text-red-600 font-medium">Copia el token ahora. No se mostrará nuevamente.</p>
          <div className="bg-gray-50 border border-gray-200 rounded-lg p-3 break-all text-sm font-mono">{newToken}</div>
          <Button onClick={copyToken}>{copied ? 'Copiado' : 'Copiar token'}</Button>
        </div>
      </Modal>
    </div>
  );
}
