package gomerchant

import "net/http"

// PaymentGateway interface
type PaymentGateway interface {
	Authorize(amount uint64, params AuthorizeParams) (AuthorizeResponse, error)
	CompleteAuthorize(transactionID string, params CompleteAuthorizeParams) (CompleteAuthorizeResponse, error)
	Capture(transactionID string, params CaptureParams) (CaptureResponse, error)

	Refund(transactionID string, amount uint, params RefundParams) (RefundResponse, error)
	Void(transactionID string, params VoidParams) (VoidResponse, error)

	Query(transactionID string) (Transaction, error)
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
	RawBody        string
	Params
}

// Complete Authorize
type CompleteAuthorizeParams struct {
	Params
}

type CompleteAuthorizeResponse struct {
	RawBody string
	Params
}

// Capture Params
type CaptureParams struct {
	Params
}

type CaptureResponse struct {
	TransactionID string
	RawBody       string
	Params
}

// Refund Params
type RefundParams struct {
	Captured bool
	Params
}

type RefundResponse struct {
	TransactionID string
	RawBody       string
	Params
}

// Void Params
type VoidParams struct {
	Captured bool
	Params
}

type VoidResponse struct {
	TransactionID string
	RawBody       string
	Params
}
