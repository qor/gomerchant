package paygent_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/configor"
	"github.com/qor/gomerchant"
	"github.com/qor/gomerchant/gateways/paygent"
)

var (
	Paygent *paygent.Paygent
)

type Config struct {
	MerchantID      string `required:"true"`
	ConnectID       string `required:"true"`
	ConnectPassword string `required:"true"`
	TelegramVersion string `required:"true" default:"1.0"`

	ClientFilePath string `required:"true" default:"paygent.pem"`
	CertPassword   string `required:"true" default:"changeit"`
	CAFilePath     string `required:"true" default:"curl-ca-bundle.crt"`

	ProductionMode bool
}

func init() {
	var config = &Config{}
	os.Setenv("CONFIGOR_ENV_PREFIX", "-")
	if err := configor.Load(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Paygent = paygent.New(&paygent.Config{
		MerchantID:      config.MerchantID,
		ConnectID:       config.ConnectID,
		ConnectPassword: config.ConnectPassword,
		ClientFilePath:  config.ClientFilePath,
		CertPassword:    config.CertPassword,
		CAFilePath:      config.CAFilePath,
		TelegramVersion: config.TelegramVersion,
		ProductionMode:  config.ProductionMode,
	})
}

func createSavedCreditCard() (gomerchant.CreditCardParamsResponse, error) {
	return Paygent.CreateCreditCard(gomerchant.CreateCreditCardParams{
		CustomerID: fmt.Sprint(time.Now().Unix()),
		CreditCard: &gomerchant.CreditCard{
			Name:     "JCB Card",
			Number:   "3580876521284076",
			ExpMonth: 1,
			ExpYear:  uint(time.Now().Year() + 1),
		},
	})
}

func TestCreateCreditCard(t *testing.T) {
	if result, err := createSavedCreditCard(); err != nil || result.CreditCardID == "" {
		t.Error(err, result)
	}
}

func TestAuthorizeAndCapture(t *testing.T) {
	authorizeResult, err := Paygent.Authorize(100, gomerchant.AuthorizeParams{
		Currency: "JPY",
		OrderID:  fmt.Sprint(time.Now().Unix()),
		PaymentMethod: &gomerchant.PaymentMethod{
			CreditCard: &gomerchant.CreditCard{
				Name:     "JCB Card",
				Number:   "3580876521284076",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
			},
		},
	})

	if err != nil || authorizeResult.TransactionID == "" {
		t.Error(err, authorizeResult)
	}

	captureResult, err := Paygent.Capture(authorizeResult.TransactionID, gomerchant.CaptureParams{})

	if err != nil || captureResult.TransactionID == "" {
		t.Error(err, captureResult)
	}
}

func TestAuthorizeAndCaptureWithSavedCreditCard(t *testing.T) {
	if savedCreditCard, err := createSavedCreditCard(); err == nil {
		authorizeResult, err := Paygent.Authorize(100, gomerchant.AuthorizeParams{
			Currency: "JPY",
			OrderID:  fmt.Sprint(time.Now().Unix()),
			PaymentMethod: &gomerchant.PaymentMethod{
				SavedCreditCard: &gomerchant.SavedCreditCard{
					CustomerID:   savedCreditCard.CustomerID,
					CreditCardID: savedCreditCard.CreditCardID,
				},
			},
		})

		if err != nil || authorizeResult.TransactionID == "" {
			t.Error(err, authorizeResult)
		}

		captureResult, err := Paygent.Capture(authorizeResult.TransactionID, gomerchant.CaptureParams{})

		if err != nil || captureResult.TransactionID == "" {
			t.Error(err, captureResult)
		}
	}
}

func Test3DAuthorizeAndCapture(t *testing.T) {
	cards := map[string]bool{
		"5123459358515820": true,
		"5123459358515821": false,
	}

	for card, is3D := range cards {
		authorizeResult, err := Paygent.SecureCodeAuthorize(100,
			paygent.SecureCodeParams{
				UserAgent: "User-Agent	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0.2 Safari/602.3.12",
				TermURL:    "http://dev.lacoste.co.jp/order/return",
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

func TestRefundAuthorizeAndCapture(t *testing.T) {
	createAuth := func() gomerchant.AuthorizeResponse {
		authorizeResponse, _ := Paygent.Authorize(1000, gomerchant.AuthorizeParams{
			Currency: "JPY",
			OrderID:  fmt.Sprint(time.Now().Unix()),
			PaymentMethod: &gomerchant.PaymentMethod{
				CreditCard: &gomerchant.CreditCard{
					Name:     "JCB Card",
					Number:   "3580876521284076",
					ExpMonth: 1,
					ExpYear:  uint(time.Now().Year() + 1),
				},
			},
		})

		return authorizeResponse
	}

	authorizeResponse1 := createAuth()
	fmt.Println(Paygent.Refund(authorizeResponse1.TransactionID, 100, gomerchant.RefundParams{}))
}
