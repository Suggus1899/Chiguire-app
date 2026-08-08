'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navGroups = [
  {
    label: 'Ventas',
    items: [
      { href: '/dashboard', label: 'Dashboard', icon: '▦' },
      { href: '/facturas', label: 'Facturas', icon: '🧾' },
      { href: '/cotizaciones', label: 'Cotizaciones', icon: '📋' },
      { href: '/notas-credito', label: 'Notas de Crédito', icon: '↩' },
      { href: '/clientes', label: 'Clientes', icon: '👥' },
    ],
  },
  {
    label: 'Inventario',
    items: [
      { href: '/productos', label: 'Productos', icon: '🏷' },
      { href: '/inventario', label: 'Inventario', icon: '📦' },
      { href: '/transferencias', label: 'Transferencias', icon: '↔' },
      { href: '/manufactura', label: 'Manufactura', icon: '⚙' },
      { href: '/picking', label: 'Picking', icon: '✔' },
    ],
  },
  {
    label: 'Compras',
    items: [
      { href: '/compras', label: 'Órdenes de Compra', icon: '🛒' },
      { href: '/proveedores', label: 'Proveedores', icon: '🏭' },
    ],
  },
  {
    label: 'Fiscal y Pagos',
    items: [
      { href: '/fiscal', label: 'Fiscal', icon: '🏛' },
      { href: '/pagos', label: 'Pagos', icon: '💳' },
      { href: '/cuentas', label: 'Cuentas por C/P', icon: '⚖' },
      { href: '/metodos-pago', label: 'Métodos de Pago', icon: '🔗' },
      { href: '/dispositivos-fiscales', label: 'Dispositivos Fiscales', icon: '🖨' },
    ],
  },
  {
    label: 'Vendedores',
    items: [
      { href: '/vendedores', label: 'Vendedores', icon: '👤' },
      { href: '/vendedores/comisiones', label: 'Comisiones', icon: '💰' },
      { href: '/rutas', label: 'Rutas', icon: '🗺' },
    ],
  },
  {
    label: 'Sistema',
    items: [
      { href: '/informes', label: 'Informes', icon: '📊' },
      { href: '/api-tokens', label: 'API Tokens', icon: '🔑' },
      { href: '/configuracion', label: 'Configuración', icon: '⚙' },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-56 bg-slate-900 min-h-screen flex flex-col shrink-0">
      {/* logo */}
      <div className="px-5 py-4 border-b border-slate-800">
        <Link href="/dashboard" className="text-white font-bold text-base tracking-tight">
          🦔 Chiguire
        </Link>
      </div>

      {/* nav */}
      <nav className="flex-1 overflow-y-auto py-3">
        {navGroups.map((group) => (
          <div key={group.label} className="mb-1">
            <p className="px-5 pt-3 pb-1 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              {group.label}
            </p>
            {group.items.map((item) => {
              const active = pathname === item.href || (item.href !== '/dashboard' && pathname.startsWith(item.href + '/'));
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`flex items-center gap-2.5 px-5 py-2 text-sm transition-colors ${
                    active
                      ? 'bg-blue-600/20 text-blue-400 border-r-2 border-blue-500 font-medium'
                      : 'text-slate-400 hover:text-white hover:bg-slate-800'
                  }`}
                >
                  <span className="text-base leading-none opacity-70">{item.icon}</span>
                  <span>{item.label}</span>
                </Link>
              );
            })}
          </div>
        ))}
      </nav>

      {/* footer */}
      <div className="px-5 py-4 border-t border-slate-800">
        <div className="flex items-center gap-2 text-xs text-slate-500">
          <span className="w-1.5 h-1.5 bg-green-400 rounded-full" />
          Sincronizado
        </div>
      </div>
    </aside>
  );
}
