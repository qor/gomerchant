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
	AppID          string `required:"true"`
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
func (alipay *Alipay) CheckoutURL(params gomerchant.CheckoutParams) (string, error) {
	type Params struct {
		OutTradeNo         string `json:"out_trade_no"`
		ProductCode        string `json:"product_code"`
		TotalAmount        uint64 `json:"total_amount"`
		Subject            string `json:"subject"`
		Body               string `json:"body,omitempty"`
		GoodsDetail        string `json:"goods_detail,omitempty"`
		PassbackParams     string `json:"passback_params,omitempty"`
		ExtendParams       string `json:"extend_params,omitempty"`
		GoodsType          string `json:"goods_type,omitempty"`
		TimeoutExpress     string `json:"timeout_express,omitempty"`
		EnablePayChannels  string `json:"enable_pay_channels,omitempty"`
		DisablePayChannels string `json:"disable_pay_channels,omitempty"`
		AuthToken          string `json:"auth_token,omitempty"`
		QRPayMode          string `json:"qr_pay_mode,omitempty"`
		QRWidth            string `json:"qrcode_width,omitempty"`
	}

	var currentParams Params

	currentParams.TotalAmount = params.Amount
	currentParams.OutTradeNo = params.OrderID
	currentParams.Subject = params.Description

	checkoutParams := Common{
		Method:     "alipay.trade.page.pay",
		BizContent: currentParams,
	}

	query, err := alipay.Sign(&checkoutParams)

	return query, err
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
