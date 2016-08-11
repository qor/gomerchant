package gomerchant

type PaymentGateway interface {
	Purchase(amount uint64, pm *PaymentMethod, params *PurchaseParams) (PurchaseResponse, error)
	Authorize(amount uint64, pm *PaymentMethod, params *AuthorizeParams) (AuthorizeResponse, error)
	Capture(transactionID string, params *CaptureParams) (CaptureResponse, error)
	Void(transactionID string, params *VoidParams) (VoidResponse, error)
}
