package gomerchant

type PaymentGateway interface {
	// The amount to charge as an integer (never a float). In the case of
	// currencies that support fractional amounts, should be the integer
	// amount of the smallest fractional (so, in the case of USD, the integer
	// number of cents).
	Purchase(amount uint64, pm *PaymentMethod, params *PurchaseParams) (PurchaseResponse, error)
	Authorize(amount uint64, pm *PaymentMethod, params *AuthorizeParams) (AuthorizeResponse, error)
	Capture(transactionID string, params *CaptureParams) (CaptureResponse, error)
	Void(transactionID string, params *VoidParams) (VoidResponse, error)
}
