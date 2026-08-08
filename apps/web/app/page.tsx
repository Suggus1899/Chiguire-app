import Link from 'next/link';

/* ───── data ───── */

const features = [
  {
    title: 'Facturación y Retenciones',
    items: [
      'Forma Libre, Máquina Fiscal e Imprenta Digital',
      'Comprobantes de retención IVA e ISLR',
      'Notas de crédito y débito',
      'Cuentas por Cobrar / por Pagar',
      'Numeración correlativa automática SENIAT',
    ],
  },
  {
    title: 'Inventario y Operaciones',
    items: [
      'Control de stock por almacén y sucursal',
      'Transferencias y guías de despacho',
      'Manufactura y órdenes de producción',
      'Ajustes de inventario y control de merma',
      'Kardex y libro de inventario (Art. 177)',
    ],
  },
  {
    title: 'Links de Pago',
    items: [
      'Cashea Link — financiamiento en cuotas',
      'Biopago BDV — débito Banco de Venezuela',
      'Spidi — bolívares y cripto, confirmación inmediata',
      'WayuPay — bolívares o dólares con validación bancaria',
      'Comparte por WhatsApp, correo o redes',
    ],
  },
  {
    title: 'Presupuestos y Cotizaciones',
    items: [
      'Crea presupuestos y cotizaciones para clientes',
      'Conviértelos en factura en un clic',
      'Listado, historial y búsqueda',
      'PDF listo para enviar por WhatsApp',
    ],
  },
  {
    title: 'Compras y Proveedores',
    items: [
      'Órdenes de compra con aprobación',
      'Registro de facturas de compra',
      'Notas de crédito y débito de proveedores',
      'Recepción parcial de mercancía',
    ],
  },
  {
    title: 'Reportes Fiscales',
    items: [
      'Libro de Ventas y Libro de Compras',
      'Reporte de IGTF percibido',
      'Exportación completa a Excel',
      'Reportes listos para el contador',
    ],
  },
  {
    title: 'Moneda y Fiscal',
    items: [
      'Tasa BCV actualizada automáticamente cada día',
      'Cálculo automático de IVA, IGTF, ISLR',
      'Facturación en USD y VES simultáneamente',
      'Retenciones configurables por cliente',
    ],
  },
  {
    title: 'Integraciones y API',
    items: [
      'API REST y webhooks por evento',
      'Recordatorios de pago por correo',
      'Importación y exportación Excel',
      'Validación automática de pagos móviles',
    ],
  },
  {
    title: 'Clientes y Vendedores',
    items: [
      'Autocompletado por cédula o RIF venezolano',
      'Perfiles de vendedores con comisiones por factura',
      'Asignación de vendedores a clientes',
      'Reportes de ventas y comisiones',
    ],
  },
  {
    title: 'Offline e Infraestructura',
    items: [
      'SQLite local en web, móvil y escritorio',
      'Sincronización automática al reconectar',
      'Sin pérdida de datos ante cortes de luz o internet',
      'App nativa para Android, iOS, Windows, Mac y Linux',
    ],
  },
];

