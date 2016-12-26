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

func TestCreateCreditCard(t *testing.T) {
	result, err := Paygent.CreateCreditCard(&gomerchant.CreditCard{
		Name:     "JCB Card",
		Number:   "3580876521284076",
		ExpMonth: 1,
		ExpYear:  uint(time.Now().Year() + 1),
		CVC:      "111",
	})

	fmt.Println(result)
	fmt.Println(err)
}
func TestAuthorize(t *testing.T) {
	Paygent.Authorize(100, &gomerchant.AuthorizeParams{
		Currency: "JPY",
		OrderID:  fmt.Sprint(time.Now().Unix()),
		PaymentMethod: &gomerchant.PaymentMethod{
			CreditCard: &gomerchant.CreditCard{
				Name:     "JCB Card",
				Number:   "3580876521284076",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
				CVC:      "111",
			},
		},
	})
}
