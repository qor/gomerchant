package paygent_test

import (
	"testing"

	"github.com/qor/gomerchant/gateways/paygent"
)

var Paygent *paygent.Paygent

func init() {
	Paygent = paygent.New(&paygent.Config{})
}

func TestPurchase(t *testing.T) {
}
