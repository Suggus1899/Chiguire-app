package payment

import "errors"

// SpidiProvider implements PaymentProvider for Spidi.
// ponytail: stub returns error until API credentials are available; real impl
// needs SPIDI_API_KEY and HTTP POST to api.spidi.app/v1/payment-links
type SpidiProvider struct{}

func (p *SpidiProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return nil, errors.New("spidi: not configured — set SPIDI_API_KEY to enable")
}

func (p *SpidiProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return nil, errors.New("spidi: not configured")
}

func init() {
	RegisterProvider("spidi", &SpidiProvider{})
}
