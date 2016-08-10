package gmstripe

import (
	"github.com/qor/gomerchant"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type Stripe struct{}

func NewStripe(key string) *Stripe {
	stripe.Key = key
	return new(Stripe)
}

type PurchaseOptions struct {
	Params                       stripe.Params
	Desc, Statement, Email, Dest string
	CardOptions                  *CardOptions
}

func isPurchaseOptions(opt interface{}) bool {
	if opt == nil {
		return false
	}
	_, ok := opt.(PurchaseOptions)
	return ok
}

func (s *Stripe) Purchase(amount uint64, pm *gomerchant.PaymentMethod, opts *gomerchant.Options) (gomerchant.Response, error) {
	cp := &stripe.ChargeParams{Amount: amount}

	if pm.Token != "" {
		cp.SetSource(pm.Token)
	} else if pm.CreditCard != nil {
		cp.SetSource(toStripeCC(pm.CreditCard, opts))
	}

	if opts != nil {
		cp.Customer = opts.Customer
		cp.Currency = opts.Currency
		if isPurchaseOptions(opts.Extra) {
			popt := opts.Extra.(PurchaseOptions)
			cp.Params = popt.Params
			cp.Desc = popt.Desc
			cp.Statement = popt.Statement
			cp.Email = popt.Email
			cp.Dest = popt.Dest
		}
	}

	scharge, err := charge.New(cp)
	var resp gomerchant.Response
	if scharge != nil {
		resp.ID = scharge.ID
		resp.Extra = scharge
	}

	return resp, err
}

type CardOptions struct {
	Params  stripe.Params
	Default bool
}

func toStripeCC(cc *gomerchant.CreditCard, opts *gomerchant.Options) *stripe.CardParams {
	cm := stripe.CardParams{
		Name:   cc.Name,
		Number: cc.Number,
		Month:  cc.ExpMonth,
		Year:   cc.ExpYear,
		CVC:    cc.CVC,
	}

	if opts == nil {
		return &cm
	}
	cm.Currency = opts.Currency
	cm.Customer = opts.Customer

	if opts.BillingAddress != nil {
		cm.Address1 = opts.BillingAddress.Address1
		cm.Address2 = opts.BillingAddress.Address2
		cm.City = opts.BillingAddress.City
		cm.State = opts.BillingAddress.State
		cm.Zip = opts.BillingAddress.ZIP
		cm.Country = opts.BillingAddress.Country
	}

	if opts.Extra == nil {
		return &cm
	}
	co, ok := opts.Extra.(CardOptions)
	if !ok {
		return &cm
	}
	cm.Params = co.Params
	cm.Default = co.Default

	return &cm
}

func (s *Stripe) Authorize(amount int, pm *gomerchant.PaymentMethod, opts *gomerchant.Options) (gomerchant.Response, error) {
}

func (s *Stripe) Capture(amount int, id string, opts *gomerchant.Options) (gomerchant.Response, error) {
}

func (s *Stripe) Void(id string, opts *gomerchant.Options) (gomerchant.Response, error) {
}

func (s *Stripe) Store(pm *gomerchant.PaymentMethod, opts *gomerchant.Options) (gomerchant.Response, error) {
}
