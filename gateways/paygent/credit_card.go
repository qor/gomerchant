package paygent

import (
	"fmt"

	"github.com/qor/gomerchant"
)

func getValidTerm(creditCard *gomerchant.CreditCard) string {
	return fmt.Sprintf("%02d", creditCard.ExpMonth) + fmt.Sprint(creditCard.ExpYear)[len(fmt.Sprint(creditCard.ExpYear))-2:]
}

func (paygent *Paygent) CreateCreditCard(creditCard *gomerchant.CreditCard) (gomerchant.Params, error) {
	return paygent.Request("025", gomerchant.Params{
		"card_number":     creditCard.Number,
		"card_valid_term": getValidTerm(creditCard),
		"cardholder_name": creditCard.Name,
		"card_token":      creditCard.CVC,
	}.IgnoreBlankFields())
}
