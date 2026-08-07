package payment

// BiopagoProvider implements PaymentProvider for Biopago.
// TODO: integrate real Biopago API in production.
type BiopagoProvider struct{}

func (p *BiopagoProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return &PaymentLink{
		ProviderRef: "biopago-stub-" + invoiceID,
		AmountUSD:   amountUSD,
		Status:      "pending",
	}, nil
}

func (p *BiopagoProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return &PaymentEvent{
		ProviderRef: "biopago-stub",
		Status:      "paid",
		AmountUSD:   0,
	}, nil
}

func init() {
	RegisterProvider("biopago", &BiopagoProvider{})
}
