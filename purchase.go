package gomerchant

type PurchaseParams struct {
	Amount          uint64
	Currency        string
	Customer        string
	Description     string
	OrderID         string
	BillingAddress  *Address
	ShippingAddress *Address
	PaymentMethod   *PaymentMethod
	Params
}

type PurchaseResponse struct {
	TransactionID string
	Params
}
