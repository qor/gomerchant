package gomerchant

type RefundParams struct {
	Extra
}

type RefundResponse struct {
	TransactionID string
	Extra
}
