# Chiguire

ERP offline-first multi-tenant para empresas venezolanas. Facturación, inventario, fiscal SENIAT, pagos locales (Cashea, Spidi, WayuPay, Biopago) y sincronización en la nube. Funciona sin internet; sincroniza cuando hay conexión.

---

## Contenido

- [¿Por qué Chiguire?](#por-qué-chiguire)
- [Arquitectura](#arquitectura)
- [Stack tecnológico](#stack-tecnológico)
- [Estructura del proyecto](#estructura-del-proyecto)
- [Configuración local](#configuración-local)
- [Comandos](#comandos)
- [Fases de desarrollo](#fases-de-desarrollo)
  - [Fase 0 — Fundación](#fase-0--fundación-completada)
  - [Fase 1 — Facturación e inventario](#fase-1--facturación-e-inventario)
  - [Fase 2 — Fiscal venezolano](#fase-2--fiscal-venezolano)
  - [Fase 3 — Impresora fiscal (Desktop)](#fase-3--impresora-fiscal-desktop)
  - [Fase 4 — Pagos](#fase-4--pagos)
  - [Fase 5 — Compras, cotizaciones y operaciones](#fase-5--compras-cotizaciones-y-operaciones)
  - [Fase 6 — Integraciones y SaaS](#fase-6--integraciones-y-saas)
- [Modelo de datos](#modelo-de-datos)
- [Multi-tenancy](#multi-tenancy)
- [Sync offline-first](#sync-offline-first)
- [Convenciones](#convenciones)

---

## ¿Por qué Chiguire?

Venezuela tiene necesidades ERP específicas que las soluciones internacionales no cubren bien:

- **Multi-moneda real**: USD como moneda funcional, VES con tasa BCV actualizada automáticamente.
- **Compliance SENIAT**: IGTF 3%, IVA 16%, retenciones ISLR/IVA, libros de ventas/compras, impresoras fiscales homologadas.
- **Pagos locales**: Cashea, Spidi, WayuPay, Biopago, además de transferencias y efectivo.
- **Conectividad intermitente**: el negocio no puede detenerse cuando se va la luz o la fibra. Chiguire funciona completamente offline y sincroniza en cuanto hay conexión.

---

## Arquitectura

```
┌──────────────────────────────────────────────────────────┐
│  Clientes                                                │
│                                                          │
│  Flutter (móvil)   Next.js PWA (web)   Tauri (desktop)  │
│       │                  │                   │           │
│  SQLite (Drift)    SQLite (WASM)     SQLite (nativo)    │
│       └──────────────────┴───────────────────┘           │
│                          │ PowerSync SDK                 │
└──────────────────────────┼───────────────────────────────┘
                           │ sync bidireccional, incremental
┌──────────────────────────┼───────────────────────────────┐
│  PowerSync Service       │                               │
│  (replica Postgres → SQLite, filtra por tenant)          │
└──────────────────────────┼───────────────────────────────┘
                           │
┌──────────────────────────┼───────────────────────────────┐
│  PostgreSQL 16           │                               │
│  Row-Level Security por tenant_id                        │
└──────────────────────────┼───────────────────────────────┘
                           │
┌──────────────────────────┼───────────────────────────────┐
│  API Go (chi)            │                               │
│  Auth JWT · Emisión fiscal · Pagos · Numeración          │
└──────────────────────────────────────────────────────────┘
```

### Reglas offline-first

| Tipo de operación | Flujo |
|---|---|
| Leer datos (clientes, productos, facturas) | SQLite local. Nunca espera red. |
| Crear/editar borrador (factura, cliente, producto) | SQLite local → PowerSync sincroniza a Postgres. |
| Emitir factura (número correlativo SENIAT) | `POST /invoices/emit` → API asigna número → PowerSync propaga. Se encola offline; emite al reconectar. |
| Cobrar con pasarela (Cashea, etc.) | `POST /payments` → API. Se marca `pending_payment` offline. |
| Stock | Solo movimientos (`stock_movements`, append-only). Nunca se modifica la cantidad absoluta. |

---

## Stack tecnológico

| Capa | Tecnología | Notas |
|---|---|---|
| Base de datos | PostgreSQL 16 | `wal_level=logical` para replicación |
| Sync engine | PowerSync (self-hosted) | Postgres → SQLite por tenant |
| API | Go 1.25 + chi + pgx/v5 + sqlc | Monolito. Migraciones con goose. |
| Web | Next.js 16 App Router + TypeScript + Tailwind CSS + shadcn/ui | PWA instalable |
| State web | TanStack Query + Zustand | |
| Móvil | Flutter + Riverpod + Drift | Android / iOS |
| Desktop | Tauri 2 + Rust | Envuelve el web. Plugin Rust para impresoras fiscales. |
| Auth | JWT propios (access 15 min + refresh 7 días) + bcrypt | Multi-tenant: `tenant_id` en el claim |
| Tasa BCV | Cron Go → tabla `exchange_rates` | Fallback: admin setea manualmente |
| Pago SaaS | Stripe | Suscripciones de Chiguire mismo |
| CI | GitHub Actions | lint + test + build por app |

---

## Estructura del proyecto

```
chiguire/
├── apps/
│   ├── api/                        # Go API
│   │   ├── cmd/server/             # main + handlers de fase 0
│   │   └── internal/
│   │       ├── auth/               # JWT, register, login, refresh
│   │       ├── middleware/         # RequireAuth, RequireTenant
│   │       ├── tenant/             # CRUD de empresas
│   │       ├── powersync/          # Token endpoint para el sync
│   │       ├── invoice/            # (Fase 1)
│   │       ├── inventory/          # (Fase 1)
│   │       ├── fiscal/             # (Fase 2) IGTF, IVA, retenciones
│   │       └── payment/            # (Fase 4) Cashea, Spidi, etc.
│   │
│   ├── web/                        # Next.js PWA
│   │   ├── app/
│   │   │   ├── (auth)/             # login, register
│   │   │   └── (app)/              # dashboard, facturas, inventario...
│   │   └── lib/
│   │       ├── api.ts              # Cliente HTTP tipado
│   │       └── powersync.ts        # Schema SQLite + instancia DB local
│   │
│   ├── mobile/                     # Flutter
│   │   └── lib/
│   │       ├── core/               # database.dart, api.dart
│   │       └── features/           # (Fase 1+) por feature
│   │
│   └── desktop/                    # Tauri 2 (Fase 3)
│       └── src-tauri/
│           └── src/
│               └── fiscal/         # Plugin Rust para impresoras fiscales
│
├── infra/
│   ├── docker/
│   │   └── compose.yml             # Solo PowerSync (Postgres corre local)
│   ├── migrations/                 # Goose SQL
│   └── powersync/
│       ├── config.yaml             # Conexión a Postgres + JWT secret
│       └── sync_rules.yaml         # Reglas de sync por tenant
│
├── packages/
│   └── shared-types/               # (futuro) tipos TS generados desde SQL
│
├── AGENTS.md                       # Referencia rápida para AI agents
├── Makefile                        # Comandos del proyecto
└── README.md                       # Este archivo
```

---

## Configuración local

### Requisitos

| Herramienta | Versión mínima |
|---|---|
| Go | 1.25+ |
| Node.js | 20+ |
| pnpm | 9+ |
| Flutter | 3.24+ |
| PostgreSQL | 16 (instalado localmente) |
| Docker Desktop | Solo para PowerSync |
| goose | `go install github.com/pressly/goose/v3/cmd/goose@latest` |

### 1. Habilitar replicación lógica en Postgres

Solo se hace una vez. Requiere reiniciar el servicio Postgres.

```sql
-- Ejecutar en psql como superusuario
ALTER SYSTEM SET wal_level = logical;
ALTER SYSTEM SET max_replication_slots = 10;
ALTER SYSTEM SET max_wal_senders = 10;
```

Después reiniciar Postgres (Windows: `Restart-Service postgresql-x64-16` o desde Services.msc).

Verificar:

```sql
SHOW wal_level;  -- debe retornar "logical"
```

### 2. Crear la base de datos

```bash
make db-create
```

O manualmente:

```sql
CREATE DATABASE chiguire;
```

### 3. Variables de entorno

```bash
cp apps/api/.env.example apps/api/.env
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

### 4. Ejecutar migraciones

```bash
make migrate
```

### 5. Iniciar el API

```bash
make api-dev
```

El API queda en `http://localhost:3001`.

### 6. Iniciar la web

```bash
make web-dev
```

La web queda en `http://localhost:3000`.

### 7. PowerSync (opcional en Fase 0)

El sync end-to-end requiere PowerSync corriendo. Solo corre en Docker.

Editar `infra/powersync/config.yaml` con el mismo `POWERSYNC_JWT_SECRET` del `.env`.

```bash
make powersync-up
```

PowerSync queda en `http://localhost:8080`.

---

## Comandos

```bash
# Base de datos (Postgres local)
make db-create              # Crea la DB chiguire
make db-drop                # Elimina la DB (cuidado)
make db-wal                 # Habilita wal_level=logical (requiere restart)

# Migraciones
make migrate                # Aplica todas las pendientes
make migrate-down           # Revierte una migración
make migrate-status         # Ver estado actual
make migrate-create name=add_customers   # Nueva migración vacía

# API Go
make api-dev                # Corre el servidor con go run
make api-build              # Compila binario en apps/api/bin/server

# Web Next.js
make web-dev                # Dev server en :3000
make web-build              # Build de producción

# PowerSync (Docker)
make powersync-up           # Inicia PowerSync service
make powersync-down         # Detiene PowerSync service

# Flutter
make mobile-analyze         # flutter analyze
make mobile-test            # flutter test
```

---

## Fases de desarrollo

### Fase 0 — Fundación ✅ Completada

**Objetivo:** probar que el sync multi-tenant funciona end-to-end antes de construir features.

**Entregables:**
- Monorepo: pnpm workspaces + Go + Flutter
- API Go: auth JWT, multi-tenant (tenants, users, user_tenants), PowerSync token endpoint
- PostgreSQL: Row-Level Security por `tenant_id` en cada tabla
- PowerSync: sync rules que aíslan datos por tenant (`token_parameters.tenant_id`)
- Next.js web: login, registro, dashboard con selector de empresa y test de sync
- Flutter mobile: skeleton con PowerSync + Riverpod
- Migración base: `tenants`, `users`, `user_tenants`, `refresh_tokens`, `exchange_rates`, `test_items`

**Test crítico de aislamiento:**

```bash
# Crear dos usuarios con empresas distintas
POST /auth/register { email: "user1@test.com", ... }
POST /auth/register { email: "user2@test.com", ... }

# Cada uno crea su empresa y agrega test_items
# Verificar que /test-items de user1 NO retorna items de user2
```

---

### Fase 1 — Facturación e inventario

**Objetivo:** un negocio puede vender y controlar stock completamente offline.

**Módulos:**

#### Clientes y proveedores
- CRUD completo (nombre, RIF, dirección, teléfono, email)
- Búsqueda offline desde SQLite local
- Historial de compras por cliente

#### Productos y catálogo
- Categorías, unidades de medida, código de barras/QR
- Precios en múltiples monedas (USD base, VES calculado)
- Foto del producto (almacenada localmente, sincronizada en background)
- Costo promedio ponderado

#### Inventario
- Múltiples almacenes por sucursal
- Movimientos append-only (entrada, salida, transferencia, ajuste)
- Stock por ubicación calculado desde movimientos (vista materializada)
- Alertas de stock mínimo

#### Facturación — Forma Libre
- Encabezado: cliente, fecha, vendedor, sucursal
- Líneas: producto, cantidad, precio, descuento
- Totales: subtotal, descuento, impuestos, total
- Estados: `borrador → emitida → pagada → anulada`
- Numeración: `XX-YYYYMMDD-NNNNN` (correlativo por sucursal, asignado por API)
- PDF generado en cliente (sin necesidad de servidor)
- Funciona completamente offline; emite el número al reconectar

#### Pagos básicos
- Efectivo USD, efectivo VES, transferencia bancaria
- Pagos parciales y combinados
- Saldo pendiente por factura

**Nuevas tablas:**
```sql
customers, vendors, product_categories, products, product_prices,
warehouses, stock_movements, stock_by_location (materialized view),
invoices, invoice_items, invoice_payments
```

**Migraciones:** `00002_customers_vendors.sql`, `00003_products_inventory.sql`, `00004_invoices.sql`

---

### Fase 2 — Fiscal venezolano

**Objetivo:** cumplimiento completo con SENIAT. Ninguna factura se emite sin cálculos correctos.

> ⚠️ Esta fase requiere revisión de un contador público venezolano antes de lanzar a producción.

#### Tasa BCV automática
- Cron Go cada 6 horas scraping BCV
- Tabla `exchange_rates` con historial
- Fallback: admin setea tasa manualmente si el scraper falla
- Todas las facturas graban la tasa del momento de emisión

#### IGTF (Impuesto a las Grandes Transacciones Financieras)
- 3% sobre pagos en divisas extranjeras (USD, EUR, etc.)
- Se calcula al momento del pago, no de la factura
- Retención opcional si el cliente es contribuyente especial

#### IVA
- Alícuota general: 16%
- Alícuota reducida: 8% (algunos alimentos y bienes)
- Exento: medicamentos, ciertos alimentos básicos
- Cada producto tiene su categoría fiscal

#### Retenciones IVA
- Agentes de retención retienen el 75% del IVA (o 100% en casos especiales)
- Comprobante de retención generado en PDF
- Registro en `tax_withholdings`

#### Retenciones ISLR
- Aplica a proveedores según actividad económica
- Porcentaje configurable por tipo de servicio/bien
- Comprobante de retención

#### Libros fiscales
- Libro de Ventas: todas las facturas del período con IVA desglosado
- Libro de Compras: todas las compras con IVA acreditable
- Exportación a PDF y Excel (formato SENIAT)
- Períodos mensuales

**Nuevas tablas:**
```sql
tax_categories, tax_withholdings, fiscal_periods, fiscal_books
```

---

### Fase 3 — Impresora fiscal (Desktop)

**Objetivo:** la app Tauri 2 puede imprimir en impresoras fiscales homologadas por SENIAT.

#### Tauri 2 desktop
- Envuelve el mismo build de Next.js (un solo frontend)
- Bundle ~10MB vs Electron ~120MB
- Autoactualización via Tauri updater

#### Plugin de impresoras fiscales
Plugin Rust (`src-tauri/src/fiscal/`) que abstrae el protocolo de cada marca:

| Marca | Protocolo | Status |
|---|---|---|
| The Factory HKA | Serial RS-232 / USB | Fase 3 (prioridad — más común en VE) |
| PNP | Serial / TCP | Fase 3 |
| Bixolon | USB | Fase 4 |
| Otras | TBD | Bajo demanda |

**Flujo de emisión:**
1. Factura llega a estado `approved` en la app
2. Desktop envía al plugin: líneas, totales, cliente, forma de pago
3. Plugin serializa en protocolo de la marca y escribe al puerto
4. Impresora retorna número Z y número de control
5. API registra el número de control SENIAT en la factura

**Modos de facturación SENIAT:**
- **Forma Libre**: PDF generado por la app (Fases 0–2)
- **Máquina Fiscal**: impresora conectada al desktop (Fase 3)
- **Imprenta Digital**: XML firmado + envío a proveedor SENIAT-autorizado (Fase 4+)

---

### Fase 4 — Pagos

**Objetivo:** el cajero puede cobrar con cualquier método de pago digital venezolano desde la app.

#### Pasarelas integradas

| Pasarela | Método | Notas |
|---|---|---|
| **Cashea** | Crédito digital en cuotas | Webhook confirma pago |
| **Spidi** | Pago móvil / transferencia | |
| **WayuPay** | Multi-método VE | |
| **Biopago** | Biométrico / punto de venta | Requiere terminal física |

**Arquitectura de pagos:**

```go
// Cada pasarela implementa esta interfaz
type PaymentProvider interface {
    CreatePaymentLink(invoice Invoice, amount Money) (PaymentLink, error)
    HandleWebhook(payload []byte, signature string) (PaymentEvent, error)
}
```

Un solo archivo por pasarela (`internal/payment/cashea.go`, etc.). El router de webhooks enruta por `provider` en la URL: `POST /webhooks/payments/{provider}`.

**Flujo offline:**
- Si no hay red al cobrar: se registra el intento localmente
- Al reconectar: se genera el link y se envía por WhatsApp/SMS al cliente
- Si el cliente pagó offline (efectivo): se registra directamente sin pasarela

#### Links de pago
- Se genera un link corto (`chiguire.app/pay/{code}`) por factura
- El cliente abre en su teléfono, escoge método, paga
- Webhook actualiza el estado de la factura en tiempo real

---

### Fase 5 — Compras, cotizaciones y operaciones

**Objetivo:** el ciclo completo de compra-venta y operaciones de campo.

#### Órdenes de compra
- Selección de proveedor, productos, cantidades
- Aprobación en flujo (borrador → aprobada → recibida)
- Recepción parcial de mercancía
- Actualización automática de stock al recibir

#### Cotizaciones
- Igual que una factura pero en estado `quotation`
- Convertible a factura con un clic
- Validez configurable (días)
- PDF enviable por WhatsApp

#### Comisiones de vendedores
- Porcentaje configurable por vendedor y/o por categoría de producto
- Cálculo sobre ventas cobradas (no solo facturadas)
- Reporte de comisiones por período

#### Rutas de reparto
- Lista de clientes ordenada por dirección (ruta del día)
- Estado de entrega por factura (pendiente, entregada, rechazada)
- Firma digital del cliente al recibir
- Funciona completamente offline; sincroniza al volver a la oficina

**Nuevas tablas:**
```sql
purchase_orders, purchase_order_items, purchase_order_receipts,
quotations, sales_commissions, delivery_routes, delivery_stops
```

---

### Fase 6 — Integraciones y SaaS

**Objetivo:** producto SaaS con signup público, facturación propia y automatizaciones.

#### Import/Export Excel
- Importar clientes, productos desde Excel (SheetJS / Apache POI)
- Exportar cualquier reporte a Excel
- Plantillas descargables para importación masiva

#### Webhooks salientes
- El cliente configura una URL
- Chiguire envía eventos: `invoice.created`, `invoice.paid`, `stock.low`, etc.
- Cola de reintentos con backoff exponencial

#### N8N templates
- Templates pre-configurados para automatizaciones comunes:
  - Factura pagada → enviar comprobante por email
  - Stock bajo → notificar al encargado por WhatsApp
  - Nuevo cliente → agregar a CRM

#### Suscripciones SaaS (Chiguire como producto)
- Planes: Trial (7 días) → Emprendedor → PyME → Gold → Enterprise
- Límites por plan: facturas/mes, sucursales, usuarios
- Cobro via Stripe (tarjeta internacional)
- Al vencer el plan: modo solo lectura (no bloqueo total)

#### Onboarding
- Wizard de 5 pasos: empresa → moneda → sucursales → primer producto → primera factura
- Datos de ejemplo precargados opcionales

---

## Modelo de datos

### Tablas núcleo (multi-tenant)

```sql
-- Identificación de empresa
tenants (id, name, slug, plan, trial_ends_at)
users (id, email, password_hash, full_name)
user_tenants (user_id, tenant_id, role)

-- Catálogo
customers (tenant_id, rif, name, address, phone, email)
vendors (tenant_id, rif, name, ...)
products (tenant_id, sku, name, category_id, tax_category)
product_prices (product_id, currency, amount, valid_from)

-- Inventario (append-only, nunca UPDATE de cantidad)
warehouses (tenant_id, branch_id, name)
stock_movements (tenant_id, warehouse_id, product_id, type, qty, reference_id)

-- Facturación
invoices (tenant_id, number, customer_id, status, currency, subtotal, tax, total, issued_at)
invoice_items (invoice_id, product_id, qty, unit_price, discount, tax_amount)
invoice_payments (invoice_id, method, provider, amount_usd, amount_ves, rate_at_payment)

-- Fiscal
tax_withholdings (tenant_id, invoice_id, type, base, rate, amount)
fiscal_books (tenant_id, type, period, entries_json)
exchange_rates (currency, rate_to_ves, source, effective_at)  -- sin tenant_id (compartida)
```

### Convención RLS

Todas las tablas de negocio tienen:

```sql
tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE
-- + índice: CREATE INDEX ON tabla(tenant_id);
-- + policy: USING (tenant_id = current_setting('app.tenant_id', true)::uuid)
```

El API setea el contexto al inicio de cada request:

```go
conn.Exec(ctx, "SET LOCAL app.tenant_id = $1", tenantID)
```

---

## Multi-tenancy

Un usuario puede pertenecer a varias empresas (ej: contador que maneja múltiples clientes). El JWT de acceso lleva el `tenant_id` de la empresa activa:

```json
{
  "user_id": "uuid",
  "tenant_id": "uuid",
  "role": "owner|admin|member|viewer",
  "exp": 1234567890
}
```

**Para cambiar de empresa:** el usuario hace login nuevamente pasando el nuevo `tenant_id`. El frontend guarda el token en sessionStorage (dev) / httpOnly cookie (prod).

**Token PowerSync:** el endpoint `GET /powersync/token` emite un token de 30 min con `parameters.tenant_id` como claim. PowerSync usa ese valor para filtrar qué filas replica a ese cliente.

---

## Sync offline-first

### ¿Cómo funciona PowerSync?

1. PowerSync lee el WAL (Write-Ahead Log) de Postgres via replicación lógica.
2. Aplica las `sync_rules.yaml` para decidir qué filas van a qué bucket.
3. Cada cliente se suscribe a su bucket (filtrado por `tenant_id`).
4. Los cambios llegan al cliente como operaciones SQLite incrementales.
5. El cliente escribe en su SQLite local y la UI reacciona reactivamente.

### Escrituras desde el cliente

Las escrituras del cliente van primero a SQLite local (instantáneo). PowerSync las propaga a Postgres via la API:

- Chiguire usa el modo **"managed uploads"**: el SDK envía un `PUT /powersync/upload` con las operaciones pendientes.
- El API valida, aplica RLS, escribe en Postgres.
- PowerSync propaga el cambio a todos los demás clientes del mismo tenant.

### Conflictos

| Caso | Resolución |
|---|---|
| Dos usuarios editan el mismo cliente | Last-Write-Wins (PowerSync por defecto) |
| Stock modificado simultáneamente | No hay conflicto: solo se hacen inserts en `stock_movements` |
| Número de factura asignado | Lo asigna exclusivamente el API (nunca el cliente) |
| Pago procesado offline | Se encola; se envía al reconectar |

---

## Convenciones

### Commits (Conventional Commits)

```
feat: agregar módulo de cotizaciones
fix: corregir cálculo de IGTF en pagos combinados
chore: actualizar dependencias Go
docs: documentar endpoints de fiscal
refactor: extraer lógica de PDF a helper
test: agregar tests de retención ISLR
```

### Branches

```
main          → producción
dev           → integración
feat/xxx      → features nuevas
fix/xxx       → bugfixes
```

### APIs Go

- Handlers en `internal/{feature}/handler.go`
- Lógica de negocio en `internal/{feature}/service.go`
- Queries SQL en `internal/db/queries/{feature}.sql` (sqlc genera el Go)
- Una interfaz por pasarela de pago; un archivo por implementación

### Nuevo módulo checklist

Al agregar una nueva tabla:

- [ ] Migración goose en `infra/migrations/`
- [ ] `tenant_id` en la tabla + índice
- [ ] Policy RLS `USING (tenant_id = ...)`
- [ ] Agregar tabla a `sync_rules.yaml` si debe sincronizar al cliente
- [ ] Agregar tabla al schema SQLite en `apps/web/lib/powersync.ts` y `apps/mobile/lib/core/database.dart`
- [ ] Queries sqlc en `internal/db/queries/`
- [ ] Handler + service en `internal/{feature}/`
- [ ] Rutas registradas en `cmd/server/main.go`
