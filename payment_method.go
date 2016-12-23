package gomerchant

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
