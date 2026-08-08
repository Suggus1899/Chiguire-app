'use client';
import { useRouter, usePathname } from 'next/navigation';
import { Sidebar } from '@/components/Sidebar';
import { api } from '@/lib/api';

function breadcrumb(pathname: string): string {
  const map: Record<string, string> = {
    '/dashboard': 'Dashboard',
    '/clientes': 'Clientes',
    '/proveedores': 'Proveedores',
    '/productos': 'Productos',
    '/productos/categorias': 'Categorías',
    '/inventario': 'Inventario',
    '/inventario/almacenes': 'Almacenes',
    '/facturas': 'Facturas',
    '/cotizaciones': 'Cotizaciones',
    '/compras': 'Compras',
    '/fiscal': 'Fiscal',
    '/pagos': 'Pagos',
    '/rutas': 'Rutas',
    '/vendedores': 'Vendedores',
    '/vendedores/comisiones': 'Comisiones',
    '/metodos-pago': 'Métodos de Pago',
    '/dispositivos-fiscales': 'Dispositivos Fiscales',
    '/cuentas': 'Cuentas por Cobrar/Pagar',
    '/informes': 'Informes',
    '/notas-credito': 'Notas de Crédito',
    '/transferencias': 'Transferencias',
    '/manufactura': 'Manufactura',
    '/picking': 'Picking',
    '/api-tokens': 'API Tokens',
    '/configuracion': 'Configuración',
    '/onboarding': 'Onboarding',
  };
  return map[pathname] ?? map[pathname.split('/').slice(0, -1).join('/')] ?? 'Chiguire';
}

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();

  function logout() {
    api.auth.logout();
    router.push('/login');
  }

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar />
      <div className="flex-1 flex flex-col min-w-0">
        {/* header */}
        <header className="sticky top-0 z-10 bg-white border-b border-gray-100 px-6 py-3 flex items-center justify-between gap-4">
          <div>
            <h1 className="text-sm font-semibold text-gray-900">{breadcrumb(pathname)}</h1>
          </div>
          <button
            onClick={logout}
            className="text-xs text-gray-400 hover:text-gray-700 px-3 py-1.5 rounded-lg hover:bg-gray-100 transition-colors font-medium"
          >
            Salir
          </button>
        </header>
        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}
