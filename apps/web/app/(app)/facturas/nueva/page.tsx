'use client';
import { useState, useEffect, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { api, type Customer, type Product } from '@/lib/api';
import { toast } from '@/lib/toast';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Select } from '@/components/Select';

type LineItem = { product_id: string; description: string; quantity: number; unit_price: number; tax_rate: number };

export default function NuevaFacturaPage() {
  const router = useRouter();
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [customerId, setCustomerId] = useState('');
  const [currency, setCurrency] = useState('VES');
  const [items, setItems] = useState<LineItem[]>([]);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    Promise.all([
      api.customers.list().catch((err) => { console.error('operation failed:', err); return [] as Customer[]; }),
      api.products.list().catch((err) => { console.error('operation failed:', err); return [] as Product[]; }),
    ]).then(([cust, prod]) => {
      setCustomers(cust);
      setProducts(prod);
      setLoading(false);
    });
  }, []);

  function addItem() {
    setItems((prev) => [...prev, { product_id: '', description: '', quantity: 1, unit_price: 0, tax_rate: 16 }]);
  }

  function updateItem(index: number, patch: Partial<LineItem>) {
    setItems((prev) => prev.map((it, i) => (i === index ? { ...it, ...patch } : it)));
  }

  function removeItem(index: number) {
    setItems((prev) => prev.filter((_, i) => i !== index));
  }

  function selectProduct(index: number, productId: string) {
    const product = products.find((p) => p.id === productId);
    updateItem(index, {
      product_id: productId,
      description: product?.name ?? '',
      unit_price: product?.cost ?? 0,
    });
  }

  const subtotal = items.reduce((s, it) => s + it.quantity * it.unit_price, 0);
  const tax = items.reduce((s, it) => s + it.quantity * it.unit_price * (it.tax_rate / 100), 0);
  const total = subtotal + tax;

  async function save(emit: boolean) {
    if (!customerId || items.length === 0) return;
    setSaving(true);
    try {
      const inv = await api.invoices.create({
        customer_id: customerId,
        currency,
        subtotal,
        tax,
        total,
        status: 'draft',
      });
      await api.invoices.addItems(inv.id, items);
      if (emit) {
        await api.invoices.emit(inv.id);
      }
      router.push('/facturas');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Error al guardar');
    } finally {
      setSaving(false);
    }
  }

  if (loading) return <p className="text-gray-500">Cargando...</p>;

  return (
    <div className="space-y-6">
      <h1 className="text-xl font-semibold text-gray-900">Nueva factura</h1>

      <Card title="Datos generales">
        <div className="grid grid-cols-2 gap-4">
          <Select label="Cliente" required value={customerId} onChange={(e) => setCustomerId(e.target.value)}>
            <option value="">Seleccionar cliente...</option>
            {customers.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
          </Select>
          <Select label="Moneda" value={currency} onChange={(e) => setCurrency(e.target.value)}>
            <option value="VES">Bolívares (VES)</option>
            <option value="USD">Dólares (USD)</option>
          </Select>
        </div>
      </Card>

      <Card title="Items" actions={<Button variant="secondary" onClick={addItem}>Agregar item</Button>}>
        {items.length === 0 ? (
          <p className="text-sm text-gray-400 py-2">Sin items. Agrega al menos uno.</p>
        ) : (
          <div className="space-y-3">
            {items.map((item, i) => (
              <div key={i} className="grid grid-cols-12 gap-2 items-end">
                <div className="col-span-4">
                  <Select label={i === 0 ? 'Producto' : undefined} value={item.product_id} onChange={(e) => selectProduct(i, e.target.value)}>
                    <option value="">Seleccionar...</option>
                    {products.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                  </Select>
                </div>
                <div className="col-span-3">
                  <Input label={i === 0 ? 'Descripción' : undefined} value={item.description} onChange={(e) => updateItem(i, { description: e.target.value })} />
                </div>
                <div className="col-span-1">
                  <Input label={i === 0 ? 'Cant.' : undefined} type="number" value={item.quantity} onChange={(e) => updateItem(i, { quantity: parseFloat(e.target.value) || 0 })} />
                </div>
                <div className="col-span-2">
                  <Input label={i === 0 ? 'Precio' : undefined} type="number" step="0.01" value={item.unit_price} onChange={(e) => updateItem(i, { unit_price: parseFloat(e.target.value) || 0 })} />
                </div>
                <div className="col-span-1">
                  <Input label={i === 0 ? 'IVA%' : undefined} type="number" value={item.tax_rate} onChange={(e) => updateItem(i, { tax_rate: parseFloat(e.target.value) || 0 })} />
                </div>
                <div className="col-span-1">
                  <Button variant="danger" onClick={() => removeItem(i)} className="!px-2">×</Button>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="mt-4 pt-4 border-t border-gray-200 space-y-1 text-sm">
          <div className="flex justify-between"><span className="text-gray-500">Subtotal:</span><span>{subtotal.toFixed(2)}</span></div>
          <div className="flex justify-between"><span className="text-gray-500">IVA:</span><span>{tax.toFixed(2)}</span></div>
          <div className="flex justify-between font-semibold text-base"><span>Total:</span><span>{total.toFixed(2)} {currency}</span></div>
        </div>
      </Card>

      <div className="flex gap-3">
        <Button variant="secondary" onClick={() => save(false)} disabled={saving || !customerId || items.length === 0}>
          {saving ? 'Guardando...' : 'Guardar borrador'}
        </Button>
        <Button onClick={() => save(true)} disabled={saving || !customerId || items.length === 0}>
          {saving ? 'Emitiendo...' : 'Guardar y emitir'}
        </Button>
      </div>
    </div>
  );
}
