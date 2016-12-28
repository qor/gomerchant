package gomerchant

// PaymentGateway interface
type PaymentGateway interface {
	Purchase(amount uint64, params *PurchaseParams) (PurchaseResponse, error)
	Authorize(amount uint64, params *AuthorizeParams) (AuthorizeResponse, error)
	Capture(transactionID string, params *CaptureParams) (CaptureResponse, error)
	Refund(transactionID string, params *RefundParams) (RefundResponse, error)
	Void(transactionID string, params *VoidParams) (VoidResponse, error)
}

// Purchase Params
type PurchaseParams struct {
	Amount          uint64
	Currency        string
	Customer        string
	Description     string
	OrderID         string
	BillingAddress  *Address
	ShippingAddress *Address
	PaymentMethod   *PaymentMethod
	Params
}

type PurchaseResponse struct {
	TransactionID string
	Params
}

// Authorize Params
type AuthorizeParams struct {
	Amount          uint64
	Currency        string
	Customer        string
	Description     string
	OrderID         string
	BillingAddress  *Address
	ShippingAddress *Address
	PaymentMethod   *PaymentMethod
	Params
}

type AuthorizeResponse struct {
	TransactionID string
	Params
}

// Capture Params
type CaptureParams struct {
	Params
}

type CaptureResponse struct {
	TransactionID string
	Params
}

// Refund Params
type RefundParams struct {
	Params
}

type RefundResponse struct {
	TransactionID string
	Params
}

// Void Params
type VoidParams struct {
	Params
}

type VoidResponse struct {
	TransactionID string
	Params
}

// CreateCreditCard Params
type CreateCreditCardParams struct {
	CustomerID string
	CreditCard *CreditCard
}

type CreditCardParamsResponse struct {
	CustomerID   string
	CreditCardID string
	Params
}
