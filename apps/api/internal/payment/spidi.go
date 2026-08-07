package payment

// SpidiProvider implements PaymentProvider for Spidi.
// TODO: integrate real Spidi API in production.
type SpidiProvider struct{}

func (p *SpidiProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return &PaymentLink{
		ProviderRef: "spidi-stub-" + invoiceID,
		AmountUSD:   amountUSD,
		Status:      "pending",
	}, nil
}

func (p *SpidiProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return &PaymentEvent{
		ProviderRef: "spidi-stub",
		Status:      "paid",
		AmountUSD:   0,
	}, nil
}

func init() {
	RegisterProvider("spidi", &SpidiProvider{})
}
