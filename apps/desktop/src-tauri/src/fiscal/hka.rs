use crate::fiscal::protocol::{FiscalInvoice, FiscalPrinter, FiscalResult, PrinterConfig};
use serialport::SerialPort;
use std::io::{Read, Write};
use std::time::Duration;

const STX: u8 = 0x02;
const ETX: u8 = 0x03;
const ACK: u8 = 0x06;
const NAK: u8 = 0x15;

pub struct HkaPrinter {
    port_name: String,
    baud_rate: u32,
    port: Option<Box<dyn SerialPort>>,
}

impl HkaPrinter {
    pub fn new(port_name: impl Into<String>, baud_rate: u32) -> Self {
        Self {
            port_name: port_name.into(),
            baud_rate,
            port: None,
        }
    }

    pub fn from_config(cfg: &PrinterConfig) -> Result<Self, String> {
        Ok(Self::new(
            cfg.port_name.clone().unwrap_or_else(|| "COM1".into()),
            cfg.baud_rate.unwrap_or(9600),
        ))
    }

    fn open(&mut self) -> Result<(), String> {
        if self.port.is_some() {
            return Ok(());
        }
        let port = serialport::new(&self.port_name, self.baud_rate)
            .timeout(Duration::from_secs(3))
            .open()
            .map_err(|e| format!("open serial {0}: {e}", self.port_name))?;
        self.port = Some(port);
        Ok(())
    }

    // HKA frame: STX | cmd | data | ETX | checksum(2 bytes)
    fn build_frame(cmd: u8, data: &[u8]) -> Vec<u8> {
        let mut frame = Vec::with_capacity(data.len() + 5);
        frame.push(STX);
        frame.push(cmd);
        frame.extend_from_slice(data);
        frame.push(ETX);
        let checksum = frame[1..frame.len()].iter().fold(0u16, |acc, b| acc.wrapping_add(*b as u16));
        frame.push((checksum >> 8) as u8);
        frame.push((checksum & 0xff) as u8);
        frame
    }

    fn send(&mut self, cmd: u8, data: &[u8]) -> Result<Vec<u8>, String> {
        self.open()?;
        let port = self.port.as_mut().unwrap();
        let frame = Self::build_frame(cmd, data);
        port.write_all(&frame).map_err(|e| format!("write: {e}"))?;
        port.flush().map_err(|e| format!("flush: {e}"))?;

        let mut ack = [0u8; 1];
        port.read_exact(&mut ack).map_err(|e| format!("read ack: {e}"))?;
        if ack[0] != ACK {
            return Err(format!("printer NAK (0x{:02x})", ack[0]));
        }

        let mut resp = Vec::new();
        let mut buf = [0u8; 64];
        loop {
            let n = port.read(&mut buf).map_err(|e| format!("read resp: {e}"))?;
            if n == 0 {
                break;
            }
            resp.extend_from_slice(&buf[..n]);
            if resp.contains(&ETX) {
                break;
            }
        }
        Ok(resp)
    }

    // CMD 0x01: open fiscal invoice header (RIF,name)
    fn cmd_open_invoice(&mut self, invoice: &FiscalInvoice) -> Result<(), String> {
        let data = format!("{0}|{1}", invoice.customer_rif, invoice.customer_name);
        self.send(0x01, data.as_bytes())?;
        Ok(())
    }

    // CMD 0x02: add line item (desc|qty|price|tax_rate)
    fn cmd_add_item(&mut self, item: &crate::fiscal::protocol::FiscalItem) -> Result<(), String> {
        let data = format!(
            "{0}|{1}|{2}|{3}",
            item.description, item.qty, item.unit_price, item.tax_rate
        );
        self.send(0x02, data.as_bytes())?;
        Ok(())
    }

    // CMD 0x03: close invoice with payment method
    fn cmd_close_invoice(&mut self, payment_method: &str) -> Result<(), String> {
        self.send(0x03, payment_method.as_bytes())?;
        Ok(())
    }

    // CMD 0x04: request last invoice numbers (Z, control, invoice)
    fn cmd_get_numbers(&mut self) -> Result<(String, String, String), String> {
        let resp = self.send(0x04, &[])?;
        let text = resp
            .into_iter()
            .filter(|b| *b != STX && *b != ETX && *b != ACK)
            .map(|b| b as char)
            .collect::<String>();
        let parts: Vec<&str> = text.split('|').collect();
        if parts.len() < 3 {
            return Err("malformed numbers response".into());
        }
        Ok((parts[0].into(), parts[1].into(), parts[2].into()))
    }

    // CMD 0x05: Z report (cierre Z)
    fn cmd_z_report(&mut self) -> Result<(String, String, String), String> {
        let resp = self.send(0x05, &[])?;
        let text = resp
            .into_iter()
            .filter(|b| *b != STX && *b != ETX && *b != ACK)
            .map(|b| b as char)
            .collect::<String>();
        let parts: Vec<&str> = text.split('|').collect();
        if parts.len() < 3 {
            return Err("malformed Z response".into());
        }
        Ok((parts[0].into(), parts[1].into(), parts[2].into()))
    }

    // CMD 0x06: status query
    fn cmd_status(&mut self) -> Result<String, String> {
        let resp = self.send(0x06, &[])?;
        let text = resp
            .into_iter()
            .filter(|b| *b != STX && *b != ETX && *b != ACK)
            .map(|b| b as char)
            .collect::<String>();
        Ok(text)
    }
}

impl FiscalPrinter for HkaPrinter {
    fn print_invoice(&mut self, invoice: &FiscalInvoice) -> Result<FiscalResult, String> {
        self.cmd_open_invoice(invoice)?;
        for item in &invoice.items {
            self.cmd_add_item(item)?;
        }
        self.cmd_close_invoice(&invoice.payment_method)?;
        let (z, ctrl, inv) = self.cmd_get_numbers()?;
        Ok(FiscalResult::ok(z, ctrl, inv))
    }

    fn print_z_report(&mut self) -> Result<FiscalResult, String> {
        let (z, ctrl, inv) = self.cmd_z_report()?;
        Ok(FiscalResult::ok(z, ctrl, inv))
    }

    fn status(&mut self) -> Result<String, String> {
        self.cmd_status()
    }

    fn name(&self) -> &str {
        "Factory HKA"
    }
}
