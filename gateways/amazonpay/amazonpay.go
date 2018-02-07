package amazonpay

import "github.com/qor/gomerchant"

var _ gomerchant.PaymentGateway = &AmazonPay{}

// AmazonPay amazon pay
type AmazonPay struct {
	*Config
}

// Config amazon pay configuration
type Config struct {
	MerchantID     string
	AccessKey      string
	SecretKey      string
	ClientID       string
	ClientSecret   string
	Region         string
	CurrencyCode   string
	ProductionMode bool
}

// New initialize amazon pay
func New(config *Config) *AmazonPay {
	if config == nil {
		config = &Config{}
	}

	return &AmazonPay{Config: config}
}

func (AmazonPay) Authorize(amount uint64, params gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	return gomerchant.AuthorizeResponse{}, nil
}

func (AmazonPay) CompleteAuthorize(paymentID string, params gomerchant.CompleteAuthorizeParams) (gomerchant.CompleteAuthorizeResponse, error) {
	return gomerchant.CompleteAuthorizeResponse{}, nil
}

func (AmazonPay) Capture(transactionID string, params gomerchant.CaptureParams) (gomerchant.CaptureResponse, error) {
	return gomerchant.CaptureResponse{}, nil
}

func (AmazonPay) Refund(transactionID string, amount uint, params gomerchant.RefundParams) (gomerchant.RefundResponse, error) {
	return gomerchant.RefundResponse{}, nil
}

func (AmazonPay) Void(transactionID string, params gomerchant.VoidParams) (gomerchant.VoidResponse, error) {
	return gomerchant.VoidResponse{}, nil
}

func (AmazonPay) Query(transactionID string) (gomerchant.Transaction, error) {
	return gomerchant.Transaction{}, nil
}
