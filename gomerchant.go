package gomerchant

import "net/http"

// PaymentGateway interface
type PaymentGateway interface {
	Authorize(amount uint64, params AuthorizeParams) (AuthorizeResponse, error)
	Capture(transactionID string, params CaptureParams) (CaptureResponse, error)
	Refund(transactionID string, amount uint, params RefundParams) (RefundResponse, error)
	Void(transactionID string, params VoidParams) (VoidResponse, error)
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
	TransactionID  string
	HandleRequest  bool                                                   // need process request after authorize or not
	RequestHandler func(http.ResponseWriter, *http.Request, Params) error // process request
	Params
}

type CompleteAuthorizeParams struct {
	Params
}

type CompleteAuthorizeResponse struct {
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
	Captured bool
	Params
}

type RefundResponse struct {
	TransactionID string
	Params
}

// Void Params
type VoidParams struct {
	Captured bool
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
