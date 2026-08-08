'use client';
import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { WifiOff, Landmark, CreditCard } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { api } from '@/lib/api';

const bullets = [
  { icon: WifiOff,    text: 'Datos en SQLite local · sincroniza al reconectar' },
  { icon: Landmark,   text: 'IVA, IGTF, ISLR y tasa BCV automática' },
  { icon: CreditCard, text: 'Cashea, Spidi, WayuPay y Biopago integrados' },
];

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
      {/* panel izquierdo */}
      <div className="hidden lg:flex lg:w-[52%] bg-slate-950 flex-col justify-between p-12 relative overflow-hidden">
        <div className="pointer-events-none absolute inset-0 bg-[linear-gradient(to_right,#ffffff06_1px,transparent_1px),linear-gradient(to_bottom,#ffffff06_1px,transparent_1px)] bg-[size:48px_48px]" />
        <div className="pointer-events-none absolute -bottom-20 -left-20 w-80 h-80 rounded-full bg-blue-600/15 blur-[80px]" />
        <div className="relative">
          <Link href="/" className="text-white font-bold text-lg tracking-tight flex items-center gap-2">
            🦔 Chiguire
          </Link>
        </div>
        <div className="relative">
          <p className="text-xs font-semibold text-blue-400 uppercase tracking-widest mb-3">ERP offline-first</p>
          <h2 className="text-3xl font-bold text-white leading-snug mb-3">
            Tu ERP funciona<br />
            <span className="text-blue-400">aunque se vaya la luz</span>
          </h2>
          <p className="text-slate-400 text-sm leading-relaxed mb-8 max-w-sm">
            Offline-first desde el día uno. Tus facturas, inventario y fiscal — siempre disponibles.
          </p>
          <ul className="space-y-4">
            {bullets.map(({ icon: Icon, text }) => (
              <li key={text} className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center shrink-0">
                  <Icon className="size-4 text-blue-400" />
                </div>
                <span className="text-sm text-slate-300">{text}</span>
              </li>
            ))}
          </ul>
        </div>
        <p className="relative text-slate-700 text-xs">© 2026 Chiguire ERP · Hecho para Venezuela</p>
      </div>

      {/* panel derecho */}
      <div className="flex-1 flex items-center justify-center bg-muted/20 px-6 py-12">
        <div className="w-full max-w-sm">
          <Link href="/" className="lg:hidden block font-bold text-xl mb-8">🦔 Chiguire</Link>

          <div className="mb-8">
            <h1 className="text-2xl font-bold tracking-tight text-foreground mb-1">Bienvenido de vuelta</h1>
            <p className="text-sm text-muted-foreground">Ingresa a tu cuenta para continuar</p>
          </div>

          {error && (
            <div className="mb-5 p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-sm text-destructive">
              {error}
            </div>
          )}

          <form onSubmit={onSubmit} className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                required
                autoComplete="email"
                placeholder="tu@empresa.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="h-10"
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="password">Contraseña</Label>
              <Input
                id="password"
                type="password"
                required
                autoComplete="current-password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="h-10"
              />
            </div>
            <Button type="submit" disabled={loading} className="w-full h-10 bg-blue-600 hover:bg-blue-500 text-white border-0 mt-2">
              {loading ? 'Ingresando...' : 'Ingresar'}
            </Button>
          </form>

          <p className="text-center text-sm text-muted-foreground mt-6">
            ¿No tienes cuenta?{' '}
            <Link href="/register" className="text-blue-500 font-medium hover:underline">Prueba 7 días gratis</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
