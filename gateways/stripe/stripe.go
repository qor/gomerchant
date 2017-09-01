// Package stripe implements GoMerchant payment gateway for Stripe.
package stripe

import "github.com/qor/gomerchant"

// Stripe implements gomerchant.PaymetGateway interface.
type Stripe struct {
	Config *Config
}

var _ gomerchant.PaymentGateway = &Stripe{}

// Config stripe config
type Config struct {
	Key string
}

// New creates Stripe struct.
func New(config *Config) *Stripe {
	return &Stripe{
		Config: config,
	}
}

func (*Stripe) Authorize(amount uint64, params gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	return gomerchant.AuthorizeResponse{}, nil
}

func (*Stripe) CompleteAuthorize(paymentID string, params gomerchant.CompleteAuthorizeParams) (gomerchant.CompleteAuthorizeResponse, error) {
	return gomerchant.CompleteAuthorizeResponse{}, nil
}

func (*Stripe) Capture(transactionID string, params gomerchant.CaptureParams) (gomerchant.CaptureResponse, error) {
	return gomerchant.CaptureResponse{}, nil
}

func (*Stripe) Refund(transactionID string, amount uint, params gomerchant.RefundParams) (gomerchant.RefundResponse, error) {
	return gomerchant.RefundResponse{}, nil
}

func (*Stripe) Void(transactionID string, params gomerchant.VoidParams) (gomerchant.VoidResponse, error) {
	return gomerchant.VoidResponse{}, nil
}

func (*Stripe) Query(transactionID string) (gomerchant.Transaction, error) {
	return gomerchant.Transaction{}, nil
}
