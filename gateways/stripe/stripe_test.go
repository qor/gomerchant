package stripe_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/configor"
	"github.com/qor/gomerchant"
	"github.com/qor/gomerchant/gateways/stripe"
)

var (
	creditCardManager gomerchant.CreditCardManager
	gateway           gomerchant.PaymentGateway
)

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

	Stripe := stripe.New(&stripe.Config{
		Key: config.Key,
	})

	// creditCardManager = Stripe
	gateway = Stripe
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
