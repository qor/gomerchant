package alipay

import (
	"encoding/json"
	"net/http"

	"github.com/qor/gomerchant"
)

var (
	APIDomain    = "https://openapi.alipay.com/gateway.do"
	DevAPIDomain = "https://openapi.alipaydev.com/gateway.do"
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
	APIDomain      string
}

var _ gomerchant.IntegrationGateway = &Alipay{}

// New initialize alipay
func New(config *Config) *Alipay {
	if config.APIDomain == "" {
		if config.ProductionMode {
			config.APIDomain = APIDomain
		} else {
			config.APIDomain = DevAPIDomain
		}
	}

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

	if params.Params != nil {
		if result, err := json.Marshal(params.Params); err == nil {
			json.Unmarshal(result, &currentParams)
		}
	}

	currentParams.TotalAmount = params.Amount
	currentParams.OutTradeNo = params.OrderID
	currentParams.Subject = params.Description

	checkoutParams := Common{
		Method:     "alipay.trade.page.pay",
		BizContent: currentParams,
	}

	query, err := alipay.Sign(&checkoutParams)

	return alipay.Config.APIDomain + "?" + query, err
}

// VerifyNotification verify notification
func (*Alipay) VerifyNotification(req *http.Request) (bool, error) {
	return false, nil
}

// Refund refund transaction
func (*Alipay) Refund(transactionID string, amount uint, params gomerchant.RefundParams) (gomerchant.RefundResponse, error) {
	type Params struct {
		OutTradeNo   string `json:"out_trade_no,omitempty"`
		TradeNo      string `json:"trade_no,omitempty"`
		RefundAmount uint   `json:"refund_amount"`
		RefundReason string `json:"refund_reason,omitempty"`
		OutRequestNo string `json:"out_request_no,omitempty"`
		OperatorID   string `json:"operator_id,omitempty"`
		StoreID      string `json:"store_id,omitempty"`
		TerminalID   string `json:"terminal_id,omitempty"`
	}

	var currentParams Params

	if params.Params != nil {
		if result, err := json.Marshal(params.Params); err == nil {
			json.Unmarshal(result, &currentParams)
		}
	}

	currentParams.RefundAmount = amount
	currentParams.OutTradeNo = transactionID

	// TODO Do really request & error check

	return gomerchant.RefundResponse{}, nil
}

// Void void transaction
func (*Alipay) Void(transactionID string, params gomerchant.VoidParams) (gomerchant.VoidResponse, error) {
	type Params struct {
		OutTradeNo string `json:"out_trade_no,omitempty"`
		TradeNo    string `json:"trade_no,omitempty"`
		OperatorID string `json:"operator_id,omitempty"`
	}

	var currentParams Params
	if params.Params != nil {
		if result, err := json.Marshal(params.Params); err == nil {
			json.Unmarshal(result, &currentParams)
		}
	}

	currentParams.OutTradeNo = transactionID

	// TODO Do really request & error check

	return gomerchant.VoidResponse{}, nil
}

// Query query transaction
func (*Alipay) Query(transactionID string) (gomerchant.Transaction, error) {
	type Params struct {
		OutTradeNo string `json:"out_trade_no,omitempty"`
		TradeNo    string `json:"trade_no,omitempty"`
	}

	currentParams := Params{OutTradeNo: transactionID}

	// TODO Do really request & error check

	return gomerchant.Transaction{}, nil
}
