package gomerchant

import "net/http"

// PaymentGateway interface
type PaymentGateway interface {
	Authorize(amount uint64, params AuthorizeParams) (AuthorizeResponse, error)
	CompleteAuthorize(paymentID string, params CompleteAuthorizeParams) (CompleteAuthorizeResponse, error)
	Capture(transactionID string, params CaptureParams) (CaptureResponse, error)
	Refund(transactionID string, amount uint, params RefundParams) (RefundResponse, error)
	Void(transactionID string, params VoidParams) (VoidResponse, error)

	Query(transactionID string) (Transaction, error)
	//RakutePayApplicationMessage(amount uint64, params RakutenPayApplicationParams) (RakutenPayApplicationResponse, error)
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

type InquiryResponse struct {
	TransactionID     string
	TradingID         string
	PaymentNoticeID   string
	PaymentInitDate   string
	PaymentChangeDate string
	PaymentAmount     string
	RelatedPaymentID  string
	PaymentStatus     string
	SuccessCode       string
	SuccessDetail     string
	Params
}

type RakutenPayApplicationParams struct {
	MerchandiseType uint64
	PCMobileType    uint64
	ButtonType      string
	ReturnUrl       string
	CancelUrl       string
	Goods           []Good
	Params
}

type RakutenPayApplicationResponse struct {
	TransactionID       string
	OrderCode           string
	TradeGenerationDate string
	RedirectHTML        string
	Params
}

type Good struct {
	ID     string
	Name   string
	Price  float64
	Amount uint64
}

const RAKUTEN_PAY_PRODUCT_ID = "WholeOrderAmount"
