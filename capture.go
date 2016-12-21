package gomerchant

type CaptureParams struct {
	Amount uint64
	Params
}

type CaptureResponse struct {
	TransactionID string
	Params
}
