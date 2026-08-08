package payment

import "errors"

// CasheaProvider implements PaymentProvider for Cashea.
// ponytail: stub returns error until API credentials are available; real impl
// needs CASHEA_API_KEY env var and HTTP POST to cashea.com.ve/api/v1/links
type CasheaProvider struct{}

func (p *CasheaProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return nil, errors.New("cashea: not configured — set CASHEA_API_KEY to enable")
}

func (p *CasheaProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return nil, errors.New("cashea: not configured")
}

func init() {
	RegisterProvider("cashea", &CasheaProvider{})
}
