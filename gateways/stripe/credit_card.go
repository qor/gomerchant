package stripe

import "github.com/qor/gomerchant"

func (*Stripe) CreateCreditCard(creditCardParams gomerchant.CreateCreditCardParams) (gomerchant.CreditCardResponse, error) {
	return gomerchant.CreditCardResponse{}, nil
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
