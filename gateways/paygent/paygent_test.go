package paygent_test

import (
	"testing"

	"github.com/qor/gomerchant"
	"github.com/qor/gomerchant/gateways/paygent"
)

var Paygent *paygent.Paygent

func init() {
	Paygent = paygent.New(&paygent.Config{ClientFilePath: "paygent.pem"})
}

func TestPurchase(t *testing.T) {
	Paygent.Client()
	Paygent.Authorize(100, &gomerchant.AuthorizeParams{})
}
