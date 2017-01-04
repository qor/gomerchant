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
	creditCardManager gomerchant.CreditCardManager
	gateway           gomerchant.PaymentGateway
	gateway3d         *paygent.Paygent
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

	Paygent := paygent.New(&paygent.Config{
		MerchantID:      config.MerchantID,
		ConnectID:       config.ConnectID,
		ConnectPassword: config.ConnectPassword,
		ClientFilePath:  config.ClientFilePath,
		CertPassword:    config.CertPassword,
		CAFilePath:      config.CAFilePath,
		ProductionMode:  config.ProductionMode,
	})

	creditCardManager = Paygent
	gateway = Paygent
	gateway3d = Paygent
}

func createSavedCreditCard() (gomerchant.CreditCardResponse, error) {
	return creditCardManager.CreateCreditCard(gomerchant.CreateCreditCardParams{
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
	authorizeResult, err := gateway.Authorize(100, gomerchant.AuthorizeParams{
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

	captureResult, err := gateway.Capture(authorizeResult.TransactionID, gomerchant.CaptureParams{})

	if err != nil || captureResult.TransactionID == "" {
		t.Error(err, captureResult)
	}
}

func TestAuthorizeAndCaptureWithSavedCreditCard(t *testing.T) {
	if savedCreditCard, err := createSavedCreditCard(); err == nil {
		authorizeResult, err := gateway.Authorize(100, gomerchant.AuthorizeParams{
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

		captureResult, err := gateway.Capture(authorizeResult.TransactionID, gomerchant.CaptureParams{})

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
		authorizeResult, err := gateway3d.SecureCodeAuthorize(100,
			paygent.SecureCodeParams{
				UserAgent: "User-Agent	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0.2 Safari/602.3.12",
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

func createAuth() gomerchant.AuthorizeResponse {
	authorizeResponse, _ := gateway.Authorize(1000, gomerchant.AuthorizeParams{
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

func TestRefund(t *testing.T) {
	// refund authorized transaction
	authorizeResponse := createAuth()
	if refundResponse, err := gateway.Refund(authorizeResponse.TransactionID, 100, gomerchant.RefundParams{}); err == nil {
		if transaction, err := gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 900 && transaction.Paid == true && transaction.Captured == false && transaction.Cancelled == false && transaction.CreatedAt != nil) {
				t.Errorf("transaction after refund authorized transaction is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}

	// refund authorized transaction, and capture it
	authorizeResponse = createAuth()
	if refundResponse, err := gateway.Refund(authorizeResponse.TransactionID, 150, gomerchant.RefundParams{Captured: true}); err == nil {
		if transaction, err := gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 850 && transaction.Paid == true && transaction.Captured == true && transaction.Cancelled == false && transaction.CreatedAt != nil) {
				t.Errorf("transaction after refund authorized transaction is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}

	// refund captured transaction
	authorizeResponse = createAuth()
	captureResponse, _ := gateway.Capture(authorizeResponse.TransactionID, gomerchant.CaptureParams{})
	if refundResponse, err := gateway.Refund(captureResponse.TransactionID, 200, gomerchant.RefundParams{Captured: true}); err == nil {
		if transaction, err := gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 800 && transaction.Paid == true && transaction.Captured == true && transaction.Cancelled == false && transaction.CreatedAt != nil) {
				t.Errorf("transaction after refund captured transaction is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}
}

func TestVoid(t *testing.T) {
	// void authorized transaction
	authorizeResponse := createAuth()
	if refundResponse, err := gateway.Void(authorizeResponse.TransactionID, gomerchant.VoidParams{}); err == nil {
		if transaction, err := gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 1000 && transaction.Paid == false && transaction.Captured == false && transaction.Cancelled == true && transaction.CreatedAt != nil) {
				t.Errorf("transaction after refund auth is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}

	// void captured transaction
	authorizeResponse = createAuth()
	captureResponse, _ := gateway.Capture(authorizeResponse.TransactionID, gomerchant.CaptureParams{})
	if refundResponse, err := gateway.Void(captureResponse.TransactionID, gomerchant.VoidParams{Captured: true}); err == nil {
		if transaction, err := gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 1000 && transaction.Paid == false && transaction.Captured == false && transaction.Cancelled == true && transaction.CreatedAt != nil) {
				t.Errorf("transaction after refund captured is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}
}

func TestListCreditCards(t *testing.T) {
	if response, err := createSavedCreditCard(); err == nil {
		// create anotther credit card
		creditCardManager.CreateCreditCard(gomerchant.CreateCreditCardParams{
			CustomerID: response.CustomerID,
			CreditCard: &gomerchant.CreditCard{
				Name:     "JCB Card",
				Number:   "3580876521284076",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
			},
		})

		if response, err := creditCardManager.ListCreditCards(gomerchant.ListCreditCardsParams{CustomerID: response.CustomerID}); err == nil {
			if len(response.CreditCards) != 2 {
				t.Errorf("Should found two saved credit cards, but got %v", response.CreditCards)
			}

			for _, creditCard := range response.CreditCards {
				if creditCard.MaskedNumber == "" || creditCard.ExpYear == 0 || creditCard.ExpMonth == 0 || creditCard.CustomerID == "" || creditCard.CreditCardID == "" {
					t.Errorf("Credit card's information should be correct, but got %v", creditCard)
				}
			}
		} else {
			t.Errorf("no error should happen when query saved credit cards, but got %v", err)
		}
	}
}

func TestGetCreditCard(t *testing.T) {
	if response, err := createSavedCreditCard(); err == nil {
		if response, err := creditCardManager.GetCreditCard(gomerchant.GetCreditCardParams{CustomerID: response.CustomerID, CreditCardID: response.CreditCardID}); err == nil {
			creditCard := response.CreditCard
			if creditCard == nil {
				t.Errorf("Should found saved credit cards, but got %v", response)
			} else if creditCard.MaskedNumber == "" || creditCard.ExpYear == 0 || creditCard.ExpMonth == 0 || creditCard.CustomerID == "" || creditCard.CreditCardID == "" {
				t.Errorf("Credit card's information should be correct, but got %v", creditCard)
			}
		} else {
			t.Errorf("no error should happen when query saved credit card, but got %v", err)
		}
	}
}

func TestDeleteCreditCard(t *testing.T) {
	if response, err := createSavedCreditCard(); err == nil {
		if _, err := creditCardManager.DeleteCreditCard(gomerchant.DeleteCreditCardParams{CustomerID: response.CustomerID, CreditCardID: response.CreditCardID}); err == nil {
			if response, err := creditCardManager.GetCreditCard(gomerchant.GetCreditCardParams{CustomerID: response.CustomerID, CreditCardID: response.CreditCardID}); err == nil {
				t.Errorf("Should failed to get credit card, but got %v", response)
			}
		} else {
			t.Errorf("no error should happen when delete saved credit card, but got %v", err)
		}
	}
}
