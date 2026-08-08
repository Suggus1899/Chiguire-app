package payment

import "errors"

// BiopagoProvider implements PaymentProvider for Biopago.
// ponytail: stub returns error until API credentials are available; real impl
// needs BIOPAGO_API_KEY and HTTP POST to api.biopago.com/v1/transactions
type BiopagoProvider struct{}

func (p *BiopagoProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return nil, errors.New("biopago: not configured — set BIOPAGO_API_KEY to enable")
}

func (p *BiopagoProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return nil, errors.New("biopago: not configured")
}

func init() {
	RegisterProvider("biopago", &BiopagoProvider{})
}
