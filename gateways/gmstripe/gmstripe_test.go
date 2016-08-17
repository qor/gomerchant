package gmstripe

import (
	"os"
	"testing"

	"github.com/stripe/stripe-go/card"

	"github.com/stripe/stripe-go/customer"

	"github.com/stripe/stripe-go/charge"

	stripe "github.com/stripe/stripe-go"

	"github.com/qor/gomerchant"
)

var cardToken string
var customerToken string
var key string

func init() {
	key = os.Getenv("STRIPE_KEY")
	if key == "" {
		panic("please specify STRIPE_KEY")
	}
	stripe.Key = key

	cus, err := customer.New(&stripe.CustomerParams{Email: "paralyzed.horse@test.com"})
	if err != nil {
		panic(err)
	}
	customerToken = cus.ID
	scard, err := card.New(&stripe.CardParams{
		Customer: cus.ID,
		Name:     "Paralyzed Horse",
		Number:   "4242424242424242",
		Month:    "07",
		Year:     "2020",
		CVC:      "314",
	})
	if err != nil {
		panic(err)
	}
	cardToken = scard.ID
}

func TestStripe(t *testing.T) {
	gmstripe := NewStripe(key)

	for i, pm := range []gomerchant.PaymentMethod{
		{Token: cardToken},
		{
			CreditCard: &gomerchant.CreditCard{
				Name:     "Paralyzed Horse",
				Number:   "4242424242424242",
				ExpMonth: "07",
				ExpYear:  "2020",
				CVC:      "314",
			},
		},
	} {
		t.Log("case", i)
		{
			t.Log("Purchase")
			var params gomerchant.PurchaseParams
			if pm.Token != "" {
				params.Customer = customerToken
			}
			params.Currency = os.Getenv("STRIPE_CURRENCY")
			params.Set(ExtraKey, &PurchaseOptions{
				Statement: "Testing Inc.",
				Params: stripe.Params{
					Meta: map[string]string{"id": "test"},
				},
			})
			resp, err := gmstripe.Purchase(1000, &pm, &params)
			if err != nil {
				t.Error(err)
			} else if _, err := charge.Get(resp.TransactionID, nil); err != nil {
				t.Error(err)
			}

			if _, err := gmstripe.Void(resp.TransactionID, nil); err != nil {
				t.Error(err)
			}
		}

		{
			t.Log("Auth & Capture")
			var params gomerchant.AuthorizeParams
			if pm.Token != "" {
				params.Customer = customerToken
			}
			params.Currency = os.Getenv("STRIPE_CURRENCY")
			params.Set(ExtraKey, &PurchaseOptions{
				Statement: "Testing Inc.",
				Params: stripe.Params{
					Meta: map[string]string{"id": "test"},
				},
			})
			aresp, err := gmstripe.Authorize(1000, &pm, &params)
			if err != nil {
				t.Error(err)
			} else if _, err := charge.Get(aresp.TransactionID, nil); err != nil {
				t.Error(err)
			}
			cresp, err := gmstripe.Capture(aresp.TransactionID, nil)
			if err != nil {
				t.Error(err)
			} else if _, err := charge.Get(cresp.TransactionID, nil); err != nil {
				t.Error(err)
			}
		}
	}
}
