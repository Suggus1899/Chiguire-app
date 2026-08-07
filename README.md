<div align="center">

# 🦔 Chiguire

### ERP offline-first multi-tenant para empresas venezolanas

Plataforma **ERP** diseñada para el contexto venezolano: facturación, inventario, fiscal SENIAT, pagos locales (Cashea, Spidi, WayuPay, Biopago) y sincronización en la nube. **Funciona sin internet**; sincroniza cuando hay conexión.

</div>

<br>

<div align="center">

## 🛠️ Tech Stack

</div>

<table align="center">
<tr>
<th colspan="5" align="center" width="600"><sub><b>Frontend</b></sub></th>
</tr>
<tr>
<td align="center" width="120">
<a href="https://nextjs.org/" target="_blank"><img src="https://cdn.simpleicons.org/nextdotjs/000000" width="48" height="48" alt="Next.js" /></a>
<br><sub><b><a href="https://nextjs.org/" target="_blank">Next.js 16</a></b></sub>
<br><sub>App Router · PWA</sub>
</td>
<td align="center" width="120">
<a href="https://react.dev/" target="_blank"><img src="https://cdn.simpleicons.org/react/61DAFB" width="48" height="48" alt="React" /></a>
<br><sub><b><a href="https://react.dev/" target="_blank">React 19</a></b></sub>
<br><sub>Server Components</sub>
</td>
<td align="center" width="120">
<a href="https://www.typescriptlang.org/" target="_blank"><img src="https://cdn.simpleicons.org/typescript/3178C6" width="48" height="48" alt="TypeScript" /></a>
<br><sub><b><a href="https://www.typescriptlang.org/" target="_blank">TypeScript 5</a></b></sub>
<br><sub>Type-safe</sub>
</td>
<td align="center" width="120">
<a href="https://tailwindcss.com/" target="_blank"><img src="https://cdn.simpleicons.org/tailwindcss/06B6D4" width="48" height="48" alt="Tailwind CSS" /></a>
<br><sub><b><a href="https://tailwindcss.com/" target="_blank">Tailwind v4</a></b></sub>
<br><sub>Utility-first</sub>
</td>
<td align="center" width="120">
<a href="https://flutter.dev/" target="_blank"><img src="https://cdn.simpleicons.org/flutter/02569B" width="48" height="48" alt="Flutter" /></a>
<br><sub><b><a href="https://flutter.dev/" target="_blank">Flutter</a></b></sub>
<br><sub>Android · iOS</sub>
</td>
</tr>
<tr>
<th colspan="5" align="center" width="600"><sub><b>Backend</b></sub></th>
</tr>
<tr>
<td align="center" width="120">
<a href="https://go.dev/" target="_blank"><img src="https://cdn.simpleicons.org/go/00ADD8" width="48" height="48" alt="Go" /></a>
<br><sub><b><a href="https://go.dev/" target="_blank">Go 1.26</a></b></sub>
<br><sub>Backend API</sub>
</td>
<td align="center" width="120">
<a href="https://github.com/go-chi/chi" target="_blank"><img src="https://raw.githubusercontent.com/go-chi/docs/master/assets/chi.png" width="48" height="48" alt="go-chi" /></a>
<br><sub><b><a href="https://github.com/go-chi/chi" target="_blank">go-chi v5</a></b></sub>
<br><sub>HTTP router</sub>
</td>
<td align="center" width="120">
<a href="https://www.postgresql.org/" target="_blank"><img src="https://cdn.simpleicons.org/postgresql/4169E1" width="48" height="48" alt="PostgreSQL" /></a>
<br><sub><b><a href="https://www.postgresql.org/" target="_blank">PostgreSQL 16</a></b></sub>
<br><sub>Primary DB · RLS</sub>
</td>
<td align="center" width="120">
<a href="https://jwt.io/" target="_blank"><img src="https://cdn.simpleicons.org/jsonwebtokens/000000" width="48" height="48" alt="JWT" /></a>
<br><sub><b><a href="https://jwt.io/" target="_blank">JWT</a></b></sub>
<br><sub>HS256 · 15min access</sub>
</td>
<td align="center" width="120">
<a href="https://github.com/pressly/goose" target="_blank"><img src="https://cdn.simpleicons.org/goose/00ADD8" width="48" height="48" alt="goose" /></a>
<br><sub><b><a href="https://github.com/pressly/goose" target="_blank">goose</a></b></sub>
<br><sub>Migrations</sub>
</td>
</tr>
<tr>
<th colspan="5" align="center" width="600"><sub><b>Sync · Desktop · Tooling</b></sub></th>
</tr>
<tr>
<td align="center" width="120">
<a href="https://www.powersync.com/" target="_blank"><img src="https://cdn.simpleicons.org/powersync/FF6B6B" width="48" height="48" alt="PowerSync" /></a>
<br><sub><b><a href="https://www.powersync.com/" target="_blank">PowerSync</a></b></sub>
<br><sub>Offline-first sync</sub>
</td>
<td align="center" width="120">
<a href="https://tauri.app/" target="_blank"><img src="https://cdn.simpleicons.org/tauri/FFC131" width="48" height="48" alt="Tauri" /></a>
<br><sub><b><a href="https://tauri.app/" target="_blank">Tauri 2</a></b></sub>
<br><sub>Desktop · Rust</sub>
</td>
<td align="center" width="120">
<a href="https://riverpod.dev/" target="_blank"><img src="https://cdn.simpleicons.org/riverpod/00ADD8" width="48" height="48" alt="Riverpod" /></a>
<br><sub><b><a href="https://riverpod.dev/" target="_blank">Riverpod 3</a></b></sub>
<br><sub>State (Flutter)</sub>
</td>
<td align="center" width="120">
<a href="https://pnpm.io/" target="_blank"><img src="https://cdn.simpleicons.org/pnpm/F69220" width="48" height="48" alt="pnpm" /></a>
<br><sub><b><a href="https://pnpm.io/" target="_blank">pnpm 11</a></b></sub>
<br><sub>Workspaces</sub>
</td>
<td align="center" width="120">
<a href="https://www.docker.com/" target="_blank"><img src="https://cdn.simpleicons.org/docker/2496ED" width="48" height="48" alt="Docker" /></a>
<br><sub><b><a href="https://www.docker.com/" target="_blank">Docker</a></b></sub>
<br><sub>PowerSync only</sub>
</td>
</tr>
</table>

