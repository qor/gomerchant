package stripe_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/configor"
	"github.com/qor/gomerchant/gateways/stripe"
	"github.com/qor/gomerchant/tests"
	"github.com/stripe/stripe-go/customer"
)

var Stripe *stripe.Stripe

type Config struct {
	Key string `required:"true"`
}

func init() {
	var config = &Config{}
	os.Setenv("CONFIGOR_ENV_PREFIX", "-")
	if err := configor.Load(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Stripe = stripe.New(&stripe.Config{
		Key: config.Key,
	})
}

func TestTestSuite(t *testing.T) {
	tests.TestSuite{
		CreditCardManager: Stripe,
		Gateway:           Stripe,
		GetRandomCustomerID: func() string {
			Customer, err := customer.New(nil)
			if err != nil {
				fmt.Printf("Get error when create customer: %v", err)
			}
			return Customer.ID
		},
	}.TestAll(t)
}