const plans = [
  {
    name: 'Emprendedor',
    price: '$15',
    period: '/mes',
    highlight: false,
    badge: null,
    items: [
      '100 facturas de venta mensuales',
      'Hasta $15.000 en ventas mensuales',
      'Facturas de compra ilimitadas',
      '1 usuario (dueño de la cuenta)',
      'Facturación Forma Libre',
      'Productos ilimitados',
      'Multi-moneda con tasa del día',
      'IVA, IGTF, retenciones ISLR',
      'Cuentas por Pagar/Cobrar',
      'Picking y rutas de entrega',
      'App offline (móvil)',
      'Acceso a la API',
    ],
    cta: 'Pruébalo 7 días gratis',
    sub: 'Ingresa si ya eres cliente',
  },
  {
    name: 'Pymes',
    price: '$28',
    period: '/mes',
    highlight: true,
    badge: 'Más Popular',
    items: [
      '800 facturas de venta mensuales',
      'Hasta $30.000 en ventas mensuales',
      'Facturas de compra ilimitadas',
      'Hasta 5 usuarios con roles y permisos',
      'Máquinas Fiscales**',
      'Facturación Digital**',
      'Facturación Forma Libre',
      'Productos ilimitados',
      'Multi-moneda con tasa del día',
      'Links de Pago (Cashea, Spidi, Biopago, WayuPay)',
      'App offline (móvil y web PWA)',
      'IVA, IGTF, ISLR, notas de crédito/débito',
      'Reportes en Excel',
      'Acceso a la API',
    ],
    cta: 'Pruébalo 7 días gratis',
    sub: 'Ingresa si ya eres cliente',
  },
  {
    name: 'Pro',
    price: '$55',
    period: '/mes',
    highlight: false,
    badge: null,
    items: [
      'Facturas ilimitadas',
      'Sin límite en el monto de ventas',
      'Usuarios ilimitados con roles y permisos',
      'Máquinas Fiscales**',
      'Facturación Digital**',
      'Manufactura y órdenes de producción',
      'Hasta 5 sucursales',
      'App offline (web, móvil y escritorio)',
      'Links de Pago completos',
      'Transferencias entre sucursales',
      'Reportes en Excel + API avanzada',
      'Soporte prioritario',
    ],
    cta: 'Pruébalo 7 días gratis',
    sub: 'Ingresa si ya eres cliente',
  },
];

/* ───── components ───── */

function DashboardMockup() {
  const bars = [
    { h: 38, label: 'Ene', val: '$2.1K' },
    { h: 52, label: 'Feb', val: '$3.4K' },
    { h: 30, label: 'Mar', val: '$1.9K' },
    { h: 65, label: 'Abr', val: '$4.2K' },
    { h: 80, label: 'May', val: '$5.1K' },
    { h: 72, label: 'Jun', val: '$4.6K' },
    { h: 90, label: 'Jul', val: '$5.8K' },
    { h: 100, label: 'Ago', val: '$6.4K' },
  ];

  return (
    <div className="rounded-2xl border border-slate-700 bg-slate-800/80 backdrop-blur p-5 shadow-2xl w-full max-w-md text-white text-xs">
      {/* header */}
      <div className="flex items-center justify-between mb-4">
        <span className="font-semibold text-sm text-white">Dashboard — Agosto 2026</span>
        <span className="text-green-400 font-semibold text-xs bg-green-400/10 px-2 py-0.5 rounded-full">↗ +31%</span>
      </div>
      {/* stat tiles */}
      <div className="grid grid-cols-3 gap-2 mb-5">
        {[
          { label: 'Facturas', val: '318' },
          { label: 'Ventas', val: '$6.4K' },
          { label: 'Clientes', val: '47' },
        ].map((s) => (
          <div key={s.label} className="bg-slate-700/60 rounded-xl p-2.5 text-center">
            <p className="text-slate-400 mb-1">{s.label}</p>
            <p className="font-bold text-white text-sm">{s.val}</p>
          </div>
        ))}
      </div>
      {/* bar chart */}
      <p className="text-slate-400 mb-2">Ventas Mensuales</p>
      <div className="flex items-end gap-1.5 h-20">
        {bars.map((b) => (
          <div key={b.label} className="flex flex-col items-center flex-1 gap-1">
            <div
              className="w-full rounded-t bg-blue-500 opacity-80 hover:opacity-100 transition-opacity"
              style={{ height: `${b.h}%` }}
              title={b.val}
            />
            <span className="text-slate-500" style={{ fontSize: 9 }}>{b.label}</span>
          </div>
        ))}
      </div>
      {/* offline badge */}
      <div className="mt-4 flex items-center gap-1.5 text-slate-400">
        <span className="w-1.5 h-1.5 bg-green-400 rounded-full" />
        <span style={{ fontSize: 10 }}>Sincronizado · último hace 2 min</span>
      </div>
    </div>
  );
}

