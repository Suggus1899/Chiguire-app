package payment

// CasheaProvider implements PaymentProvider for Cashea.
// TODO: integrate real Cashea API in production.
type CasheaProvider struct{}

func (p *CasheaProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return &PaymentLink{
		ProviderRef: "cashea-stub-" + invoiceID,
		AmountUSD:   amountUSD,
		Status:      "pending",
	}, nil
}

func (p *CasheaProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return &PaymentEvent{
		ProviderRef: "cashea-stub",
		Status:      "paid",
		AmountUSD:   0,
	}, nil
}

func init() {
	RegisterProvider("cashea", &CasheaProvider{})
}
