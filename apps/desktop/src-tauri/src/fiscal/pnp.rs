use crate::fiscal::protocol::{FiscalInvoice, FiscalItem, FiscalPrinter, FiscalResult, PrinterConfig};
use serialport::SerialPort;
use std::io::{Read, Write};
use std::net::TcpStream;
use std::time::Duration;

const ENQ: u8 = 0x05;
const STX: u8 = 0x02;
const ETX: u8 = 0x03;
const ACK: u8 = 0x06;

enum Transport {
    Serial(Box<dyn SerialPort>),
    Tcp(TcpStream),
}

pub struct PnpPrinter {
    host: Option<String>,
    tcp_port: Option<u16>,
    port_name: Option<String>,
    baud_rate: u32,
    transport: Option<Transport>,
}

impl PnpPrinter {
    pub fn new_serial(port_name: impl Into<String>, baud_rate: u32) -> Self {
        Self {
            host: None,
            tcp_port: None,
            port_name: Some(port_name.into()),
            baud_rate,
            transport: None,
        }
    }

    pub fn new_tcp(host: impl Into<String>, port: u16) -> Self {
        Self {
            host: Some(host.into()),
            tcp_port: Some(port),
            port_name: None,
            baud_rate: 9600,
            transport: None,
        }
    }

    pub fn from_config(cfg: &PrinterConfig) -> Result<Self, String> {
        if let (Some(h), Some(p)) = (&cfg.tcp_host, cfg.tcp_port) {
            Ok(Self::new_tcp(h, p))
        } else if let Some(name) = &cfg.port_name {
            Ok(Self::new_serial(name, cfg.baud_rate.unwrap_or(9600)))
        } else {
            Err("pnp requires tcp_host+tcp_port or port_name".into())
        }
    }

    fn connect(&mut self) -> Result<(), String> {
        if self.transport.is_some() {
            return Ok(());
        }
        if let (Some(h), Some(p)) = (&self.host, self.tcp_port) {
            let stream = TcpStream::connect((h.as_str(), p))
                .map_err(|e| format!("tcp connect {h}:{p}: {e}"))?;
            stream
                .set_read_timeout(Some(Duration::from_secs(5)))
                .map_err(|e| format!("set_read_timeout: {e}"))?;
            stream
                .set_write_timeout(Some(Duration::from_secs(5)))
                .map_err(|e| format!("set_write_timeout: {e}"))?;
            self.transport = Some(Transport::Tcp(stream));
        } else if let Some(name) = &self.port_name {
            let port = serialport::new(name, self.baud_rate)
                .timeout(Duration::from_secs(5))
                .open()
                .map_err(|e| format!("open serial {name}: {e}"))?;
            self.transport = Some(Transport::Serial(port));
        } else {
            return Err("no transport configured".into());
        }
        Ok(())
    }

    // PNP frame: ENQ | STX | cmd | data | ETX | checksum(1 byte)
    fn build_frame(cmd: u8, data: &[u8]) -> Vec<u8> {
        let mut frame = Vec::with_capacity(data.len() + 5);
        frame.push(ENQ);
        frame.push(STX);
        frame.push(cmd);
        frame.extend_from_slice(data);
        frame.push(ETX);
        let checksum = frame[2..frame.len()]
            .iter()
            .fold(0u8, |acc, b| acc.wrapping_add(*b));
        frame.push(checksum);
        frame
    }

    fn write_all(&mut self, buf: &[u8]) -> Result<(), String> {
        match self.transport.as_mut().unwrap() {
            Transport::Serial(p) => {
                p.write_all(buf).map_err(|e| format!("serial write: {e}"))?;
                p.flush().map_err(|e| format!("serial flush: {e}"))
            }
            Transport::Tcp(s) => {
                s.write_all(buf).map_err(|e| format!("tcp write: {e}"))?;
                s.flush().map_err(|e| format!("tcp flush: {e}"))
            }
        }
    }

