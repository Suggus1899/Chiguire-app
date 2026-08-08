'use client';
import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      await api.auth.login(email, password);
      router.push('/dashboard');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Credenciales incorrectas');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex">
      {/* ── Panel izquierdo (branding) ── */}
      <div className="hidden lg:flex lg:w-1/2 bg-slate-900 flex-col justify-between p-12">
        <Link href="/" className="text-white font-bold text-xl tracking-tight">🦔 Chiguire</Link>

        <div>
          <h2 className="text-3xl font-bold text-white leading-snug mb-4">
            Tu ERP funciona<br />
            <span className="text-blue-400">aunque se vaya la luz</span>
          </h2>
          <p className="text-slate-400 text-sm leading-relaxed mb-10">
            Offline-first desde el día uno. Tus facturas, inventario y fiscal siempre disponibles — con o sin internet.
          </p>
          <ul className="space-y-4">
            {[
              { icon: '📡', text: 'Datos en SQLite local · sincroniza al reconectar' },
              { icon: '🏛️', text: 'IVA, IGTF, ISLR y tasa BCV automática' },
              { icon: '💳', text: 'Cashea, Spidi, WayuPay y Biopago integrados' },
            ].map((f) => (
              <li key={f.text} className="flex items-start gap-3 text-sm text-slate-300">
                <span className="text-lg leading-none mt-0.5">{f.icon}</span>
                {f.text}
              </li>
            ))}
          </ul>
        </div>

        <p className="text-slate-600 text-xs">© 2026 Chiguire ERP · Hecho para Venezuela</p>
      </div>

      {/* ── Panel derecho (form) ── */}
      <div className="flex-1 flex items-center justify-center bg-gray-50 px-6 py-12">
        <div className="w-full max-w-sm">
          {/* mobile logo */}
          <Link href="/" className="lg:hidden block text-slate-900 font-bold text-xl mb-8">🦔 Chiguire</Link>

          <h1 className="text-2xl font-bold text-gray-900 mb-1">Bienvenido de vuelta</h1>
          <p className="text-sm text-gray-500 mb-8">Ingresa a tu cuenta para continuar</p>

          {error && (
            <div className="mb-5 p-3 rounded-xl bg-red-50 border border-red-100 text-sm text-red-600">
              {error}
            </div>
          )}

          <form onSubmit={onSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
              <input
                type="email"
                required
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="tu@empresa.com"
                className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder:text-gray-300 transition-shadow"
              />
            </div>

            <div>
              <div className="flex items-center justify-between mb-1.5">
                <label className="text-sm font-medium text-gray-700">Contraseña</label>
              </div>
              <input
                type="password"
                required
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder:text-gray-300 transition-shadow"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-slate-900 text-white py-2.5 rounded-xl text-sm font-semibold hover:bg-slate-700 disabled:opacity-50 transition-colors mt-2"
            >
              {loading ? 'Ingresando...' : 'Ingresar'}
            </button>
          </form>

          <p className="text-center text-sm text-gray-500 mt-6">
            ¿No tienes cuenta?{' '}
            <Link href="/register" className="text-blue-600 font-medium hover:underline">
              Prueba 7 días gratis
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
