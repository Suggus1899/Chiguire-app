import Link from 'next/link';
import {
  WifiOff, FileText, Landmark, CreditCard, Package, Building2,
  Code2, Users, Truck, BarChart3, Globe, Monitor, Smartphone,
  CheckCircle2, ArrowRight, Zap, ShieldCheck, Clock,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { buttonVariants } from '@/components/ui/button';
import { cn } from '@/lib/utils';

/* ─── datos ─── */

const features = [
  { icon: FileText,   title: 'Facturación y Retenciones',   items: ['Forma Libre, Máquina Fiscal e Imprenta Digital', 'Comprobantes IVA e ISLR', 'Notas de crédito y débito', 'Cuentas por Cobrar / Pagar', 'Numeración correlativa SENIAT'] },
  { icon: Package,    title: 'Inventario y Operaciones',     items: ['Control de stock por almacén y sucursal', 'Transferencias y guías de despacho', 'Manufactura y órdenes de producción', 'Ajustes y control de merma', 'Kardex y libro Art. 177'] },
  { icon: CreditCard, title: 'Links de Pago',                items: ['Cashea Link — cuotas para tus clientes', 'Biopago BDV — débito Banco de Venezuela', 'Spidi — bolívares y cripto inmediato', 'WayuPay — validación bancaria automática', 'Comparte por WhatsApp o correo'] },
  { icon: FileText,   title: 'Presupuestos y Cotizaciones',  items: ['Cotizaciones con un clic a factura', 'Listado, historial y búsqueda', 'PDF listo para WhatsApp', 'Control de estado (vigente / vencido)'] },
  { icon: Package,    title: 'Compras y Proveedores',        items: ['Órdenes de compra con aprobación', 'Facturas de compra y notas de crédito', 'Recepción parcial de mercancía', 'Historial por proveedor'] },
  { icon: BarChart3,  title: 'Reportes Fiscales',            items: ['Libro de Ventas y Libro de Compras', 'Reporte de IGTF percibido', 'Exportación completa a Excel', 'Reportes listos para el contador'] },
  { icon: Landmark,   title: 'Moneda y Fiscal',              items: ['Tasa BCV actualizada cada día', 'Cálculo automático IVA, IGTF, ISLR', 'Facturación en USD y VES simultánea', 'Retenciones configurables por cliente'] },
  { icon: Code2,      title: 'Integraciones y API',          items: ['API REST y webhooks por evento', 'Recordatorios de pago por correo', 'Importación / exportación Excel', 'Validación automática pagos móviles'] },
  { icon: Users,      title: 'Clientes y Vendedores',        items: ['Autocompletado por cédula o RIF', 'Comisiones configurables por factura', 'Asignación de vendedores a clientes', 'Reportes de ventas y comisiones'] },
  { icon: WifiOff,    title: 'Offline e Infraestructura',    items: ['SQLite local en web, móvil y escritorio', 'Sincronización al reconectar', 'Sin pérdida ante cortes de luz o red', 'Android, iOS, Windows, Mac y Linux'] },
];

const plans = [
  {
    name: 'Emprendedor', price: '$15', period: '/mes', popular: false,
    items: ['100 facturas de venta mensuales', 'Hasta $15.000 en ventas', 'Facturas de compra ilimitadas', '1 usuario', 'Facturación Forma Libre', 'Multi-moneda con tasa BCV', 'IVA · IGTF · ISLR', 'App offline (móvil)', 'API access'],
  },
  {
    name: 'Pymes', price: '$28', period: '/mes', popular: true,
    items: ['800 facturas de venta mensuales', 'Hasta $30.000 en ventas', 'Facturas de compra ilimitadas', 'Hasta 5 usuarios con roles', 'Máquinas Fiscales + Imprenta Digital', 'Links de Pago (Cashea, Spidi, Biopago, WayuPay)', 'App offline (móvil y web PWA)', 'Reportes en Excel', 'API access'],
  },
  {
    name: 'Pro', price: '$55', period: '/mes', popular: false,
    items: ['Facturas ilimitadas', 'Sin límite de ventas', 'Usuarios ilimitados', 'Hasta 5 sucursales', 'Manufactura y órdenes de producción', 'App offline (web, móvil y escritorio)', 'Transferencias entre sucursales', 'Reportes avanzados + API', 'Soporte prioritario'],
  },
];

const faqs = [
  { q: '¿Cómo funciona el modo offline?', a: 'Tus datos se guardan en SQLite local (en el dispositivo o el navegador). Al reconectar, la sincronización ocurre automáticamente en segundo plano. Sin pérdida de datos, sin interrupciones.' },
  { q: '¿Está homologado por el SENIAT?', a: 'Sí. Chiguire está integrado con imprentas digitales autorizadas (The Factory HKA, Unidigital) y es compatible con máquinas fiscales HKA, PNP y Bixolon para facturación Forma Libre, Máquina Fiscal e Imprenta Digital.' },
  { q: '¿Qué métodos de pago acepta?', a: 'Cashea Link (cuotas), Biopago BDV (débito banco), Spidi (bolívares y cripto), y WayuPay (validación bancaria automática). Todos integrados nativamente — el enlace se genera desde la misma factura.' },
  { q: '¿En qué dispositivos funciona?', a: 'Web (PWA instalable en Chrome/Edge), app nativa para Android e iOS (Flutter), y app de escritorio para Windows, macOS y Linux (Tauri 2). El mismo dato en todos los dispositivos.' },
  { q: '¿Puedo gestionar varias empresas?', a: 'Sí. Cada empresa es un tenant independiente con su propio inventario, facturas y usuarios. Cambias de empresa en un clic desde el menú principal.' },
  { q: '¿Cómo funciona el período de prueba?', a: '7 días completos con todas las funciones activas. Sin tarjeta de crédito. Al vencerse elige el plan que prefieras — si no, la cuenta queda inactiva sin cobrar nada.' },
];

/* ─── subcomponentes ─── */

function DashboardMockup() {
  const bars = [38, 52, 30, 65, 80, 72, 90, 100];
  const labels = ['E', 'F', 'M', 'A', 'M', 'J', 'J', 'A'];
  return (
    <div className="w-full max-w-sm rounded-2xl border border-white/10 bg-white/5 backdrop-blur-sm p-5 shadow-2xl text-white text-xs ring-1 ring-white/10">
      <div className="flex items-center justify-between mb-4">
        <div>
          <p className="text-white/60 text-[10px] mb-0.5">Dashboard · Agosto 2026</p>
          <p className="font-semibold text-sm">Resumen mensual</p>
        </div>
        <span className="flex items-center gap-1 text-emerald-400 text-xs bg-emerald-400/10 px-2 py-1 rounded-full border border-emerald-400/20">
          <span className="w-1.5 h-1.5 bg-emerald-400 rounded-full" />
          +31%
        </span>
      </div>
      <div className="grid grid-cols-3 gap-2 mb-5">
        {[['Facturas', '318'], ['Ventas', '$6.4K'], ['Clientes', '47']].map(([l, v]) => (
          <div key={l} className="rounded-xl bg-white/5 border border-white/5 p-2.5 text-center">
            <p className="text-white/40 mb-1 text-[10px]">{l}</p>
            <p className="font-bold text-sm">{v}</p>
          </div>
        ))}
      </div>
      <p className="text-white/40 mb-2 text-[10px] uppercase tracking-wider">Ventas Mensuales</p>
      <div className="flex items-end gap-1.5 h-16">
        {bars.map((h, i) => (
          <div key={i} className="flex flex-col items-center flex-1 gap-1">
            <div className="w-full rounded-t bg-blue-500/80" style={{ height: `${h}%` }} />
            <span className="text-white/30" style={{ fontSize: 8 }}>{labels[i]}</span>
          </div>
        ))}
      </div>
      <div className="mt-4 pt-3 border-t border-white/5 flex items-center gap-1.5 text-white/30">
        <span className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-pulse" />
        <span style={{ fontSize: 10 }}>Sincronizado · hace 2 min</span>
      </div>
    </div>
  );
}

/* ─── página ─── */

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-background text-foreground">

      {/* NAV */}
      <nav className="fixed top-0 inset-x-0 z-50 border-b border-white/5 bg-slate-950/90 backdrop-blur-md">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <Link href="/" className="font-bold text-white text-base tracking-tight flex items-center gap-2">
            🦔 <span>Chiguire</span>
          </Link>
          <div className="hidden md:flex items-center gap-6 text-sm text-slate-400">
            <a href="#caracteristicas" className="hover:text-white transition-colors">Características</a>
            <a href="#precios" className="hover:text-white transition-colors">Precios</a>
            <a href="#faq" className="hover:text-white transition-colors">FAQ</a>
            <a href="#contacto" className="hover:text-white transition-colors">Contacto</a>
          </div>
          <div className="flex items-center gap-2">
            <Link href="/login" className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'text-slate-300 hover:text-white hover:bg-white/10')}>
              Ingresar
            </Link>
            <Link href="/register" className={cn(buttonVariants({ size: 'sm' }), 'bg-blue-600 hover:bg-blue-500 text-white border-0')}>
              Prueba gratis
            </Link>
          </div>
        </div>
      </nav>

      {/* HERO */}
      <section className="relative bg-slate-950 pt-28 pb-24 px-6 overflow-hidden">
        {/* grid bg */}
        <div className="pointer-events-none absolute inset-0 bg-[linear-gradient(to_right,#ffffff08_1px,transparent_1px),linear-gradient(to_bottom,#ffffff08_1px,transparent_1px)] bg-[size:64px_64px]" />
        {/* glow */}
        <div className="pointer-events-none absolute -top-40 left-1/2 -translate-x-1/2 w-[600px] h-[400px] rounded-full bg-blue-600/20 blur-[120px]" />

        <div className="relative max-w-6xl mx-auto grid lg:grid-cols-2 gap-14 items-center">
          <div>
            <div className="flex items-center gap-2 mb-6">
              <Badge variant="outline" className="border-blue-500/30 text-blue-400 bg-blue-500/10 gap-1.5 px-3 py-1 h-auto text-xs">
                <span className="w-1.5 h-1.5 bg-blue-400 rounded-full animate-pulse" />
                ERP offline-first · Venezuela
              </Badge>
            </div>
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-white leading-[1.1] tracking-tight mb-6">
              Software que trabaja<br />
              <span className="text-blue-400">aunque se vaya la luz</span>
            </h1>
            <p className="text-slate-400 text-lg leading-relaxed mb-3 max-w-lg">
              Factura, controla inventario y cumple con el SENIAT — con o sin internet. Sincronización automática al reconectar.
            </p>
            <p className="text-slate-500 text-sm mb-8 max-w-lg">
              Incluye <span className="text-slate-300">API REST</span> y <span className="text-slate-300">webhooks</span>. IVA, IGTF, ISLR y tasa BCV actualizada diariamente.
            </p>
            <div className="flex flex-wrap gap-3">
              <Link href="/register" className={cn(buttonVariants({ size: 'lg' }), 'bg-blue-600 hover:bg-blue-500 text-white border-0 gap-2')}>
                Prueba 7 días gratis <ArrowRight className="size-4" />
              </Link>
              <a href="#caracteristicas" className={cn(buttonVariants({ variant: 'outline', size: 'lg' }), 'border-white/20 text-slate-300 hover:bg-white/5 hover:text-white hover:border-white/30')}>
                Ver características
              </a>
            </div>
            <p className="text-slate-600 text-xs mt-5">
              ¿Ya eres cliente?{' '}
              <Link href="/login" className="text-blue-400 hover:underline">Inicia sesión aquí</Link>
            </p>
          </div>
          <div className="flex justify-center lg:justify-end">
            <DashboardMockup />
          </div>
        </div>
      </section>

      {/* STATS */}
      <section className="border-y border-border bg-muted/30 py-6 px-6">
        <div className="max-w-5xl mx-auto grid grid-cols-2 md:grid-cols-4 gap-6 text-center">
          {[
            { icon: Zap,          val: '< 500 ms', label: 'Carga inicial offline' },
            { icon: ShieldCheck,  val: '100%',     label: 'Cumplimiento fiscal SENIAT' },
            { icon: Clock,        val: '7 días',   label: 'Prueba gratuita completa' },
            { icon: Globe,        val: '5',         label: 'Plataformas soportadas' },
          ].map(({ icon: Icon, val, label }) => (
            <div key={label} className="flex flex-col items-center gap-1">
              <Icon className="size-4 text-blue-500 mb-1" />
              <p className="font-bold text-xl text-foreground">{val}</p>
              <p className="text-xs text-muted-foreground">{label}</p>
            </div>
          ))}
        </div>
      </section>

      {/* QUÉ OFRECEMOS */}
      <section className="py-20 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-14">
            <Badge variant="secondary" className="mb-4">¿Qué Ofrecemos?</Badge>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">Un ERP completo para Venezuela</h2>
            <p className="text-muted-foreground text-lg max-w-xl mx-auto">Desde emprendedores hasta grandes empresas, con o sin internet.</p>
          </div>
          <div className="grid md:grid-cols-3 gap-6">
            {[
              { icon: WifiOff,    title: 'Offline Primero',       desc: 'SQLite local en web, móvil y escritorio. Tus datos siempre disponibles. Sincronización automática al reconectar, sin pérdida de información.' },
              { icon: Landmark,   title: 'Multi-moneda y Fiscal', desc: 'Tasa BCV diaria automática. IGTF, IVA e ISLR calculados al instante. Facturación homologada Forma Libre, Máquina Fiscal e Imprenta Digital.' },
              { icon: Code2,      title: 'Integración Total',     desc: 'Cashea, Spidi, WayuPay y Biopago nativos. API REST y webhooks para conectar con cualquier sistema contable o herramienta de automatización.' },
            ].map(({ icon: Icon, title, desc }) => (
              <Card key={title} className="border-border/60 hover:border-blue-500/30 hover:shadow-md transition-all">
                <CardHeader>
                  <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center mb-2">
                    <Icon className="size-5 text-blue-500" />
                  </div>
                  <CardTitle className="text-base font-semibold">{title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-sm leading-relaxed">{desc}</CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* PLATAFORMAS */}
      <section className="py-16 px-6 bg-slate-950">
        <div className="max-w-5xl mx-auto text-center">
          <h2 className="text-2xl font-bold text-white mb-2">Disponible en todas tus pantallas</h2>
          <p className="text-slate-400 mb-10 text-sm">Mismo dato, en tiempo real, desde cualquier dispositivo.</p>
          <div className="grid sm:grid-cols-3 gap-4">
            {[
              { icon: Globe,      label: 'Web / PWA',            sub: 'Instálala como app · Chrome, Edge, Firefox' },
              { icon: Smartphone, label: 'Android e iOS',        sub: 'App nativa Flutter con SQLite offline' },
              { icon: Monitor,    label: 'Windows · Mac · Linux', sub: 'App de escritorio Tauri + impresora fiscal' },
            ].map(({ icon: Icon, label, sub }) => (
              <div key={label} className="rounded-2xl border border-white/10 bg-white/5 px-6 py-7 flex flex-col items-center gap-3 hover:border-blue-500/30 hover:bg-blue-500/5 transition-all">
                <div className="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center">
                  <Icon className="size-6 text-blue-400" />
                </div>
                <span className="font-semibold text-white text-sm">{label}</span>
                <span className="text-xs text-slate-500 text-center">{sub}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CARACTERÍSTICAS DETALLADAS */}
      <section id="caracteristicas" className="py-20 px-6 bg-muted/20">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-14">
            <Badge variant="secondary" className="mb-4">Características detalladas</Badge>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">Todo lo que necesitas para operar</h2>
            <p className="text-muted-foreground text-lg max-w-lg mx-auto">En una sola plataforma, pensada para el contexto venezolano.</p>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
            {features.map(({ icon: Icon, title, items }) => (
              <Card key={title} className="border-border/60">
                <CardHeader className="pb-2">
                  <div className="flex items-center gap-2.5">
                    <div className="w-8 h-8 rounded-lg bg-blue-500/10 flex items-center justify-center shrink-0">
                      <Icon className="size-4 text-blue-500" />
                    </div>
                    <CardTitle className="text-sm font-semibold">{title}</CardTitle>
                  </div>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-1.5 mt-1">
                    {items.map((item) => (
                      <li key={item} className="flex items-start gap-2 text-xs text-muted-foreground">
                        <CheckCircle2 className="size-3.5 text-emerald-500 mt-0.5 shrink-0" />
                        {item}
                      </li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* PRECIOS */}
      <section id="precios" className="py-20 px-6">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-14">
            <Badge variant="secondary" className="mb-4">Planes y Precios</Badge>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">Sin contratos, cancela cuando quieras</h2>
            <p className="text-muted-foreground text-lg max-w-lg mx-auto">Elige el plan que mejor se adapte a tu negocio.</p>
          </div>
          <div className="grid md:grid-cols-3 gap-6 items-start">
            {plans.map((p) => (
              p.popular ? (
                /* Popular — dark custom */
                <div key={p.name} className="relative rounded-2xl bg-slate-900 text-white p-7 flex flex-col gap-5 ring-2 ring-blue-500 shadow-xl shadow-blue-500/10">
                  <span className="absolute -top-3.5 left-1/2 -translate-x-1/2 bg-blue-600 text-white text-xs font-semibold px-4 py-1 rounded-full">Más Popular</span>
                  <div>
                    <p className="text-slate-400 text-sm mb-1">{p.name}</p>
                    <p className="text-4xl font-bold">{p.price}<span className="text-base font-normal text-slate-400 ml-1">{p.period}</span></p>
                  </div>
                  <ul className="space-y-2 flex-1">
                    {p.items.map((item) => (
                      <li key={item} className="flex items-start gap-2 text-xs text-slate-300">
                        <CheckCircle2 className="size-3.5 text-emerald-400 mt-0.5 shrink-0" />
                        {item}
                      </li>
                    ))}
                  </ul>
                  <div className="flex flex-col gap-2">
                    <Link href="/register" className={cn(buttonVariants({ size: 'default' }), 'w-full bg-blue-600 hover:bg-blue-500 text-white border-0 justify-center')}>Prueba 7 días gratis</Link>
                    <Link href="/login" className="text-center text-xs text-slate-500 hover:text-slate-300 py-1.5 transition-colors">Ingresa si ya eres cliente</Link>
                  </div>
                </div>
              ) : (
                <Card key={p.name} className="border-border/60 flex flex-col">
                  <CardHeader>
                    <CardDescription>{p.name}</CardDescription>
                    <CardTitle className="text-4xl font-bold text-foreground">
                      {p.price}<span className="text-base font-normal text-muted-foreground ml-1">{p.period}</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="flex-1 flex flex-col gap-5">
                    <ul className="space-y-2 flex-1">
                      {p.items.map((item) => (
                        <li key={item} className="flex items-start gap-2 text-xs text-muted-foreground">
                          <CheckCircle2 className="size-3.5 text-emerald-500 mt-0.5 shrink-0" />
                          {item}
                        </li>
                      ))}
                    </ul>
                    <div className="flex flex-col gap-2">
                      <Link href="/register" className={cn(buttonVariants({ variant: 'outline', size: 'default' }), 'w-full justify-center')}>Prueba 7 días gratis</Link>
                      <Link href="/login" className="text-center text-xs text-muted-foreground hover:text-foreground py-1.5 transition-colors">Ingresa si ya eres cliente</Link>
                    </div>
                  </CardContent>
                </Card>
              )
            ))}
          </div>

          {/* footnotes */}
          <div className="mt-8 rounded-xl border border-border bg-muted/40 p-5 text-xs text-muted-foreground space-y-1.5">
            <p><strong className="text-foreground">** Máquinas fiscales:</strong> The Factory HKA, PNP, Bixolon. Requiere Windows 8.1+.</p>
            <p><strong className="text-foreground">** Imprenta Digital:</strong> Integración con The Factory HKA y Unidigital. La autorización la otorga el SENIAT.</p>
            <p>Todos los precios incluyen IVA según tasa oficial al momento del pago.</p>
          </div>

          <div className="mt-5 rounded-xl border border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-950/30 p-5 text-center">
            <p className="font-medium text-foreground mb-1">¿Necesitas un plan personalizado?</p>
            <p className="text-sm text-muted-foreground mb-3">Más facturas, sucursales adicionales o integraciones específicas — lo diseñamos a tu medida.</p>
            <a href="mailto:hola@chiguire.app" className={cn(buttonVariants({ variant: 'outline', size: 'sm' }))}>
              Solicitar cotización sin compromiso
            </a>
          </div>
        </div>
      </section>

      {/* FAQ */}
      <section id="faq" className="py-20 px-6 bg-muted/20">
        <div className="max-w-3xl mx-auto">
          <div className="text-center mb-12">
            <Badge variant="secondary" className="mb-4">Preguntas Frecuentes</Badge>
            <h2 className="text-3xl font-bold tracking-tight mb-3">Todo lo que necesitas saber</h2>
          </div>
          <div className="divide-y divide-border rounded-xl border border-border overflow-hidden bg-card">
            {faqs.map((f) => (
              <details key={f.q} className="group px-6 py-1">
                <summary className="flex cursor-pointer items-center justify-between py-4 text-sm font-medium text-foreground hover:text-blue-500 transition-colors list-none [&::-webkit-details-marker]:hidden">
                  {f.q}
                  <span className="ml-4 shrink-0 text-muted-foreground transition-transform group-open:rotate-45 text-lg leading-none">+</span>
                </summary>
                <div className="pb-4 text-sm text-muted-foreground leading-relaxed">{f.a}</div>
              </details>
            ))}
          </div>
        </div>
      </section>

      {/* CTA FINAL */}
      <section className="py-24 px-6 bg-slate-950 relative overflow-hidden">
        <div className="pointer-events-none absolute inset-0 bg-[linear-gradient(to_right,#ffffff06_1px,transparent_1px),linear-gradient(to_bottom,#ffffff06_1px,transparent_1px)] bg-[size:64px_64px]" />
        <div className="pointer-events-none absolute bottom-0 left-1/2 -translate-x-1/2 w-[400px] h-[200px] rounded-full bg-blue-600/15 blur-[80px]" />
        <div className="relative max-w-2xl mx-auto text-center">
          <h2 className="text-4xl md:text-5xl font-bold text-white tracking-tight mb-4">¿Listo para comenzar?</h2>
          <p className="text-slate-400 text-lg mb-8">7 días con todas las funciones, sin tarjeta de crédito. Sin compromisos.</p>
          <Link href="/register" className={cn(buttonVariants({ size: 'lg' }), 'bg-blue-600 hover:bg-blue-500 text-white border-0 gap-2')}>
            Prueba 7 días gratis <ArrowRight className="size-4" />
          </Link>
        </div>
      </section>

      {/* CONTACTO */}
      <section id="contacto" className="py-20 px-6">
        <div className="max-w-3xl mx-auto text-center">
          <Badge variant="secondary" className="mb-6">Contáctanos</Badge>
          <h2 className="text-3xl font-bold tracking-tight mb-2">Estamos para ayudarte</h2>
          <p className="text-muted-foreground mb-10">Lunes a viernes · 8:00 AM – 6:00 PM</p>
          <div className="grid sm:grid-cols-2 gap-5">
            <Card className="border-border/60 text-left hover:border-blue-500/30 transition-colors">
              <CardHeader><CardTitle className="text-sm flex items-center gap-2"><span>📧</span> Email</CardTitle></CardHeader>
              <CardContent>
                <a href="mailto:hola@chiguire.app" className="text-blue-500 hover:underline text-sm">hola@chiguire.app</a>
                <p className="text-xs text-muted-foreground mt-2">Ventas, soporte y facturación</p>
              </CardContent>
            </Card>
            <Card className="border-border/60 text-left hover:border-blue-500/30 transition-colors">
              <CardHeader><CardTitle className="text-sm flex items-center gap-2"><span>💬</span> WhatsApp</CardTitle></CardHeader>
              <CardContent>
                <a href="https://wa.me/584140000000" className="text-blue-500 hover:underline text-sm">+58 414-000-0000</a>
                <p className="text-xs text-muted-foreground mt-2">Respuesta en menos de 1 hora en horario hábil</p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      <Separator />

      {/* FOOTER */}
      <footer className="py-8 px-6 bg-background">
        <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-muted-foreground">
          <span>🦔 Chiguire ERP · Hecho para Venezuela · © 2026</span>
          <div className="flex gap-6">
            <a href="#caracteristicas" className="hover:text-foreground transition-colors">Características</a>
            <a href="#precios" className="hover:text-foreground transition-colors">Precios</a>
            <a href="#faq" className="hover:text-foreground transition-colors">FAQ</a>
            <Link href="/login" className="hover:text-foreground transition-colors">Ingresar</Link>
            <Link href="/register" className="hover:text-foreground transition-colors">Registro</Link>
          </div>
        </div>
      </footer>

    </div>
  );
}
