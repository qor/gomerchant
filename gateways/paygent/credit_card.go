package paygent

import (
	"encoding/csv"
	"fmt"
	"strings"

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
	var response = gomerchant.DeleteCreditCardResponse{}

	results, err := paygent.Request("026", gomerchant.Params{"customer_id": deleteCreditCardParams.CustomerID, "customer_card_id": deleteCreditCardParams.CreditCardID}.IgnoreBlankFields())
	response.Params = response.Params
	return response, err
}

func (paygent *Paygent) ListCreditCards(listCreditCardsParams gomerchant.ListCreditCardsParams) (gomerchant.ListCreditCardsResponse, error) {
	var response = gomerchant.ListCreditCardsResponse{CustomerID: listCreditCardsParams.CustomerID}

	results, err := paygent.Request("027", gomerchant.Params{"customer_id": listCreditCardsParams.CustomerID})

	if err == nil {
	}

	return response, err
}

func parseListCreditCardsResponse(response *Response) ([]gomerchant.CreditCard, error) {
	for _, str := range strings.Split(response.RawBody, "\r\n") {
		row := csv.NewReader(strings.NewReader(str))
		if record, err := row.Read(); err == nil {
			switch record[0] {
			case "1":
				// response information
				response.Result = record[1]
				response.ResponseCode = record[2]
				response.ResponseDetail = record[3]
			case "2":
				// card header
			case "3":
				// card information
			case "4":
				// card numbers
			}
		} else {
			return response, err
		}
	}
}
