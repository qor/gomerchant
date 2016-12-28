package paygent

import (
	"fmt"

	"github.com/qor/gomerchant"
)

func getValidTerm(creditCard *gomerchant.CreditCard) string {
	return fmt.Sprintf("%02d", creditCard.ExpMonth) + fmt.Sprint(creditCard.ExpYear)[len(fmt.Sprint(creditCard.ExpYear))-2:]
}

var issuersMap = map[string]string{"visa": "V", "master": "M", "american_express": "X", "diners_club": "C", "jcb": "J"}

func (paygent *Paygent) CreateCreditCard(creditCardParams *gomerchant.CreateCreditCardParams) (gomerchant.Params, error) {
	var (
		creditCard = creditCardParams.CreditCard
		issuer, _  = issuersMap[creditCard.Issuer()]
	)

	return paygent.Request("025", gomerchant.Params{
		"customer_id":     creditCardParams.CustomerID,
		"card_number":     creditCard.Number,
		"card_valid_term": getValidTerm(creditCard),
		"cardholder_name": creditCard.Name,
		"card_brand":      issuer,
	}.IgnoreBlankFields())
}
