package payment

import "errors"

// WayuPayProvider implements PaymentProvider for WayuPay.
// ponytail: stub returns error until API credentials are available; real impl
// needs WAYUPAY_API_KEY and HTTP POST to api.wayupay.app/charges
type WayuPayProvider struct{}

func (p *WayuPayProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return nil, errors.New("wayupay: not configured — set WAYUPAY_API_KEY to enable")
}

func (p *WayuPayProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return nil, errors.New("wayupay: not configured")
}

func init() {
	RegisterProvider("wayupay", &WayuPayProvider{})
}
