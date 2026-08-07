'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navItems = [
  { href: '/dashboard', label: 'Dashboard' },
  { href: '/clientes', label: 'Clientes' },
  { href: '/proveedores', label: 'Proveedores' },
  { href: '/productos', label: 'Productos' },
  { href: '/inventario', label: 'Inventario' },
  { href: '/facturas', label: 'Facturas' },
  { href: '/cotizaciones', label: 'Cotizaciones' },
  { href: '/compras', label: 'Compras' },
  { href: '/fiscal', label: 'Fiscal' },
  { href: '/pagos', label: 'Pagos' },
  { href: '/rutas', label: 'Rutas' },
  { href: '/vendedores', label: 'Vendedores' },
  { href: '/metodos-pago', label: 'Métodos de Pago' },
  { href: '/dispositivos-fiscales', label: 'Dispositivos Fiscales' },
  { href: '/cuentas', label: 'Cuentas por Cobrar/Pagar' },
  { href: '/informes', label: 'Informes' },
  { href: '/notas-credito', label: 'Notas de Crédito' },
  { href: '/transferencias', label: 'Transferencias' },
  { href: '/manufactura', label: 'Manufactura' },
  { href: '/picking', label: 'Picking' },
  { href: '/api-tokens', label: 'API Tokens' },
  { href: '/onboarding', label: 'Onboarding' },
  { href: '/configuracion', label: 'Configuración' },
];

export function Sidebar() {
  const pathname = usePathname();
  return (
    <aside className="w-56 bg-white border-r border-gray-200 min-h-screen flex flex-col">
      <div className="px-5 py-4 border-b border-gray-200">
        <Link href="/dashboard" className="text-lg font-semibold text-gray-900">Chiguire</Link>
      </div>
      <nav className="flex-1 py-2 space-y-0.5">
        {navItems.map((item) => {
          const active = pathname === item.href || pathname.startsWith(item.href + '/');
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`block px-5 py-2 text-sm font-medium transition-colors ${
                active ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-600' : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              {item.label}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
