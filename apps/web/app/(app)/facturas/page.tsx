'use client';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import { api, type Invoice } from '@/lib/api';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Table } from '@/components/Table';

export default function FacturasPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.invoices.list()
      .then(setInvoices)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-semibold text-gray-900">Facturas</h1>
        <Link href="/facturas/nueva"><Button>Nueva factura</Button></Link>
      </div>

      <Card>
        {loading ? (
          <p className="text-gray-500">Cargando...</p>
        ) : (
          <Table<Invoice>
            columns={[
              {
                key: 'number',
                label: 'Número',
                render: (r) => <Link href={`/facturas/${r.id}`} className="text-blue-600 hover:underline">{r.number || r.id.slice(0, 8)}</Link>,
              },
              { key: 'status', label: 'Estado' },
              { key: 'currency', label: 'Moneda' },
              { key: 'total', label: 'Total', render: (r) => (r.total ?? 0).toFixed(2) },
              { key: 'issue_date', label: 'Fecha', render: (r) => r.issue_date ? new Date(r.issue_date).toLocaleDateString('es-VE') : '—' },
            ]}
            data={invoices}
            empty="Sin facturas registradas"
          />
        )}
      </Card>
    </div>
  );
}
