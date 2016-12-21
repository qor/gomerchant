package paygent

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

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

func (paygent *Paygent) Request(telegramKind string, params gomerchant.Extra) gomerchant.Extra {
	client := paygent.Client()
	client.PostForm()
	return gomerchant.Extra{}
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
