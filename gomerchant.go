package gomerchant

type Gomerchant struct {
	PaymentGateway PaymentGateway
	Config         Config
}

type Config struct {
}

func New(paymentGateway PaymentGateway, config *Config) *Gomerchant {
	return &Gomerchant{PaymentGateway: paymentGateway, Config: config}
}

func (gomerchant *Gomerchant) Purchase(amount uint64, pm *PaymentMethod, params *PurchaseParams) (PurchaseResponse, error) {
	return gomerchant.Purchase(amount, pm, params)
}

func (gomerchant *Gomerchant) Authorize(amount uint64, pm *PaymentMethod, params *AuthorizeParams) (AuthorizeResponse, error) {
	return gomerchant.Authorize(amount, pm, params)
}

func (gomerchant *Gomerchant) Capture(transactionID string, params *CaptureParams) (CaptureResponse, error) {
	return gomerchant.Capture(amount, pm, params)
}

func (gomerchant *Gomerchant) Void(transactionID string, params *VoidParams) (VoidResponse, error) {
	return gomerchant.Void(amount, pm, params)
}
