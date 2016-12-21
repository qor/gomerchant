package paygent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/qor/gomerchant"
)

type Paygent struct {
	Config *Config
}

type Config struct {
	Account        string
	Password       string
	MerchantID     string
	ClientFilePath string
	ProductionMode bool
}

func New(config *Config) *Paygent {
	return &Paygent{
		Config: config,
	}
}

func (paygent *Paygent) Client() *http.Client {
	// Load CA cert
	caCertPool := x509.NewCertPool()
	if pemData, err := ioutil.ReadFile(paygent.Config.ClientFilePath); err == nil {
		caCertPool.AppendCertsFromPEM(pemData)
	} else {
		panic(err)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}

func (paygent *Paygent) serviceURLOfTelegramKind(telegramKind string) (string, error) {
	var (
		domain  = TelegramServiceSandboxDomain
		urlPath string
	)

	if paygent.Config.ProductionMode {
		domain = TelegramServiceDomain
	}

	for i := 0; i < len(telegramKind)-1; i++ {
		if p, ok := TelegramServiceURLs[telegramKind[0:len(telegramKind)-i]]; ok {
			urlPath = p
		}
	}
	return path.Join(domain, urlPath), nil
}

func (paygent *Paygent) Request(telegramKind string, params gomerchant.Params) (gomerchant.Params, error) {
	client := paygent.Client()
	if serviceURL, err := paygent.serviceURLOfTelegramKind(telegramKind); err == nil {
		var urlValues url.Values
		for key, value := range params {
			urlValues.Add(key, fmt.Sprint(value))
		}
		client.PostForm(serviceURL, urlValues)
	}
	return gomerchant.Params{}, nil
}

func (*Paygent) Purchase(amount uint64, params *gomerchant.PurchaseParams) (gomerchant.PurchaseResponse, error) {
	return gomerchant.PurchaseResponse{}, nil
}

func (*Paygent) Authorize(amount uint64, params *gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	return gomerchant.AuthorizeResponse{}, nil
}

func (*Paygent) Capture(transactionID string, params *gomerchant.CaptureParams) (gomerchant.CaptureResponse, error) {
	return gomerchant.CaptureResponse{}, nil
}

func (*Paygent) Refund(transactionID string, params *gomerchant.RefundParams) (gomerchant.RefundResponse, error) {
	return gomerchant.RefundResponse{}, nil
}

func (*Paygent) Void(transactionID string, params *gomerchant.VoidParams) (gomerchant.VoidResponse, error) {
	return gomerchant.VoidResponse{}, nil
}