<br>

## 📐 Arquitectura

```
┌──────────────────────────────────────────────────────────────┐
│  Clientes                                                     │
│                                                               │
│  Flutter (móvil)    Next.js PWA (web)    Tauri 2 (desktop)   │
│       │                   │                      │            │
│  SQLite (Drift)     SQLite (WASM)      SQLite (nativo)       │
│       └───────────────────┴────────────────────┘             │
│                           │ PowerSync SDK                     │
└───────────────────────────┼──────────────────────────────────┘
                            │ sync bidireccional, incremental
┌───────────────────────────┼──────────────────────────────────┐
│  PowerSync Service         │                                  │
│  (replica Postgres → SQLite, filtra por tenant)               │
└───────────────────────────┼──────────────────────────────────┘
                            │
┌───────────────────────────┼──────────────────────────────────┐
│  PostgreSQL 16             │                                  │
│  Row-Level Security por tenant_id                             │
└───────────────────────────┼──────────────────────────────────┘
                            │
┌───────────────────────────┼──────────────────────────────────┐
│  API Go (chi)              │                                  │
│  Auth JWT · Emisión fiscal · Pagos · Numeración SENIAT        │
└───────────────────────────────────────────────────────────────┘
```

### Reglas offline-first

| Tipo de operación | Flujo |
|---|---|
| Leer datos (clientes, productos, facturas) | SQLite local. Nunca espera red. |
| Crear/editar borrador | SQLite local → PowerSync sincroniza a Postgres. |
| Emitir factura (número correlativo SENIAT) | `POST /invoices/emit` → API asigna número. Se encola offline; emite al reconectar. |
| Cobrar con pasarela (Cashea, etc.) | `POST /payments` → API. Se marca `pending_payment` offline. |
| Stock | Solo movimientos (`stock_movements`, append-only). Nunca se modifica la cantidad absoluta. |

<br>

## 📱 Aplicaciones