    fn read_exact(&mut self, n: usize) -> Result<Vec<u8>, String> {
        let mut out = vec![0u8; n];
        match self.transport.as_mut().unwrap() {
            Transport::Serial(p) => {
                p.read_exact(&mut out).map_err(|e| format!("serial read: {e}"))
            }
            Transport::Tcp(s) => {
                s.read_exact(&mut out).map_err(|e| format!("tcp read: {e}"))
            }
        }?;
        Ok(out)
    }

    fn read_until_etx(&mut self) -> Result<Vec<u8>, String> {
        let mut resp = Vec::new();
        let mut byte = [0u8; 1];
        loop {
            let n = match self.transport.as_mut().unwrap() {
                Transport::Serial(p) => p.read(&mut byte).map_err(|e| format!("serial read: {e}"))?,
                Transport::Tcp(s) => s.read(&mut byte).map_err(|e| format!("tcp read: {e}"))?,
            };
            if n == 0 {
                break;
            }
            resp.push(byte[0]);
            if byte[0] == ETX {
                break;
            }
        }
        Ok(resp)
    }

    fn send(&mut self, cmd: u8, data: &[u8]) -> Result<Vec<u8>, String> {
        self.connect()?;
        let frame = Self::build_frame(cmd, data);
        self.write_all(&frame)?;
        let ack = self.read_exact(1)?;
        if ack[0] != ACK {
            return Err(format!("NAK 0x{:02x}", ack[0]));
        }
        self.read_until_etx()
    }

    fn parse_response(resp: Vec<u8>) -> String {
        resp.into_iter()
            .filter(|b| !matches!(*b, ENQ | STX | ETX | ACK))
            .map(|b| b as char)
            .collect()
    }

    fn cmd_open(&mut self, inv: &FiscalInvoice) -> Result<(), String> {
        let data = format!("{0}|{1}", inv.customer_rif, inv.customer_name);
        self.send(0x10, data.as_bytes())?;
        Ok(())
    }

    fn cmd_item(&mut self, item: &FiscalItem) -> Result<(), String> {
        let data = format!(
            "{0}|{1}|{2}|{3}",
            item.description, item.qty, item.unit_price, item.tax_rate
        );
        self.send(0x11, data.as_bytes())?;
        Ok(())
    }

    fn cmd_close(&mut self, payment: &str) -> Result<(), String> {
        self.send(0x12, payment.as_bytes())?;
        Ok(())
    }

    fn cmd_numbers(&mut self) -> Result<(String, String, String), String> {
        let resp = self.send(0x13, &[])?;
        let text = Self::parse_response(resp);
        let parts: Vec<&str> = text.split('|').collect();
        if parts.len() < 3 {
            return Err("malformed numbers response".into());
        }
        Ok((parts[0].into(), parts[1].into(), parts[2].into()))
    }

    fn cmd_z(&mut self) -> Result<(String, String, String), String> {
        let resp = self.send(0x14, &[])?;
        let text = Self::parse_response(resp);
        let parts: Vec<&str> = text.split('|').collect();
        if parts.len() < 3 {
            return Err("malformed Z response".into());
        }
        Ok((parts[0].into(), parts[1].into(), parts[2].into()))
    }

    fn cmd_status(&mut self) -> Result<String, String> {
        let resp = self.send(0x15, &[])?;
        Ok(Self::parse_response(resp))
    }
}

impl FiscalPrinter for PnpPrinter {
    fn print_invoice(&mut self, invoice: &FiscalInvoice) -> Result<FiscalResult, String> {
        self.cmd_open(invoice)?;
        for item in &invoice.items {
            self.cmd_item(item)?;
        }
        self.cmd_close(&invoice.payment_method)?;
        let (z, ctrl, inv) = self.cmd_numbers()?;
        Ok(FiscalResult::ok(z, ctrl, inv))
    }

    fn print_z_report(&mut self) -> Result<FiscalResult, String> {
        let (z, ctrl, inv) = self.cmd_z()?;
        Ok(FiscalResult::ok(z, ctrl, inv))
    }

    fn status(&mut self) -> Result<String, String> {
        self.cmd_status()
    }

    fn name(&self) -> &str {
        "PNP"
    }
}
