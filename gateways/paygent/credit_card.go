package paygent

import (
	"fmt"

	"github.com/qor/gomerchant"
)

func getValidTerm(creditCard *gomerchant.CreditCard) string {
	return fmt.Sprintf("%02d", creditCard.ExpMonth) + fmt.Sprint(creditCard.ExpYear)[len(fmt.Sprint(creditCard.ExpYear))-2:]
}

var issuersMap = map[string]string{"visa": "V", "master": "M", "american_express": "X", "diners_club": "C", "jcb": "J"}

func (paygent *Paygent) CreateCreditCard(creditCardParams gomerchant.CreateCreditCardParams) (gomerchant.CreditCardResponse, error) {
	var (
		response   = gomerchant.CreditCardResponse{CustomerID: creditCardParams.CustomerID}
		creditCard = creditCardParams.CreditCard
		issuer, _  = issuersMap[creditCard.Issuer()]
	)

	results, err := paygent.Request("025", gomerchant.Params{
		"customer_id":     creditCardParams.CustomerID,
		"card_number":     creditCard.Number,
		"card_valid_term": getValidTerm(creditCard),
		"cardholder_name": creditCard.Name,
		"card_brand":      issuer,
	}.IgnoreBlankFields())

	if err == nil {
		if customerCardID, ok := results.Get("customer_card_id"); ok {
			response.CreditCardID = fmt.Sprint(customerCardID)
		}
	}
	response.Params = results.Params

	return response, err
}

func (paygent *Paygent) DeleteCreditCard(deleteCreditCardParams gomerchant.DeleteCreditCardParams) (gomerchant.DeleteCreditCardResponse, error) {
	var (
		response = gomerchant.DeleteCreditCardResponse{CustomerID: deleteCreditCardParams.CustomerID}
	)

	results, err := paygent.Request("026", gomerchant.Params{})
	return response, err
}

func (paygent *Paygent) ListCreditCards(listCreditCardsParams gomerchant.ListCreditCardsParams) (gomerchant.ListCreditCardsResponse, error) {
	var (
		response = gomerchant.ListCreditCardsResponse{CustomerID: listCreditCardsParams.CustomerID}
	)

	results, err := paygent.Request("027", gomerchant.Params{})
	return response, err
}
