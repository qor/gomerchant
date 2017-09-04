package stripe

import (
	"fmt"

	"github.com/qor/gomerchant"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
)

func (*Stripe) CreateCreditCard(creditCardParams gomerchant.CreateCreditCardParams) (gomerchant.CreditCardResponse, error) {
	c, err := card.New(&stripe.CardParams{
		Customer: creditCardParams.CustomerID,
		Name:     creditCardParams.CreditCard.Name,
		Number:   creditCardParams.CreditCard.Number,
		Month:    fmt.Sprint(creditCardParams.CreditCard.ExpMonth),
		Year:     fmt.Sprint(creditCardParams.CreditCard.ExpYear),
		CVC:      creditCardParams.CreditCard.CVC,
	})

	return gomerchant.CreditCardResponse{CustomerID: c.Customer.ID, CreditCardID: c.ID}, err
}

func (*Stripe) GetCreditCard(creditCardParams gomerchant.GetCreditCardParams) (gomerchant.GetCreditCardResponse, error) {
	return gomerchant.GetCreditCardResponse{}, nil
}

func (*Stripe) ListCreditCards(listCreditCardsParams gomerchant.ListCreditCardsParams) (gomerchant.ListCreditCardsResponse, error) {
	return gomerchant.ListCreditCardsResponse{}, nil
}

func (*Stripe) DeleteCreditCard(deleteCreditCardParams gomerchant.DeleteCreditCardParams) (gomerchant.DeleteCreditCardResponse, error) {
	return gomerchant.DeleteCreditCardResponse{}, nil
}
