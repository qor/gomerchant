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
	return Paygent.CreateCreditCard(&gomerchant.CreateCreditCardParams{
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
	authorizeResult, err := Paygent.Authorize(100, &gomerchant.AuthorizeParams{
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

	captureResult, err := Paygent.Capture(authorizeResult.TransactionID, &gomerchant.CaptureParams{})

	if err != nil || captureResult.TransactionID == "" {
		t.Error(err, captureResult)
	}
}

func TestAuthorizeAndCaptureWithSavedCreditCard(t *testing.T) {
	if savedCreditCard, err := createSavedCreditCard(); err == nil {
		authorizeResult, err := Paygent.Authorize(100, &gomerchant.AuthorizeParams{
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

		captureResult, err := Paygent.Capture(authorizeResult.TransactionID, &gomerchant.CaptureParams{})

		if err != nil || captureResult.TransactionID == "" {
			t.Error(err, captureResult)
		}
	}
}
