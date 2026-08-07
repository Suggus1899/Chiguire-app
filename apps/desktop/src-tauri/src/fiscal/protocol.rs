use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FiscalInvoice {
    pub customer_rif: String,
    pub customer_name: String,
    pub items: Vec<FiscalItem>,
    pub total: f64,
    pub payment_method: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FiscalItem {
    pub description: String,
    pub qty: f64,
    pub unit_price: f64,
    pub tax_rate: f64,
    pub line_total: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FiscalResult {
    pub z_number: String,
    pub control_number: String,
    pub invoice_number: String,
    pub success: bool,
    pub error: Option<String>,
}

impl FiscalResult {
    pub fn ok(z: impl Into<String>, ctrl: impl Into<String>, inv: impl Into<String>) -> Self {
        Self {
            z_number: z.into(),
            control_number: ctrl.into(),
            invoice_number: inv.into(),
            success: true,
            error: None,
        }
    }

    pub fn err(msg: impl Into<String>) -> Self {
        Self {
            z_number: String::new(),
            control_number: String::new(),
            invoice_number: String::new(),
            success: false,
            error: Some(msg.into()),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PrinterConfig {
    #[serde(rename = "type")]
    pub printer_type: String,
    pub port_name: Option<String>,
    pub baud_rate: Option<u32>,
    pub tcp_host: Option<String>,
    pub tcp_port: Option<u16>,
}

pub trait FiscalPrinter {
    fn print_invoice(&mut self, invoice: &FiscalInvoice) -> Result<FiscalResult, String>;
    fn print_z_report(&mut self) -> Result<FiscalResult, String>;
    fn status(&mut self) -> Result<String, String>;
    fn name(&self) -> &str;
}

pub fn create_printer(config: &PrinterConfig) -> Result<Box<dyn FiscalPrinter>, String> {
    match config.printer_type.to_lowercase().as_str() {
        "hka" => Ok(Box::new(hka::HkaPrinter::from_config(config)?)),
        "pnp" => Ok(Box::new(pnp::PnpPrinter::from_config(config)?)),
        other => Err(format!("unknown printer type: {other}")),
    }
}
