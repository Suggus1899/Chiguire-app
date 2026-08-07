# Chiguire Desktop

Tauri 2 desktop wrapper for the Chiguire ERP web app, with a Rust plugin for Venezuelan fiscal printers (Factory HKA, PNP).

## Structure

```
apps/desktop/
├── package.json
├── src-tauri/
│   ├── Cargo.toml
│   ├── tauri.conf.json
│   ├── build.rs
│   └── src/
│       ├── main.rs
│       ├── lib.rs
│       └── fiscal/
│           ├── mod.rs
│           ├── plugin.rs      # Tauri commands (IPC)
│           ├── protocol.rs    # FiscalPrinter trait + types
│           ├── hka.rs         # Factory HKA (RS-232/USB)
│           └── pnp.rs         # PNP (Serial/TCP)
```

## Development

Start the Next.js web app first (`pnpm --filter web dev`), then:

```bash
cd apps/desktop
pnpm install
pnpm tauri dev
```

In production, `tauri build` bundles the built Next.js assets (`apps/web/out`).

## Fiscal printer plugin

Exposes Tauri commands to the frontend:

| Command | Description |
|---|---|
| `print_fiscal_invoice` | Print a fiscal invoice; returns Z, control and invoice numbers |
| `print_z_report` | Run cierre Z (daily close) |
| `get_printer_status` | Query printer status |
| `list_printers` | Enumerate available printer types and serial ports |

A `PrinterConfig` selects the backend:

```json
{ "type": "hka", "port_name": "COM1", "baud_rate": 9600 }
{ "type": "pnp", "tcp_host": "192.168.1.50", "tcp_port": 9100 }
```

### Supported printers

- **Factory HKA** — RS-232/USB, STX/ETX frames with 16-bit checksum.
- **PNP** — Serial or TCP, ENQ-prefixed frames with 8-bit checksum.

The byte-level protocol is structured but stubbed; real command bytes and checksums must be validated against each printer's manual before production use.