| App | Descripción | Puerto | Plataforma |
|-----|-------------|--------|------------|
| **API** | Go + chi + pgx/v5 — auth, facturación, fiscal, pagos, compras, SaaS | `3001` | Backend |
| **Web** | Next.js 16 PWA — 33 rutas, dashboard, facturas, inventario, fiscal, informes | `3000` | Web · PWA instalable |
| **Mobile** | Flutter + Riverpod — 30+ pantallas, offline-first, escáner de códigos | — | Android · iOS |
| **Desktop** | Tauri 2 + Rust — impresora fiscal (The Factory HKA, PNP, Bixolon) | — | Windows · macOS · Linux |

<br>

## 📦 Estructura del Monorepo

```
chiguire/
├── apps/
│   ├── api/                          ← Go API (chi + pgx/v5 + goose)
│   │   ├── cmd/server/               ← Entry point + router
│   │   └── internal/
│   │       ├── auth/                 ← JWT, register, login, refresh
│   │       ├── middleware/           ← RequireAuth, RequireTenant, CORS
│   │       ├── tenant/               ← CRUD de empresas
│   │       ├── customer/             ← Clientes
│   │       ├── vendor/               ← Proveedores
│   │       ├── product/              ← Productos + catálogo
│   │       ├── inventory/            ← Stock + movimientos
│   │       ├── invoice/              ← Facturación forma libre
│   │       ├── fiscal/               ← IGTF, IVA, retenciones, BCV cron
│   │       ├── fiscaldevice/         ← Dispositivos fiscales + talonarios
│   │       ├── payment/              ← Cashea, Spidi, WayuPay, Biopago
│   │       ├── paymentmethod/        ← Métodos de pago configurables
│   │       ├── purchase/             ← Órdenes de compra
│   │       ├── quotation/            ← Cotizaciones
│   │       ├── commission/           ← Comisiones de vendedores
│   │       ├── seller/               ← Vendedores + límites por producto
│   │       ├── delivery/             ← Rutas de reparto
│   │       ├── creditnote/           ← Notas de crédito/débito
│   │       ├── transfer/             ← Transferencias de inventario
│   │       ├── manufacturing/        ← Órdenes de manufactura + BOM
│   │       ├── picking/              ← Listas de picking
│   │       ├── accountspayable/      ← Cuentas por cobrar/pagar
│   │       ├── report/               ← Libros ventas/compras, kardex, art177
│   │       ├── apitoken/             ← API tokens con hash SHA256
│   │       ├── saas/                 ← Suscripciones SaaS
│   │       ├── webhook/              ← Webhooks salientes + dispatcher
│   │       ├── importexport/         ← Import/export Excel
│   │       └── powersync/            ← Token endpoint para sync
│   │
│   ├── web/                          ← Next.js 16 PWA
│   │   ├── app/
│   │   │   ├── (auth)/               ← login, register
│   │   │   └── (app)/                ← 33 rutas: dashboard, facturas,
│   │   │       clientes, productos, inventario, fiscal, pagos,
│   │   │       compras, cotizaciones, vendedores, comisiones,
│   │   │       métodos-pago, dispositivos-fiscales, notas-credito,
│   │   │       transferencias, manufactura, picking, cuentas,
│   │   │       informes, api-tokens, onboarding, configuración
│   │   ├── components/               ← Button, Input, Select, Table, Card, Modal, Sidebar
│   │   └── lib/
│   │       ├── api.ts                ← Cliente HTTP tipado
│   │       └── powersync.ts          ← Schema SQLite + instancia DB local
│   │
│   ├── mobile/                       ← Flutter + Riverpod
│   │   └── lib/
│   │       ├── core/                 ← api.dart, database.dart
│   │       ├── providers/            ← auth, customer, invoice, product
│   │       ├── widgets/              ← app_scaffold
│   │       └── features/             ← 30+ pantallas por feature
│   │           ├── auth/             ← login, register
│   │           ├── dashboard/        ← resumen
│   │           ├── customers/        ← lista + formulario
│   │           ├── products/         ← lista + formulario
│   │           ├── inventory/        ← stock + movimientos
│   │           ├── invoices/         ← lista + formulario + detalle
│   │           ├── quotations/       ← cotizaciones
│   │           ├── purchases/        ← órdenes de compra
│   │           ├── fiscal/           ← libros + BCV
│   │           ├── payments/         ← links de pago
│   │           ├── delivery/         ← rutas de reparto
│   │           ├── sellers/          ← vendedores + comisiones
│   │           ├── payment_methods/  ← métodos de pago
│   │           ├── fiscal_devices/   ← dispositivos fiscales
│   │           ├── credit_notes/     ← notas de crédito
│   │           ├── transfers/        ← transferencias
│   │           ├── manufacturing/    ← manufactura + BOM
│   │           ├── picking/          ← picking lists
│   │           ├── accounts/         ← CxC/CxP
│   │           ├── reports/          ← informes
│   │           ├── api_tokens/       ← API tokens
│   │           ├── scanner/          ← escáner de códigos
│   │           └── settings/         ← configuración
│   │
│   └── desktop/                      ← Tauri 2 (Fase 3)
│       └── src-tauri/
│           ├── src/fiscal/           ← Plugin Rust: impresoras fiscales
│           ├── Cargo.toml            ← Dependencias Rust
│           └── tauri.conf.json       ← Config + CSP
│
├── infra/
│   ├── docker/
│   │   └── compose.yml               ← Solo PowerSync (Postgres corre local)
│   ├── migrations/                   ← Goose SQL (00001–00009)
│   └── powersync/
│       ├── config.yaml               ← Conexión a Postgres + JWT secret
│       └── sync_rules.yaml           ← 50+ tablas con filtro por tenant
│
├── AGENTS.md                         ← Referencia rápida para AI agents
├── Makefile                          ← Comandos del proyecto
├── pnpm-workspace.yaml               ← Workspace definition
└── README.md                         ← Este archivo
```

