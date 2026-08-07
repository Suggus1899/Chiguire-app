'use client';
import { useState, useEffect, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';

type Tenant = { id: string; name: string; slug: string; plan: string };
type TestItem = { id: string; label: string; created_at: string };

export default function Dashboard() {
  const router = useRouter();
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [selected, setSelected] = useState<Tenant | null>(null);
  const [items, setItems] = useState<TestItem[]>([]);
  const [newTenantName, setNewTenantName] = useState('');
  const [newItemLabel, setNewItemLabel] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.tenants.list()
      .then(setTenants)
      .catch(() => router.push('/login'))
      .finally(() => setLoading(false));
  }, [router]);

  useEffect(() => {
    if (!selected) return;
    api.testItems.list().then(setItems).catch(console.error);
  }, [selected]);

  async function createTenant(e: FormEvent) {
    e.preventDefault();
    if (!newTenantName) return;
    const slug = newTenantName.toLowerCase().replace(/\s+/g, '-');
    const t = await api.tenants.create(newTenantName, slug);
    setTenants((prev) => [...prev, { ...t, plan: 'trial' }]);
    setNewTenantName('');
  }

  async function addItem(e: FormEvent) {
    e.preventDefault();
    if (!newItemLabel || !selected) return;
    const item = await api.testItems.create(newItemLabel);
    setItems((prev) => [item, ...prev]);
    setNewItemLabel('');
  }

  function logout() {
    api.auth.logout();
    router.push('/login');
  }

  if (loading) return <div className="p-8 text-gray-500">Cargando...</div>;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center">
        <h1 className="text-lg font-semibold text-gray-900">Chiguire</h1>
        <button onClick={logout} className="text-sm text-gray-500 hover:text-gray-900">
          Salir
        </button>
      </header>

      <main className="max-w-4xl mx-auto p-6 space-y-6">
        {/* Tenant selector */}
        <section className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="font-medium text-gray-900 mb-4">Empresas</h2>
          <div className="flex flex-wrap gap-2 mb-4">
            {tenants.map((t) => (
              <button
                key={t.id}
                onClick={() => setSelected(t)}
                className={`px-3 py-1.5 rounded-lg text-sm font-medium border transition-colors ${
                  selected?.id === t.id
                    ? 'bg-blue-600 text-white border-blue-600'
                    : 'bg-white text-gray-700 border-gray-300 hover:border-blue-400'
                }`}
              >
                {t.name}
              </button>
            ))}
          </div>
          <form onSubmit={createTenant} className="flex gap-2">
            <input
              type="text"
              placeholder="Nueva empresa..."
              value={newTenantName}
              onChange={(e) => setNewTenantName(e.target.value)}
              className="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              type="submit"
              className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700"
            >
              Crear
            </button>
          </form>
        </section>

        {/* Phase 0 sync test */}
        {selected && (
          <section className="bg-white rounded-xl border border-gray-200 p-6">
            <h2 className="font-medium text-gray-900 mb-1">
              Test de sync — <span className="text-blue-600">{selected.name}</span>
            </h2>
            <p className="text-xs text-gray-400 mb-4">
              Estos items se sincronizan via PowerSync. Solo son visibles para este tenant.
            </p>

            <form onSubmit={addItem} className="flex gap-2 mb-4">
              <input
                type="text"
                placeholder="Etiqueta del item..."
                value={newItemLabel}
                onChange={(e) => setNewItemLabel(e.target.value)}
                className="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="submit"
                className="bg-green-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-green-700"
              >
                Agregar
              </button>
            </form>

            {items.length === 0 ? (
              <p className="text-sm text-gray-400">Sin items todavía.</p>
            ) : (
              <ul className="space-y-2">
                {items.map((item) => (
                  <li
                    key={item.id}
                    className="flex justify-between text-sm py-2 border-b border-gray-100"
                  >
                    <span className="text-gray-800">{item.label}</span>
                    <span className="text-gray-400">{new Date(item.created_at).toLocaleString('es-VE')}</span>
                  </li>
                ))}
              </ul>
            )}
          </section>
        )}
      </main>
    </div>
  );
}
