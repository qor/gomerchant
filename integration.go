package gomerchant

import "net/http"

// IntegrationGateway integration gateway
type IntegrationGateway interface {
	CheckoutURL(params CheckoutParams) (string, error)
	VerifyNotification(req *http.Request) (bool, error)
	Refund(transactionID string, amount uint, params RefundParams) (RefundResponse, error)
	Void(transactionID string, params VoidParams) (VoidResponse, error)
	Query(transactionID string) (Transaction, error)
}

// CheckoutParams checkout params
type CheckoutParams struct {
	Amount      uint64
	Currency    string
	OrderID     string
	Description string
	Params
}
