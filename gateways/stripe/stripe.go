// Package stripe implements GoMerchant payment gateway for Stripe.
package stripe

import (
	"fmt"

	"github.com/qor/gomerchant"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
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
	_, err := charge.Capture(transactionID, nil)
	return gomerchant.CaptureResponse{TransactionID: transactionID}, err
}

func (*Stripe) Refund(transactionID string, amount uint, params gomerchant.RefundParams) (gomerchant.RefundResponse, error) {
	refundParams := &stripe.RefundParams{
		Charge: transactionID,
		Amount: uint64(amount),
	}
	_, err := refund.New(refundParams)
	return gomerchant.RefundResponse{TransactionID: transactionID}, err
}

func (*Stripe) Void(transactionID string, params gomerchant.VoidParams) (gomerchant.VoidResponse, error) {
	refundParams := &stripe.RefundParams{
		Charge: transactionID,
	}
	_, err := refund.New(refundParams)
	return gomerchant.VoidResponse{TransactionID: transactionID}, err
}

func (*Stripe) Query(transactionID string) (gomerchant.Transaction, error) {
	c, err := charge.Get(transactionID, nil)
	transaction := gomerchant.Transaction{
		ID:        c.ID,
		Amount:    int(c.Amount),
		Currency:  string(c.Currency),
		Captured:  c.Captured,
		Paid:      c.Paid,
		Cancelled: c.Refunded,
		Status:    c.Status,
	}

	return transaction, err
}

func toStripeCC(customer string, cc *gomerchant.CreditCard, billingAddress *gomerchant.Address) *stripe.CardParams {
	cm := stripe.CardParams{
		Customer: customer,
		Name:     cc.Name,
		Number:   cc.Number,
		Month:    fmt.Sprint(cc.ExpMonth),
		Year:     fmt.Sprint(cc.ExpYear),
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
