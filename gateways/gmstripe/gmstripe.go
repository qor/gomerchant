package gmstripe

import (
	"github.com/qor/gomerchant"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
)

type Stripe struct{}

var _ gomerchant.PaymentGateway = Stripe{}

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
	_, ok := opt.(*PurchaseOptions)
	return ok
}

const ExtraKey = "stripe"

func (s Stripe) Purchase(amount uint64, pm *gomerchant.PaymentMethod, params *gomerchant.PurchaseParams) (gomerchant.PurchaseResponse, error) {
	var cparams *chargeParams
	if params != nil {
		cparams = &chargeParams{
			Amount:          params.Amount,
			Currency:        params.Currency,
			Customer:        params.Customer,
			Description:     params.Description,
			OrderID:         params.OrderID,
			BillingAddress:  params.BillingAddress,
			ShippingAddress: params.ShippingAddress,
			Extra:           params.Extra,
		}
	}

	scharge, err := makeCharge(amount, true, pm, cparams)
	var resp gomerchant.PurchaseResponse
	if scharge != nil {
		resp.TransactionID = scharge.ID
		resp.Extra.Set(ExtraKey, scharge)
	}
	return resp, err
}

type chargeParams struct {
	Amount          uint64
	Currency        string
	Customer        string
	Description     string
	OrderID         string
	BillingAddress  *gomerchant.Address
	ShippingAddress *gomerchant.Address
	gomerchant.Extra
}

func makeCharge(amount uint64, capture bool, pm *gomerchant.PaymentMethod, params *chargeParams) (*stripe.Charge, error) {
	cp := &stripe.ChargeParams{Amount: amount, NoCapture: !capture}

	if pm.Token != "" {
		cp.SetSource(pm.Token)
	} else if pm.CreditCard != nil {
		cp.SetSource(toStripeCC(pm.CreditCard, params))
	}

	if params != nil {
		cp.Customer = params.Customer
		cp.Currency = stripe.Currency(params.Currency)
		if extra, ok := params.Get(ExtraKey); ok && isPurchaseOptions(extra) {
			popt := extra.(*PurchaseOptions)
			cp.Params = popt.Params
			cp.Desc = popt.Desc
			cp.Statement = popt.Statement
			cp.Email = popt.Email
			cp.Dest = popt.Dest
		}
	}

	return charge.New(cp)
}

type CardOptions struct {
	Params  stripe.Params
	Default bool
}

func toStripeCC(cc *gomerchant.CreditCard, params *chargeParams) *stripe.CardParams {
	cm := stripe.CardParams{
		Name:   cc.Name,
		Number: cc.Number,
		Month:  cc.ExpMonth,
		Year:   cc.ExpYear,
		CVC:    cc.CVC,
	}

	if params == nil {
		return &cm
	}
	// cm.Currency = params.Currency
	cm.Customer = params.Customer

	if params.BillingAddress != nil {
		cm.Address1 = params.BillingAddress.Address1
		cm.Address2 = params.BillingAddress.Address2
		cm.City = params.BillingAddress.City
		cm.State = params.BillingAddress.State
		cm.Zip = params.BillingAddress.ZIP
		cm.Country = params.BillingAddress.Country
	}

	opts, ok := params.Get(ExtraKey)
	if !ok {
		return &cm
	}

	var co *CardOptions
	switch extra := opts.(type) {
	case *PurchaseOptions:
		co = extra.CardOptions
	default:
		return &cm
	}
	if co != nil {
		cm.Params = co.Params
		cm.Default = co.Default
	}

	return &cm
}

func (s Stripe) Authorize(amount uint64, pm *gomerchant.PaymentMethod, params *gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	var cparams *chargeParams
	if params != nil {
		cparams = &chargeParams{
			Amount:          params.Amount,
			Currency:        params.Currency,
			Customer:        params.Customer,
			Description:     params.Description,
			OrderID:         params.OrderID,
			BillingAddress:  params.BillingAddress,
			ShippingAddress: params.ShippingAddress,
			Extra:           params.Extra,
		}
	}

	scharge, err := makeCharge(amount, false, pm, cparams)
	var resp gomerchant.AuthorizeResponse
	if scharge != nil {
		resp.TransactionID = scharge.ID
		resp.Extra.Set(ExtraKey, scharge)
	}
	return resp, err
}

func (s Stripe) Capture(id string, params *gomerchant.CaptureParams) (gomerchant.CaptureResponse, error) {
	var sparams *stripe.CaptureParams
	if params != nil {
		if x, ok := params.Get(ExtraKey); ok {
			sparams = x.(*stripe.CaptureParams)
		}
		if sparams == nil {
			sparams = new(stripe.CaptureParams)
		}
		if params.Amount > 0 {
			sparams.Amount = params.Amount
		}
	}
	scharge, err := charge.Capture(id, sparams)
	var resp gomerchant.CaptureResponse
	if scharge != nil {
		resp.TransactionID = scharge.ID
		resp.Extra.Set(ExtraKey, scharge)
	}
	return resp, err
}

func (s Stripe) Void(id string, params *gomerchant.VoidParams) (gomerchant.VoidResponse, error) {
	var sparams *stripe.RefundParams
	if params != nil {
		if x, ok := params.Get(ExtraKey); ok {
			sparams = x.(*stripe.RefundParams)
		}
	}
	if sparams == nil {
		sparams = new(stripe.RefundParams)
	}
	sparams.Charge = id
	srefund, err := refund.New(sparams)
	var resp gomerchant.VoidResponse
	if srefund != nil {
		resp.TransactionID = srefund.ID
		resp.Extra.Set(ExtraKey, srefund)
	}
	return resp, err
}
