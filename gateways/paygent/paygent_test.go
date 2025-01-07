package paygent_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/configor"
	"github.com/qor/gomerchant"
	"github.com/qor/gomerchant/gateways/paygent"
	"github.com/qor/gomerchant/tests"
)

var Paygent *paygent.Paygent

type Config struct {
	MerchantID      string `required:"true"`
	ConnectID       string `required:"true"`
	ConnectPassword string `required:"true"`
	TelegramVersion string `required:"true" default:"1.0"`
	MerchantName    string `required:"true"`

	ClientFilePath string `required:"true" default:"paygent.pem"`
	CertPassword   string `required:"true" default:"changeit"`
	CAFilePath     string `required:"true" default:"curl-ca-bundle.crt"`

	ProductionMode  bool
	SecurityCodeUse bool
}

func init() {
	var config = &Config{}
	if err := configor.New(&configor.Config{ENVPrefix: "PAYGENT_CONFIG"}).Load(config); err != nil {
		fmt.Println(config)
		os.Exit(1)
	}

	Paygent = paygent.New(&paygent.Config{
		MerchantID:      config.MerchantID,
		MerchantName:    config.MerchantName,
		ConnectID:       config.ConnectID,
		ConnectPassword: config.ConnectPassword,
		ClientFilePath:  config.ClientFilePath,
		CertPassword:    config.CertPassword,
		CAFilePath:      config.CAFilePath,
		ProductionMode:  config.ProductionMode,
		SecurityCodeUse: config.SecurityCodeUse,
	})
}

func TestTestSuite(t *testing.T) {
	tests.TestSuite{
		CreditCardManager: Paygent,
		Gateway:           Paygent,
		GetRandomCustomerID: func() string {
			return fmt.Sprint(time.Now().Unix())
		},
	}.TestAll(t)
}

func Test3DAuthorizeAndCapture(t *testing.T) {
	cards := map[string]bool{
		"5123459358515820": true,
		"5123459358515821": false,
	}

	for card, is3D := range cards {
		authorizeResult, err := Paygent.SecureCodeAuthorize(100,
			paygent.SecureCodeParams{
				UserAgent:  "User-Agent	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0.2 Safari/602.3.12",
				TermURL:    "http://getqor.com/order/return",
				HttpAccept: "http",
			},
			gomerchant.AuthorizeParams{
				Currency: "JPY",
				OrderID:  fmt.Sprint(time.Now().Unix()),
				PaymentMethod: &gomerchant.PaymentMethod{
					CreditCard: &gomerchant.CreditCard{
						Name:     "JCB Card",
						Number:   card,
						ExpMonth: 1,
						ExpYear:  uint(time.Now().Year() + 1),
						CVC:      "1234",
					},
				},
			})

		if err != nil || authorizeResult.TransactionID == "" {
			t.Error(err, authorizeResult)
		}

		if is3D != authorizeResult.HandleRequest {
			t.Errorf("HandleRequest for card %v should be %v, but got %v", card, is3D, authorizeResult.HandleRequest)
		}

		if is3D {
			if result, ok := authorizeResult.Get("out_acs_html"); !ok || result.(string) == "" {
				t.Errorf("should get HTML, but %v", authorizeResult)
			}
		}
	}
}

func TestStart3DS2Authentication(t *testing.T) {
	// for new creditcard
	res, err := Paygent.Start3DS2Authentication(context.Background(), gomerchant.Start3DS2AuthenticationParams{
		OrderID: fmt.Sprint(time.Now().Unix()),
		TermURL: "http://getqor.com/order/return",
		Amount:  10,
		PaymentMethod: &gomerchant.PaymentMethod{
			CreditCard: &gomerchant.CreditCard{
				Name:     "JCB Card",
				Number:   "5123459358515821",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
				CVC:      "1234",
			},
		},
	})
	if err != nil {
		t.Error("for new creditcard: ", err, res)
		return
	}
	if res.OutAcsHTML == "" {
		t.Error("out_acs_html is empty")
		return
	}
	if res.Result != "0" {
		t.Error("result should be 0")
		return
	}
	t.Logf("%+v", res)

	// for saved creditcard
	customerID := "customerid111aigletest"
	response, err := Paygent.ListCreditCards(gomerchant.ListCreditCardsParams{CustomerID: customerID})
	if err != nil {
		t.Errorf("failed to ListCreditCards err: %+v", err)
		return
	}
	if len(response.CreditCards) == 0 {
		t.Errorf("no saved credit cards for customer %v", customerID)
		return
	}
	res, err = Paygent.Start3DS2Authentication(context.Background(), gomerchant.Start3DS2AuthenticationParams{
		OrderID: fmt.Sprint(time.Now().Unix()),
		TermURL: "https://dev-lacoste-frontend.aldt.theplant-dev.com/",
		Amount:  10,
		PaymentMethod: &gomerchant.PaymentMethod{
			SavedCreditCard: &gomerchant.SavedCreditCard{
				CustomerID:   "customerid111aigletest",
				CreditCardID: response.CreditCards[0].CreditCardID,
			},
		},
	})
	if err != nil {
		t.Errorf("for saved creditcard err: %+v, res: %+v", err, res)
		return
	}
	if res.OutAcsHTML == "" {
		t.Error("out_acs_html is empty")
		return
	}
	if res.Result != "0" {
		t.Error("result should be 0")
		return
	}
	t.Logf("%+v", res)
}

func Test3DS2Authorization(t *testing.T) {
	resp, err := Paygent.Authorize(200000, gomerchant.AuthorizeParams{
		Currency: "JPY",
		OrderID:  fmt.Sprint(time.Now().Unix()),
		PaymentMethod: &gomerchant.PaymentMethod{
			SavedCreditCard: &gomerchant.SavedCreditCard{
				CustomerID:   "customerid111aigletest",
				CreditCardID: "14340385",
				// CVC:          "1234",
				ThreeDSAuthID: "6b1f2a0e-fe7f-4ccd-98de-7af69cdf996d",
			},
		},
	})
	if err != nil {
		t.Errorf("error: %+v", err)
	}
	t.Logf("result: %+v", resp)
}
