package gomerchant

import "regexp"

var Issuers = map[string]*regexp.Regexp{
	`visa`:               regexp.MustCompile(`^4\d{12}(\d{3})?$`),
	`master`:             regexp.MustCompile(`^(5[1-5]\d{4}|677189)\d{10}$`),
	`discover`:           regexp.MustCompile(`^(6011|65\d{2}|64[4-9]\d)\d{12}|(62\d{14})$`),
	`american_express`:   regexp.MustCompile(`^3[47]\d{13}$`),
	`diners_club`:        regexp.MustCompile(`^3(0[0-5]|[68]\d)\d{11}$`),
	`jcb`:                regexp.MustCompile(`^35(27|28|29|[3-8]\d)\d{12}$`),
	`switch`:             regexp.MustCompile(`^6759\d{12}(\d{2,3})?$`),
	`solo`:               regexp.MustCompile(`^6767\d{12}(\d{2,3})?$`),
	`dankort`:            regexp.MustCompile(`^5019\d{12}$`),
	`maestro`:            regexp.MustCompile(`^(5[06-8]|6\d)\d{10,17}$`),
	`forbrugsforeningen`: regexp.MustCompile(`^600722\d{10}$`),
	`laser`:              regexp.MustCompile(`^(6304|6706|6771|6709)\d{8}(\d{4}|\d{6,7})?$`),
}

type PaymentMethod struct {
	Token      string
	CreditCard *CreditCard
}

type CreditCard struct {
	Name     string
	Number   string
	ExpMonth uint
	ExpYear  uint
	CVC      string
}

func (creditCard CreditCard) Issuer() string {
	for name, match := range Issuers {
		if match.MatchString(creditCard.Number) {
			return name
		}
	}
	return ""
}
