package alipay

import (
	"net/http"

	"github.com/qor/gomerchant"
)

// Alipay alipay struct
type Alipay struct {
	Config *Config
}

// Config alipay config
type Config struct {
	APPID          string `required:"true"`
	PrivateKey     string `required:"true"`
	PublicKey      string `required:"true"`
	ProductionMode bool
}

var _ gomerchant.IntegrationGateway = &Alipay{}

// New initialize alipay
func New(config *Config) *Alipay {
	return &Alipay{
		Config: config,
	}
}

// CheckoutURL generate CheckoutURL for alipay
func (*Alipay) CheckoutURL(params gomerchant.CheckoutParams) (string, error) {
	return "", nil
}

// VerifyNotification verify notification
func (*Alipay) VerifyNotification(req *http.Request) (bool, error) {
	return false, nil
}

// Refund refund transaction
func (*Alipay) Refund(transactionID string, amount uint, params gomerchant.RefundParams) (gomerchant.RefundResponse, error) {
	return gomerchant.RefundResponse{}, nil
}

// Void void transaction
func (*Alipay) Void(transactionID string, params gomerchant.VoidParams) (gomerchant.VoidResponse, error) {
	return gomerchant.VoidResponse{}, nil
}

// Query query transaction
func (*Alipay) Query(transactionID string) (gomerchant.Transaction, error) {
	return gomerchant.Transaction{}, nil
}
