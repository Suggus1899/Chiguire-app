use crate::fiscal::protocol::{
    create_printer, FiscalInvoice, FiscalResult, PrinterConfig,
};
use tauri::{plugin::{Builder, TauriPlugin}, Runtime};

#[tauri::command]
fn print_fiscal_invoice(config: PrinterConfig, invoice: FiscalInvoice) -> FiscalResult {
    match create_printer(&config) {
        Ok(mut printer) => printer
            .print_invoice(&invoice)
            .unwrap_or_else(|e| FiscalResult::err(e)),
        Err(e) => FiscalResult::err(e),
    }
}

#[tauri::command]
fn print_z_report(config: PrinterConfig) -> FiscalResult {
    match create_printer(&config) {
        Ok(mut printer) => printer
            .print_z_report()
            .unwrap_or_else(|e| FiscalResult::err(e)),
        Err(e) => FiscalResult::err(e),
    }
}

#[tauri::command]
fn get_printer_status(config: PrinterConfig) -> String {
    match create_printer(&config) {
        Ok(mut printer) => printer.status().unwrap_or_else(|e| format!("{{\"error\":\"{e}\"}}")),
        Err(e) => format!("{{\"error\":\"{e}\"}}"),
    }
}

#[tauri::command]
fn list_printers() -> Vec<PrinterInfo> {
    let serial_ports = serialport::available_ports()
        .unwrap_or_default()
        .into_iter()
        .map(|p| p.port_name)
        .collect::<Vec<_>>();

    vec![
        PrinterInfo {
            printer_type: "hka".into(),
            label: "Factory HKA".into(),
            ports: serial_ports.clone(),
        },
        PrinterInfo {
            printer_type: "pnp".into(),
            label: "PNP".into(),
            ports: serial_ports,
        },
    ]
}

#[derive(serde::Serialize)]
struct PrinterInfo {
    #[serde(rename = "type")]
    printer_type: String,
    label: String,
    ports: Vec<String>,
}

pub fn init<R: Runtime>() -> TauriPlugin<R> {
    Builder::<R>::new("fiscal-printer")
        .invoke_handler(tauri::generate_handler![
            print_fiscal_invoice,
            print_z_report,
            get_printer_status,
            list_printers,
        ])
        .setup(|_app, _api| {
            log::info!("fiscal-printer plugin initialized");
            Ok(())
        })
        .build()
}
