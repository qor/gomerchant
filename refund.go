package gomerchant

type RefundParams struct {
	Params
}

type RefundResponse struct {
	TransactionID string
	Params
}
