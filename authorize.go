package gomerchant

type AuthorizeParams struct {
	Amount          uint64
	Currency        string
	Customer        string
	Description     string
	OrderID         string
	BillingAddress  *Address
	ShippingAddress *Address
	Extra
}

type AuthorizeResponse struct {
	TransactionID string
	Extra
}
