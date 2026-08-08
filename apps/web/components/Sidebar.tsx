'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  LayoutDashboard, Receipt, ClipboardList, CornerDownLeft,
  Users, Tag, Package, ArrowLeftRight, Cog, ScanLine,
  ShoppingCart, Factory, Landmark, CreditCard, Scale,
  Wallet, Printer, UserCheck, BadgeDollarSign, Route,
  BarChart2, KeyRound, Settings2, type LucideIcon,
} from 'lucide-react';
import { cn } from '@/lib/utils';

type NavItem = { href: string; label: string; icon: LucideIcon };
type NavGroup = { label: string; items: NavItem[] };

const groups: NavGroup[] = [
  {
    label: 'Ventas',
    items: [
      { href: '/dashboard',       label: 'Dashboard',        icon: LayoutDashboard },
      { href: '/facturas',        label: 'Facturas',         icon: Receipt },
      { href: '/cotizaciones',    label: 'Cotizaciones',     icon: ClipboardList },
      { href: '/notas-credito',   label: 'Notas de Crédito', icon: CornerDownLeft },
      { href: '/clientes',        label: 'Clientes',         icon: Users },
    ],
  },
  {
    label: 'Inventario',
    items: [
      { href: '/productos',       label: 'Productos',        icon: Tag },
      { href: '/inventario',      label: 'Inventario',       icon: Package },
      { href: '/transferencias',  label: 'Transferencias',   icon: ArrowLeftRight },
      { href: '/manufactura',     label: 'Manufactura',      icon: Cog },
      { href: '/picking',         label: 'Picking',          icon: ScanLine },
    ],
  },
  {
    label: 'Compras',
    items: [
      { href: '/compras',         label: 'Órdenes de Compra', icon: ShoppingCart },
      { href: '/proveedores',     label: 'Proveedores',       icon: Factory },
    ],
  },
  {
    label: 'Fiscal y Pagos',
    items: [
      { href: '/fiscal',              label: 'Fiscal',                icon: Landmark },
      { href: '/pagos',               label: 'Pagos',                 icon: CreditCard },
      { href: '/cuentas',             label: 'Cuentas por C/P',       icon: Scale },
      { href: '/metodos-pago',        label: 'Métodos de Pago',       icon: Wallet },
      { href: '/dispositivos-fiscales', label: 'Dispositivos Fiscales', icon: Printer },
    ],
  },
  {
    label: 'Vendedores',
    items: [
      { href: '/vendedores',            label: 'Vendedores',  icon: UserCheck },
      { href: '/vendedores/comisiones', label: 'Comisiones',  icon: BadgeDollarSign },
      { href: '/rutas',                 label: 'Rutas',       icon: Route },
    ],
  },
  {
    label: 'Sistema',
    items: [
      { href: '/informes',      label: 'Informes',       icon: BarChart2 },
      { href: '/api-tokens',    label: 'API Tokens',     icon: KeyRound },
      { href: '/configuracion', label: 'Configuración',  icon: Settings2 },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();

  function isActive(href: string) {
    if (href === '/dashboard') return pathname === href;
    return pathname === href || pathname.startsWith(href + '/');
  }

  return (
    <aside className="w-56 bg-slate-950 min-h-screen flex flex-col shrink-0 border-r border-white/[0.06]">
      {/* logo */}
      <div className="px-4 py-4 border-b border-white/[0.06]">
        <Link href="/dashboard" className="flex items-center gap-2 text-white font-bold text-sm tracking-tight hover:opacity-80 transition-opacity">
          🦔 <span>Chiguire</span>
        </Link>
      </div>

      {/* nav */}
      <nav className="flex-1 overflow-y-auto py-2 scrollbar-thin">
        {groups.map((group) => (
          <div key={group.label} className="mb-1">
            <p className="px-4 pt-4 pb-1.5 text-[10px] font-semibold text-slate-600 uppercase tracking-widest">
              {group.label}
            </p>
            {group.items.map(({ href, label, icon: Icon }) => {
              const active = isActive(href);
              return (
                <Link
                  key={href}
                  href={href}
                  className={cn(
                    'flex items-center gap-2.5 px-4 py-2 text-[13px] transition-colors relative',
                    active
                      ? 'text-white bg-white/[0.07] before:absolute before:left-0 before:top-1 before:bottom-1 before:w-0.5 before:bg-blue-500 before:rounded-full'
                      : 'text-slate-500 hover:text-slate-200 hover:bg-white/[0.04]'
                  )}
                >
                  <Icon className={cn('size-3.5 shrink-0', active ? 'text-blue-400' : 'text-slate-600')} />
                  <span>{label}</span>
                </Link>
              );
            })}
          </div>
        ))}
      </nav>

      {/* sync status */}
      <div className="px-4 py-3 border-t border-white/[0.06]">
        <div className="flex items-center gap-2 text-[11px] text-slate-600">
          <span className="w-1.5 h-1.5 bg-emerald-500 rounded-full animate-pulse shrink-0" />
          Sincronizado
        </div>
      </div>
    </aside>
  );
}
