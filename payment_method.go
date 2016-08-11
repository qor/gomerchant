package gomerchant

type PaymentMethod struct {
	Token      string
	CreditCard *CreditCard
	// TBD
	// BankAccount
	// Identifier
}

type CreditCard struct {
	Name     string
	Number   string
	ExpMonth int
	ExpYear  int
	CVC      string
}
