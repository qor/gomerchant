package gomerchant

import (
	"log"
	"os"
)

type Gomerchant struct {
	PaymentGateway PaymentGateway
	Config         *Config
}

type Config struct {
	Logger Logger
}

type Logger interface {
	Print(values ...interface{})
}

func New(paymentGateway PaymentGateway, config *Config) *Gomerchant {
	if config.Logger == nil {
		config.Logger = log.New(os.Stdout, "\r\n", 0)
	}

	return &Gomerchant{PaymentGateway: paymentGateway, Config: config}
}

// TBD: logger, common error validations...

func (gomerchant *Gomerchant) Purchase(amount uint64, pm *PaymentMethod, params *PurchaseParams) (PurchaseResponse, error) {
	response, err := gomerchant.Purchase(amount, pm, params)
	gomerchant.Config.Logger.Print("Purchase", pm, params, response, err)
	return response, err
}

func (gomerchant *Gomerchant) Authorize(amount uint64, pm *PaymentMethod, params *AuthorizeParams) (AuthorizeResponse, error) {
	response, err := gomerchant.Authorize(amount, pm, params)
	gomerchant.Config.Logger.Print("Authorize", pm, params, response, err)
	return response, err
}

func (gomerchant *Gomerchant) Capture(transactionID string, params *CaptureParams) (CaptureResponse, error) {
	response, err := gomerchant.Capture(transactionID, params)
	gomerchant.Config.Logger.Print("Capture", transactionID, response, err)
	return response, err
}

func (gomerchant *Gomerchant) Refund(transactionID string, params *RefundParams) (RefundResponse, error) {
	response, err := gomerchant.Refund(transactionID, params)
	gomerchant.Config.Logger.Print("Refund", transactionID, params, response, err)
	return response, err
}

func (gomerchant *Gomerchant) Void(transactionID string, params *VoidParams) (VoidResponse, error) {
	response, err := gomerchant.Void(transactionID, params)
	gomerchant.Config.Logger.Print("Void", transactionID, params, response, err)
	return response, err
}
