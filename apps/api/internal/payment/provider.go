package payment

// PaymentProvider abstracts a payment gateway integration.
type PaymentProvider interface {
	CreatePaymentLink(invoiceID string, amountUSD float64) (*PaymentLink, error)
	HandleWebhook(payload []byte, signature string) (*PaymentEvent, error)
}

type PaymentLink struct {
	ID          string  `json:"id"`
	ShortCode   string  `json:"short_code"`
	ProviderRef string  `json:"provider_ref"`
	AmountUSD   float64 `json:"amount_usd"`
	Status      string  `json:"status"`
}

type PaymentEvent struct {
	ProviderRef string  `json:"provider_ref"`
	Status      string  `json:"status"`
	AmountUSD   float64 `json:"amount_usd"`
}

// providers maps provider name to implementation.
var providers = map[string]PaymentProvider{}

// RegisterProvider adds a payment provider implementation by name.
func RegisterProvider(name string, p PaymentProvider) {
	providers[name] = p
}

// GetProvider returns the provider for the given name, or nil.
func GetProvider(name string) PaymentProvider {
	return providers[name]
}
