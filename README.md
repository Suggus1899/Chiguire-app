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
<a href="https://github.com/pressly/goose" target="_blank"><img src="https://github.com/pressly.png?size=96" width="48" height="48" alt="goose" /></a>
<br><sub><b><a href="https://github.com/pressly/goose" target="_blank">goose</a></b></sub>
<br><sub>Migrations</sub>
</td>
</tr>
<tr>
<th colspan="5" align="center" width="600"><sub><b>Sync · Desktop · Tooling</b></sub></th>
</tr>
<tr>
<td align="center" width="120">
<a href="https://www.powersync.com/" target="_blank"><img src="https://github.com/powersync-ja.png?size=96" width="48" height="48" alt="PowerSync" /></a>
<br><sub><b><a href="https://www.powersync.com/" target="_blank">PowerSync</a></b></sub>
<br><sub>Offline-first sync</sub>
</td>
<td align="center" width="120">
<a href="https://tauri.app/" target="_blank"><img src="https://cdn.simpleicons.org/tauri/FFC131" width="48" height="48" alt="Tauri" /></a>
<br><sub><b><a href="https://tauri.app/" target="_blank">Tauri 2</a></b></sub>
<br><sub>Desktop · Rust</sub>
</td>
<td align="center" width="120">
<a href="https://riverpod.dev/" target="_blank"><img src="https://cdn.simpleicons.org/riverpod/4D7DFB" width="48" height="48" alt="Riverpod" /></a>
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

<br>

## 📱 Aplicaciones

| App | Descripción | Puerto | Plataforma |
|-----|-------------|--------|------------|
| **API** | Go + chi + pgx/v5 — auth, facturación, fiscal, pagos, compras, SaaS | `3001` | Backend |
| **Web** | Next.js 16 PWA — 33 rutas, dashboard, facturas, inventario, fiscal, informes | `3000` | Web · PWA instalable |
| **Mobile** | Flutter + Riverpod — 30+ pantallas, offline-first, escáner de códigos | — | Android · iOS |
| **Desktop** | Tauri 2 + Rust — impresora fiscal (The Factory HKA, PNP, Bixolon) | — | Windows · macOS · Linux |

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
