'use client';
import { useRouter } from 'next/navigation';
import { Sidebar } from '@/components/Sidebar';
import { api } from '@/lib/api';

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  function logout() {
    api.auth.logout();
    router.push('/login');
  }

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <header className="bg-white border-b border-gray-200 px-6 py-3 flex justify-between items-center">
          <h1 className="text-sm font-medium text-gray-500">ERP Chiguire</h1>
          <button onClick={logout} className="text-sm text-gray-500 hover:text-gray-900">
            Salir
          </button>
        </header>
        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}
