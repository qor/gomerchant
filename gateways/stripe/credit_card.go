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
	c, err := card.Get(creditCardParams.CreditCardID, &stripe.CardParams{Customer: creditCardParams.CustomerID})

	return gomerchant.GetCreditCardResponse{
		CreditCard: &gomerchant.CustomerCreditCard{
			CustomerID:   c.Customer.ID,
			CustomerName: c.Name,
			CreditCardID: c.ID,
			MaskedNumber: c.LastFour,
			ExpMonth:     uint(c.Month),
			ExpYear:      uint(c.Year),
			Brand:        string(c.Brand),
		},
	}, err
}

func (*Stripe) ListCreditCards(listCreditCardsParams gomerchant.ListCreditCardsParams) (gomerchant.ListCreditCardsResponse, error) {
	iter := card.List(&stripe.CardListParams{Customer: listCreditCardsParams.CustomerID})
	resp := gomerchant.ListCreditCardsResponse{}
	for iter.Next() {
		c := iter.Card()
		resp.CreditCards = append(resp.CreditCards, &gomerchant.CustomerCreditCard{
			CustomerID:   c.Customer.ID,
			CustomerName: c.Name,
			CreditCardID: c.ID,
			MaskedNumber: c.LastFour,
			ExpMonth:     uint(c.Month),
			ExpYear:      uint(c.Year),
			Brand:        string(c.Brand),
		})
	}
	return resp, iter.Err()
}

func (*Stripe) DeleteCreditCard(deleteCreditCardParams gomerchant.DeleteCreditCardParams) (gomerchant.DeleteCreditCardResponse, error) {
	_, err := card.Del(deleteCreditCardParams.CreditCardID, &stripe.CardParams{Customer: deleteCreditCardParams.CustomerID})
	return gomerchant.DeleteCreditCardResponse{}, err
}
