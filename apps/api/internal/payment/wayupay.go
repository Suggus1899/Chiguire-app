package payment

// WayuPayProvider implements PaymentProvider for WayuPay.
// TODO: integrate real WayuPay API in production.
type WayuPayProvider struct{}

func (p *WayuPayProvider) CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error) {
	return &PaymentLink{
		ProviderRef: "wayupay-stub-" + invoiceID,
		AmountUSD:   amountUSD,
		Status:      "pending",
	}, nil
}

func (p *WayuPayProvider) HandleWebhook(payload []byte, signature string) (*PaymentEvent, error) {
	return &PaymentEvent{
		ProviderRef: "wayupay-stub",
		Status:      "paid",
		AmountUSD:   0,
	}, nil
}

func init() {
	RegisterProvider("wayupay", &WayuPayProvider{})
}
