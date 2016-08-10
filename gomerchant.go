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

type Options struct {
	OrderID         string
	Currency        string
	Email           string
	IP              string
	Customer        string
	Invoice         string
	Merchant        string
	Description     string
	BillingAddress  *Address
	ShippingAddress *Address

	Extra interface{}
}

type Address struct {
	Name     string
	Company  string
	Address1 string
	Address2 string
	City     string
	State    string
	Country  string
	ZIP      string
	Phone    string
}

type Response struct {
	ID    string
	Extra interface{}
}

type Payer interface {
	Purchase(amount uint64, pm *PaymentMethod, opts *Options) (Response, error)
	Authorize(amount uint64, pm *PaymentMethod, opts *Options) (Response, error)
	Capture(amount uint64, id string, opts *Options) (Response, error)
	Void(id string, opts *Options) (Response, error)
	Store(pm *PaymentMethod, opts *Options) (Response, error)
}

// TBD
type CCManager interface {
	Store(pm *PaymentMethod, opts *Options) error
	Unstore(pm *PaymentMethod, opts *Options) error
	Update(pm *PaymentMethod, opts *Options) error
}
