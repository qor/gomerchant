package gomerchant

import (
	"net/http"
	"time"
)

// PaymentGateway interface
type PaymentGateway interface {
	Authorize(amount uint64, params AuthorizeParams) (AuthorizeResponse, error)
	CompleteAuthorize(paymentID string, params CompleteAuthorizeParams) (CompleteAuthorizeResponse, error)
	Capture(transactionID string, params CaptureParams) (CaptureResponse, error)
	Refund(transactionID string, amount uint, params RefundParams) (RefundResponse, error)
	Void(transactionID string, params VoidParams) (VoidResponse, error)

	ConveniencePay(amount uint64, params ConveniencePayParams) (*ConveniencePayResponse, error)
	Query(transactionID string) (Transaction, error)
}

// AuthorizeParams authorize params
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

// AuthorizeResponse authorize response
type AuthorizeResponse struct {
	TransactionID  string
	HandleRequest  bool                                                   // need process request after authorize or not
	RequestHandler func(http.ResponseWriter, *http.Request, Params) error // process request
	Params
}

// CompleteAuthorizeParams complete authorize params
type CompleteAuthorizeParams struct {
	Params
}

// CompleteAuthorizeResponse complete authorize response
type CompleteAuthorizeResponse struct {
	Params
}

// CaptureParams capture params
type CaptureParams struct {
	Params
}

// CaptureResponse capture response
type CaptureResponse struct {
	TransactionID string
	Params
}

// RefundParams refund params
type RefundParams struct {
	Captured bool
	Params
}

// RefundResponse refund response
type RefundResponse struct {
	TransactionID string
	Params
}

// VoidParams void params
type VoidParams struct {
	Captured bool
	Params
}

// VoidResponse void response
type VoidResponse struct {
	TransactionID string
	Params
}

type CvsType uint

const (
	// 7-11
	CvsType_SevenEleven CvsType = iota + 1
	// lawson, ministop, family mart, daily yamazaki,
	CvsType_Lawson
	// seico mart
	CvsType_Seicomart
)

type ConveniencePayParams struct {
	CvsType            CvsType
	CustomerFamilyName string
	CustomerName       string
	// no '-'
	CustomerTel string
	// 0-60
	PaymentLimitDate *uint
}

type ConveniencePayResponse struct {
	ReceiptNumber    string
	PaymentID        string
	PrintURL         string
	PaymentLimitDate *time.Time
	Params
}
