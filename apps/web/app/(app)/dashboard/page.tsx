'use client';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import { api, type Invoice, type Product } from '@/lib/api';
import { Card } from '@/components/Card';

export default function DashboardPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.invoices.list().catch(() => [] as Invoice[]),
      api.products.list().catch(() => [] as Product[]),
    ]).then(([inv, prod]) => {
      setInvoices(inv);
      setProducts(prod);
      setLoading(false);
    });
  }, []);

  if (loading) return <div className="text-gray-500">Cargando...</div>;

  const recentInvoices = invoices.slice(0, 5);
  const totalInvoiced = invoices.reduce((s, i) => s + (i.total || 0), 0);

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card title="Facturas">
          <p className="text-3xl font-bold text-gray-900">{invoices.length}</p>
          <p className="text-sm text-gray-500 mt-1">Total facturado: {totalInvoiced.toFixed(2)}</p>
        </Card>
        <Card title="Productos">
          <p className="text-3xl font-bold text-gray-900">{products.length}</p>
          <p className="text-sm text-gray-500 mt-1">Productos registrados</p>
        </Card>
        <Card title="Alertas de stock">
          <p className="text-3xl font-bold text-gray-900">0</p>
          <p className="text-sm text-gray-500 mt-1">Productos con stock bajo</p>
        </Card>
      </div>

      <Card title="Facturas recientes" actions={<Link href="/facturas" className="text-sm text-blue-600 hover:underline">Ver todas</Link>}>
        {recentInvoices.length === 0 ? (
          <p className="text-sm text-gray-400">Sin facturas todavía.</p>
        ) : (
          <ul className="space-y-2">
            {recentInvoices.map((inv) => (
              <li key={inv.id} className="flex justify-between text-sm py-2 border-b border-gray-100">
                <Link href={`/facturas/${inv.id}`} className="text-blue-600 hover:underline">
                  {inv.number || inv.id.slice(0, 8)}
                </Link>
                <span className="text-gray-500">{inv.status} — {inv.total?.toFixed(2)} {inv.currency}</span>
              </li>
            ))}
          </ul>
        )}
      </Card>
    </div>
  );
}
