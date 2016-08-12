package gomerchant

type CaptureParams struct {
	Amount uint64
	Extra
}

type CaptureResponse struct {
	TransactionID string
	Extra
}
