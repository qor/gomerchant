package paygent_test

import (
	"fmt"
	"testing"

	"github.com/qor/gomerchant"
	"github.com/qor/gomerchant/gateways/paygent"
)

var Paygent *paygent.Paygent

func init() {
	Paygent = paygent.New(&paygent.Config{ClientFilePath: "paygent.pem", CAFilePath: "curl-ca-bundle.crt", CertPassword: "changeit"})
}

func TestPurchase(t *testing.T) {
	fmt.Println(Paygent.Request("094", nil))
	Paygent.Authorize(100, &gomerchant.AuthorizeParams{})
}