/* ───── page ───── */

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white text-gray-900">

      {/* ── Nav ── */}
      <nav className="fixed top-0 inset-x-0 z-50 bg-slate-900/95 backdrop-blur border-b border-slate-800">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <span className="font-bold text-white text-base tracking-tight">🦔 Chiguire</span>
          <div className="hidden md:flex items-center gap-6 text-sm text-slate-300">
            <a href="#caracteristicas" className="hover:text-white transition-colors">Características</a>
            <a href="#precios" className="hover:text-white transition-colors">Precios</a>
            <a href="#contacto" className="hover:text-white transition-colors">Contacto</a>
          </div>
          <div className="flex items-center gap-2">
            <Link href="/login" className="text-sm text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition-colors">
              Ingresar
            </Link>
            <Link href="/register" className="text-sm bg-blue-600 text-white px-4 py-1.5 rounded-lg hover:bg-blue-500 transition-colors font-medium">
              Prueba gratis
            </Link>
          </div>
        </div>
      </nav>

      {/* ── Hero ── */}
      <section className="bg-slate-900 pt-28 pb-20 px-6">
        <div className="max-w-6xl mx-auto grid lg:grid-cols-2 gap-12 items-center">
          {/* left */}
          <div>
            <div className="inline-flex items-center gap-2 bg-blue-500/10 text-blue-400 text-xs font-medium px-3 py-1.5 rounded-full mb-6 border border-blue-500/20">
              <span className="w-1.5 h-1.5 bg-blue-400 rounded-full animate-pulse" />
              ERP offline-first para Venezuela
            </div>
            <h1 className="text-4xl md:text-5xl font-bold text-white leading-tight mb-5">
              Software Administrativo<br />
              <span className="text-blue-400">que funciona sin internet</span>
            </h1>
            <p className="text-slate-400 text-lg leading-relaxed mb-4">
              Facturación, inventario y cumplimiento SENIAT — con o sin conexión. Tus datos se sincronizan automáticamente cuando vuelve la red.
            </p>
            <p className="text-slate-500 text-sm mb-8">
              Incluye <span className="text-slate-300">API REST</span> y <span className="text-slate-300">webhooks</span> para integrarse con sistemas contables y automatizaciones. IVA, IGTF, ISLR y tasa BCV actualizada a diario.
            </p>
            <div className="flex flex-col sm:flex-row gap-3">
              <Link href="/register" className="bg-blue-600 text-white px-6 py-3 rounded-xl text-sm font-semibold hover:bg-blue-500 transition-colors text-center">
                Prueba 7 días gratis
              </Link>
              <a href="#caracteristicas" className="border border-slate-700 text-slate-300 px-6 py-3 rounded-xl text-sm font-medium hover:border-slate-600 hover:text-white transition-colors text-center">
                Ver características
              </a>
            </div>
            <p className="text-slate-500 text-xs mt-4">
              ¿Ya eres cliente?{' '}
              <Link href="/login" className="text-blue-400 hover:underline">Inicia sesión aquí</Link>
            </p>
          </div>
          {/* right — dashboard mockup */}
          <div className="flex justify-center lg:justify-end">
            <DashboardMockup />
          </div>
        </div>
      </section>

      {/* ── Ofrecemos ── */}
      <section className="py-20 px-6 bg-gray-50">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 mb-3">¿Qué Ofrecemos?</h2>
          <p className="text-center text-gray-500 mb-12 text-lg">Un sistema ERP completo que se adapta a tu negocio, desde emprendedores hasta grandes empresas.</p>
          <div className="grid md:grid-cols-3 gap-6">
            {[
              {
                icon: '📡',
                title: 'Offline Primero',
                desc: 'Trabaja sin internet. SQLite local en web, móvil y escritorio. Sincronización automática en segundo plano al reconectar.',
              },
              {
                icon: '🏛️',
                title: 'Multi-moneda y Fiscal',
                desc: 'Tasa BCV diaria automática. IGTF, retenciones IVA e ISLR. Facturación homologada Forma Libre, Máquina Fiscal e Imprenta Digital.',
              },
              {
                icon: '🔗',
                title: 'Integración Total',
                desc: 'Cashea, Spidi, WayuPay y Biopago nativos. API REST y webhooks para conectar con cualquier sistema contable o automatización.',
              },
            ].map((c) => (
              <div key={c.title} className="bg-white rounded-2xl p-7 border border-gray-100 shadow-sm">
                <div className="text-4xl mb-4">{c.icon}</div>
                <h3 className="font-semibold text-gray-900 text-lg mb-2">{c.title}</h3>
                <p className="text-gray-500 text-sm leading-relaxed">{c.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Plataformas ── */}
      <section className="py-16 px-6 bg-slate-900">
        <div className="max-w-5xl mx-auto text-center">
          <h2 className="text-2xl font-bold text-white mb-3">Disponible en todas tus pantallas</h2>
          <p className="text-slate-400 mb-10">Mismo dato, en tiempo real, desde cualquier dispositivo.</p>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            {[
              { icon: '🌐', label: 'Web / PWA', sub: 'Instálala como app — Chrome, Edge, Firefox' },
              { icon: '📱', label: 'Android e iOS', sub: 'App nativa Flutter con SQLite offline' },
              { icon: '💻', label: 'Windows · Mac · Linux', sub: 'App de escritorio Tauri con impresora fiscal' },
            ].map((p) => (
              <div key={p.label} className="bg-slate-800 rounded-2xl px-6 py-7 border border-slate-700 flex flex-col items-center gap-3">
                <span className="text-4xl">{p.icon}</span>
                <span className="font-semibold text-white text-sm">{p.label}</span>
                <span className="text-xs text-slate-400 text-center">{p.sub}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Características detalladas ── */}
      <section id="caracteristicas" className="py-20 px-6">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 mb-3">Características detalladas</h2>
          <p className="text-center text-gray-500 mb-14 text-lg">Todo lo que necesitas para operar tu negocio en Venezuela, en una sola plataforma.</p>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((f) => (
              <div key={f.title} className="bg-gray-50 rounded-2xl p-6 border border-gray-100">
                <h3 className="font-semibold text-gray-900 mb-3 text-sm">{f.title}</h3>
                <ul className="space-y-1.5">
                  {f.items.map((item) => (
                    <li key={item} className="flex items-start gap-2 text-sm text-gray-600">
                      <span className="text-green-500 mt-0.5 shrink-0">✓</span>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Precios ── */}
      <section id="precios" className="py-20 px-6 bg-gray-50">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 mb-3">Planes y Precios</h2>
          <p className="text-center text-gray-500 mb-14 text-lg">Elige el plan que mejor se adapte a tu negocio. Sin contratos, cancela cuando quieras.</p>
          <div className="grid md:grid-cols-3 gap-6 items-start">
            {plans.map((p) => (
              <div
                key={p.name}
                className={`rounded-2xl p-7 flex flex-col gap-5 relative ${
                  p.highlight
                    ? 'bg-slate-900 text-white ring-2 ring-blue-500'
                    : 'bg-white border border-gray-100'
                }`}
              >
                {p.badge && (
                  <span className="absolute -top-3 left-1/2 -translate-x-1/2 bg-blue-600 text-white text-xs font-semibold px-3 py-1 rounded-full">
                    {p.badge}
                  </span>
                )}
                <div>
                  <p className={`text-sm font-medium mb-1 ${p.highlight ? 'text-slate-400' : 'text-gray-500'}`}>{p.name}</p>
                  <p className={`text-4xl font-bold ${p.highlight ? 'text-white' : 'text-gray-900'}`}>
                    {p.price}
                    <span className={`text-base font-normal ml-1 ${p.highlight ? 'text-slate-400' : 'text-gray-400'}`}>{p.period}</span>
                  </p>
                </div>
                <ul className="space-y-2 flex-1">
                  {p.items.map((item) => (
                    <li key={item} className={`text-xs flex items-start gap-2 ${p.highlight ? 'text-slate-300' : 'text-gray-600'}`}>
                      <span className={`mt-0.5 shrink-0 ${p.highlight ? 'text-green-400' : 'text-green-500'}`}>✓</span>
                      {item}
                    </li>
                  ))}
                </ul>
                <div className="flex flex-col gap-2">
                  <Link
                    href="/register"
                    className={`text-center text-sm font-semibold py-2.5 rounded-xl transition-colors ${
                      p.highlight
                        ? 'bg-blue-600 text-white hover:bg-blue-500'
                        : 'bg-slate-900 text-white hover:bg-slate-700'
                    }`}
                  >
                    {p.cta}
                  </Link>
                  <Link
                    href="/login"
                    className={`text-center text-xs py-1.5 rounded-xl transition-colors ${
                      p.highlight ? 'text-slate-400 hover:text-slate-200' : 'text-gray-400 hover:text-gray-600'
                    }`}
                  >
                    {p.sub}
                  </Link>
                </div>
              </div>
            ))}
          </div>

          {/* footnote */}
          <div className="mt-10 bg-white rounded-2xl border border-gray-100 p-6 text-sm text-gray-500 space-y-2">
            <p>** <strong className="text-gray-700">Máquinas fiscales compatibles:</strong> The Factory HKA (SRP812, DT230, HKA80), PNP y Bixolon. Requiere Windows 8.1+ para el software de conexión.</p>
            <p>** <strong className="text-gray-700">Facturación Digital:</strong> Integración con imprentas digitales autorizadas (The Factory HKA, Unidigital). La autorización la otorga el SENIAT — Chiguire facilita la integración.</p>
            <p>Todos los precios incluyen IVA y se calculan según la tasa oficial vigente al momento del pago.</p>
          </div>

          {/* custom plan */}
          <div className="mt-6 bg-blue-50 border border-blue-100 rounded-2xl p-6 text-center">
            <p className="text-gray-700 font-medium mb-1">¿Necesitas un plan personalizado?</p>
            <p className="text-gray-500 text-sm mb-4">Si tu negocio requiere más facturas, sucursales adicionales o funcionalidades específicas, podemos diseñar un plan a tu medida.</p>
            <a href="mailto:hola@chiguire.app" className="inline-block bg-white border border-gray-200 text-gray-700 text-sm font-medium px-5 py-2.5 rounded-xl hover:border-gray-300 hover:shadow-sm transition-all">
              Solicitar cotización sin compromiso
            </a>
          </div>
        </div>
      </section>

      {/* ── CTA final ── */}
      <section className="py-24 px-6 bg-slate-900 text-center">
        <div className="max-w-2xl mx-auto">
          <h2 className="text-4xl font-bold text-white mb-4">¿Listo para comenzar?</h2>
          <p className="text-slate-400 text-lg mb-8">7 días con todas las funciones, sin tarjeta de crédito. Sin compromisos.</p>
          <Link href="/register" className="inline-block bg-blue-600 text-white px-8 py-4 rounded-xl text-sm font-semibold hover:bg-blue-500 transition-colors">
            Prueba 7 días gratis →
          </Link>
        </div>
      </section>

      {/* ── Contacto ── */}
      <section id="contacto" className="py-20 px-6 bg-gray-50">
        <div className="max-w-3xl mx-auto text-center">
          <h2 className="text-3xl font-bold text-gray-900 mb-3">Contáctanos</h2>
          <p className="text-gray-500 mb-10">Atendemos de lunes a viernes, 8:00 AM – 6:00 PM.</p>
          <div className="grid sm:grid-cols-2 gap-6">
            <div className="bg-white rounded-2xl p-7 border border-gray-100 text-left">
              <p className="text-sm font-semibold text-gray-900 mb-1">📧 Email</p>
              <a href="mailto:hola@chiguire.app" className="text-blue-600 hover:underline text-sm">hola@chiguire.app</a>
              <p className="text-gray-500 text-xs mt-3">Ventas, soporte y facturación</p>
            </div>
            <div className="bg-white rounded-2xl p-7 border border-gray-100 text-left">
              <p className="text-sm font-semibold text-gray-900 mb-1">💬 WhatsApp</p>
              <a href="https://wa.me/584140000000" className="text-blue-600 hover:underline text-sm">+58 414-000-0000</a>
              <p className="text-gray-500 text-xs mt-3">Respuesta en menos de 1 hora en horario hábil</p>
            </div>
          </div>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="border-t border-gray-100 py-8 px-6 bg-white">
        <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-gray-400">
          <span>🦔 Chiguire ERP — Hecho para Venezuela · © 2026</span>
          <div className="flex gap-6">
            <a href="#caracteristicas" className="hover:text-gray-600 transition-colors">Características</a>
            <a href="#precios" className="hover:text-gray-600 transition-colors">Precios</a>
            <Link href="/login" className="hover:text-gray-600 transition-colors">Ingresar</Link>
            <Link href="/register" className="hover:text-gray-600 transition-colors">Registro</Link>
          </div>
        </div>
      </footer>

    </div>
  );
}