<br>

## 🚀 Quick Start

### Prerrequisitos

| Herramienta | Versión | Instalación |
|-------------|---------|-------------|
| Go | 1.26+ | [go.dev/dl](https://go.dev/dl/) |
| Node.js | 20+ | [nodejs.org](https://nodejs.org/) |
| pnpm | 11+ | `npm install -g pnpm` |
| Flutter | 3.24+ | [flutter.dev](https://flutter.dev/) |
| PostgreSQL | 16+ | [postgresql.org](https://www.postgresql.org/download/) |
| Docker | — | [docker.com](https://www.docker.com/) (solo PowerSync) |
| goose | latest | `go install github.com/pressly/goose/v3/cmd/goose@latest` |

### 1. Habilitar replicación lógica en Postgres

Solo se hace una vez. Requiere reiniciar el servicio Postgres.

```sql
-- Ejecutar en psql como superusuario
ALTER SYSTEM SET wal_level = logical;
ALTER SYSTEM SET max_replication_slots = 10;
ALTER SYSTEM SET max_wal_senders = 10;
```

Después reiniciar Postgres y verificar:

```sql
SHOW wal_level;  -- debe retornar "logical"
```

### 2. Base de datos

```bash
make db-create
```

### 3. Variables de entorno

```bash
cp apps/api/.env.example apps/api/.env
cp apps/web/.env.example  apps/web/.env.local
```

Editar `apps/api/.env`:

```env
DATABASE_URL=postgresql://postgres:1234@localhost:5432/chiguire
JWT_SECRET=genera_con_openssl_rand_hex_32
POWERSYNC_JWT_SECRET=mismo_valor_que_en_infra/powersync/config.yaml
PORT=3001
```

> **Generar secrets seguros:**
> ```bash
> openssl rand -hex 32
> ```

### 4. Migraciones

```bash
make migrate
```

### 5. Iniciar servicios

```bash
# PowerSync (Docker)
make infra-up

# API Go
make api-dev          # → http://localhost:3001

# Web Next.js
make web-dev          # → http://localhost:3000

# Flutter
cd apps/mobile && flutter run
```

<br>

## 📋 Comandos

```bash
# Base de datos (Postgres local)
make db-create              # Crea la DB chiguire
make db-drop                # Elimina la DB (cuidado)
make db-wal                 # Habilita wal_level=logical (requiere restart)

# Migraciones
make migrate                # Aplica todas las pendientes
make migrate-down           # Revierte una migración
make migrate-status         # Ver estado actual
make migrate-create name=x  # Nueva migración vacía

# API Go
make api-dev                # Corre el servidor con go run
make api-build              # Compila binario en apps/api/bin/server

# Web Next.js
make web-dev                # Dev server en :3000
make web-build              # Build de producción

# PowerSync (Docker)
make infra-up               # Alias de powersync-up
make powersync-up           # Inicia PowerSync service
make powersync-down         # Detiene PowerSync service

# Flutter
make mobile-analyze         # flutter analyze
make mobile-test            # flutter test
```

<br>

## 📊 Fases de Desarrollo

### Fase 0 — Fundación ✅

**Objetivo:** probar que el sync multi-tenant funciona end-to-end antes de construir features.

- Monorepo: pnpm workspaces + Go + Flutter
- API Go: auth JWT, multi-tenant (tenants, users, user_tenants), PowerSync token endpoint
- PostgreSQL: Row-Level Security por `tenant_id` en cada tabla
- PowerSync: sync rules que aíslan datos por tenant (`token_parameters.tenant_id`)
- Next.js web: login, registro, dashboard con selector de empresa
- Flutter mobile: skeleton con PowerSync + Riverpod

---

### Fase 1 — Facturación e inventario ✅

**Objetivo:** un negocio puede vender y controlar stock completamente offline.

| Módulo | Features |
|--------|----------|
| **Clientes** | CRUD completo, RIF, búsqueda offline, historial de compras, vendedor asignado |
| **Productos** | Categorías, unidades de medida, código de barras, precios multi-moneda, costo promedio ponderado |
| **Inventario** | Múltiples almacenes, movimientos append-only, stock por ubicación (vista materializada), alertas de stock mínimo |
| **Facturación** | Forma libre, líneas con descuento, totales automáticos, estados `borrador → emitida → pagada → anulada`, numeración `XX-YYYYMMDD-NNNNN`, PDF en cliente |
| **Pagos básicos** | Efectivo USD/VES, transferencia, pagos parciales y combinados |

**Migraciones:** `00002_customers_vendors.sql`, `00003_products_inventory.sql`, `00004_invoices.sql`

---

### Fase 2 — Fiscal venezolano ✅

**Objetivo:** cumplimiento completo con SENIAT.

> ⚠️ Requiere revisión de un contador público venezolano antes de producción.

| Feature | Descripción |
|---------|-------------|
| **Tasa BCV** | Cron Go cada 6h scraping BCV → tabla `exchange_rates`. Fallback manual. |
| **IGTF** | 3% sobre pagos en divisas extranjeras. Se calcula al momento del pago. |
| **IVA** | General 16%, reducido 8%, exento (medicamentos). Categoría fiscal por producto. |
| **Retenciones IVA** | Agentes de retención: 75% del IVA. Comprobante PDF. |
| **Retenciones ISLR** | Por tipo de servicio/bien. Comprobante PDF. |
| **Libros fiscales** | Libro de Ventas y Compras con IVA desglosado. Export PDF/Excel. Períodos mensuales. |
| **Dispositivos fiscales** | CRUD + secuencias de documentos + talonarios de contingencia |
| **Notas crédito/débito** | Vinculadas a factura, con items y cálculo de totales |

**Migraciones:** `00005_fiscal.sql`, `00009_sellers_payment_methods.sql`

---

### Fase 3 — Impresora fiscal (Desktop) ✅

**Objetivo:** la app Tauri 2 imprime en impresoras fiscales homologadas por SENIAT.

| Marca | Protocolo | Status |
|-------|-----------|--------|
| **The Factory HKA** | Serial RS-232 / USB | ✅ Implementado |
| **PNP** | Serial / TCP | ✅ Implementado |
| **Bixolon** | USB | ✅ Implementado |

**Flujo:**
1. Factura llega a estado `approved`
2. Desktop envía al plugin Rust: líneas, totales, cliente, forma de pago
3. Plugin serializa en protocolo de la marca y escribe al puerto
4. Impresora retorna número Z y número de control
5. API registra el número de control SENIAT en la factura

**Modos de facturación SENIAT:**
- **Forma Libre**: PDF generado por la app (Fases 0–2)
- **Máquina Fiscal**: impresora conectada al desktop (Fase 3)
- **Imprenta Digital**: XML firmado + envío a proveedor SENIAT-autorizado (futuro)

---

### Fase 4 — Pagos ✅

**Objetivo:** el cajero cobra con cualquier método de pago digital venezolano desde la app.

| Pasarela | Método | Notas |
|----------|--------|-------|
| **Cashea** | Crédito digital en cuotas | Webhook confirma pago |
| **Spidi** | Pago móvil / transferencia | |
| **WayuPay** | Multi-método VE | |
| **Biopago** | Biométrico / punto de venta | Requiere terminal física |

**Arquitectura:**

```go
type PaymentProvider interface {
    CreatePaymentLink(invoice Invoice, amount Money) (PaymentLink, error)
    HandleWebhook(payload []byte, signature string) (PaymentEvent, error)
}
```

**Migración:** `00006_payments.sql`

---

### Fase 5 — Compras, cotizaciones y operaciones ✅

**Objetivo:** el ciclo completo de compra-venta y operaciones de campo.

| Módulo | Features |
|--------|----------|
| **Órdenes de compra** | Proveedor, productos, cantidades, aprobación `borrador → aprobada → recibida`, recepción parcial, stock automático |
| **Cotizaciones** | Igual que factura en estado `quotation`, convertible a factura, validez configurable, PDF por WhatsApp |
| **Comisiones** | Porcentaje por vendedor y/o categoría, cálculo sobre ventas cobradas, marcar como pagado, referencia de pago |
| **Rutas de reparto** | Lista de clientes por dirección, estado de entrega, firma digital, offline-first |
| **Transferencias** | Entre almacenes con ship/receive + `stock_movements` automáticos |
| **Manufactura** | Órdenes con BOM, start/complete con movimientos de stock automáticos |
| **Picking** | Listas de consolidación con verificación de items por cantidad |
| **Cuentas CxC/CxP** | Cuentas por cobrar/pagar con pagos parciales y recordatorios SMS/email |

**Migraciones:** `00007_purchases_quotations.sql`, `00009_sellers_payment_methods.sql`

---

### Fase 6 — Integraciones y SaaS ✅

**Objetivo:** producto SaaS con signup público, facturación propia y automatizaciones.

| Feature | Descripción |
|---------|-------------|
| **Import/Export Excel** | Clientes, productos desde Excel. Export de reportes. Plantillas descargables. |
| **Webhooks salientes** | URL configurable. Eventos: `invoice.created`, `invoice.paid`, `stock.low`. Cola con backoff exponencial. |
| **Suscripciones SaaS** | Trial 7 días → Emprendedor → PyME → Gold → Enterprise. Límites por plan. Cobro via Stripe. |
| **API Tokens** | Generación con hash SHA256, revocación, `last_used_at` |
| **Informes** | Libro ventas/compras, inventario actual/valorizado, kardex, art177, IGTF |
| **Onboarding** | Wizard de 5 pasos: empresa → sucursal → impuestos → producto → dispositivo fiscal |
| **Escáner** | Búsqueda de producto por barcode (mobile) |

**Migración:** `00008_saas_webhooks.sql`, `00009_sellers_payment_methods.sql`

<br>

## 🗄️ Modelo de Datos

### Tablas núcleo (multi-tenant)

```sql
-- Identificación de empresa
tenants (id, name, slug, plan, trial_ends_at)
users (id, email, password_hash, full_name)
user_tenants (user_id, tenant_id, role)

-- Catálogo
customers (tenant_id, rif, name, address, phone, email, default_seller_id)
vendors (tenant_id, rif, name, ...)
products (tenant_id, sku, name, category_id, tax_category, barcode)
product_prices (product_id, currency, amount, valid_from)

-- Inventario (append-only, nunca UPDATE de cantidad)
warehouses (tenant_id, branch_id, name)
stock_movements (tenant_id, warehouse_id, product_id, type, qty, reference_id)

-- Facturación
invoices (tenant_id, number, customer_id, seller_id, status, currency, subtotal, tax, total, issued_at)
invoice_items (invoice_id, product_id, qty, unit_price, discount, tax_amount)
invoice_payments (invoice_id, method, provider, amount_usd, amount_ves, rate_at_payment)

-- Fiscal
tax_categories (tenant_id, name, rate, type)
tax_withholdings (tenant_id, invoice_id, type, base, rate, amount)
fiscal_books (tenant_id, type, period, entries_json)
exchange_rates (currency, rate_to_ves, source, effective_at)  -- sin tenant_id (compartida)

-- Cachicamo parity (Fase 5-6)
sellers (tenant_id, user_id, name, commission_pct, is_active)
payment_methods (tenant_id, name, type, currency, ...)
fiscal_devices (tenant_id, serial, brand, model, status)
credit_notes (tenant_id, invoice_id, type, number, total)
inventory_transfers (tenant_id, from_warehouse, to_warehouse, status)
manufacturing_orders (tenant_id, product_id, qty, status)
picking_lists (tenant_id, status, assigned_to)
accounts_receivable (tenant_id, customer_id, amount, due_date)
accounts_payable (tenant_id, vendor_id, amount, due_date)
api_tokens (tenant_id, name, token_hash, last_used_at)
```

### Convención RLS

Todas las tablas de negocio tienen:

```sql
tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE
CREATE INDEX ON tabla(tenant_id);
ALTER TABLE tabla ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON tabla
    USING (tenant_id = current_setting('app.tenant_id', true)::uuid);
```

El API setea el contexto al inicio de cada request:

```go
conn.Exec(ctx, "SET LOCAL app.tenant_id = $1", tenantID)
```

<br>

## 🔐 Multi-tenancy

Un usuario puede pertenecer a varias empresas (ej: contador que maneja múltiples clientes). El JWT de acceso lleva el `tenant_id` de la empresa activa:

```json
{
  "user_id": "uuid",
  "tenant_id": "uuid",
  "role": "owner|admin|member|viewer",
  "exp": 1234567890
}
```

**Aislamiento garantizado en 3 capas:**

1. **PostgreSQL RLS** — policy `USING (tenant_id = current_setting('app.tenant_id', true)::uuid)`
2. **PowerSync sync rules** — `WHERE tenant_id = bucket.tenant_id` (filtro por JWT)
3. **API middleware** — `RequireTenant` setea `app.tenant_id` en cada request

<br>

## 🔧 Variables de Entorno

### API Go

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `DATABASE_URL` | ✅ | — | Connection string PostgreSQL |
| `JWT_SECRET` | ✅ | — | Secret para firmar JWT (HS256) |
| `POWERSYNC_JWT_SECRET` | ✅ | — | Secret para tokens de PowerSync |
| `PORT` | ❌ | `3001` | Puerto del servidor |

### Web (Next.js)

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | ✅ | `http://localhost:3001` | URL base del backend Go |
| `NEXT_PUBLIC_POWERSYNC_URL` | ✅ | `http://localhost:8080` | URL de PowerSync |

### Mobile (Flutter)

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `API_URL` | ✅ | `http://10.0.2.2:3001` | URL del backend (Android emulator) |
| `POWERSYNC_URL` | ✅ | `http://10.0.2.2:8080` | URL de PowerSync |

> Pasar via `--dart-define` en build flags.

<br>

## 🧪 Testing

```bash
# API Go
cd apps/api && go test ./...

# Flutter
cd apps/mobile && flutter test

# Web
cd apps/web && pnpm build
```

| Test | Cobertura |
|------|-----------|
| `invoice/handler_test.go` | Formato de número de factura |
| `payment/provider_test.go` | Providers de pago |
| `flutter test` | Widget test — login render |

<br>

## 📜 Licencia

**© 2026 Gustavo Colina (@Suggus1899). Todos los derechos reservados.**

Este software y su código fuente son **propiedad exclusiva** de Gustavo Colina (@Suggus1899).

- **No** está permitido copiar, modificar, distribuir, sublicenciar ni usar este código, total o parcialmente, sin autorización expresa y por escrito del autor.
- **No** está permitido usar este código con fines comerciales ni privados sin una licencia válida.
- Cualquier uso no autorizado constituye una violación de los derechos de autor y será perseguido conforme a la ley.

**Este es un software propietario. No es código abierto (open source) ni software libre.**

---

<div align="center">

<sub>Hecho con 🦔 para el sector empresarial de Venezuela</sub>

</div>
