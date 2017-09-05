// Package stripe implements GoMerchant payment gateway for Stripe.
package stripe

import (
	"github.com/qor/gomerchant"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

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
	stripe.Key = config.Key

	return &Stripe{
		Config: config,
	}
}

func (*Stripe) Authorize(amount uint64, params gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	chargeParams := &stripe.ChargeParams{
		Amount:    amount,
		Currency:  stripe.Currency(params.Currency),
		Desc:      params.Description,
		NoCapture: true,
	}
	chargeParams.AddMeta("order_id", params.OrderID)

	if params.PaymentMethod != nil {
		if params.PaymentMethod.CreditCard == nil {
			chargeParams.SetSource(toStripeCC(params.Customer, params.PaymentMethod.CreditCard, params.BillingAddress))
		}
		// TODO token
	}

	charge, err := charge.New(chargeParams)
	return gomerchant.AuthorizeResponse{TransactionID: charge.ID}, err
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

func toStripeCC(customer string, cc *gomerchant.CreditCard, billingAddress *gomerchant.Address) *stripe.CardParams {
	cm := stripe.CardParams{
		Customer: customer,
		Name:     cc.Name,
		Number:   cc.Number,
		Month:    cc.ExpMonth,
		Year:     cc.ExpYear,
		CVC:      cc.CVC,
	}

	if billingAddress != nil {
		cm.Address1 = billingAddress.Address1
		cm.Address2 = billingAddress.Address2
		cm.City = billingAddress.City
		cm.State = billingAddress.State
		cm.Zip = billingAddress.ZIP
		cm.Country = billingAddress.Country
	}

	return &cm
}
