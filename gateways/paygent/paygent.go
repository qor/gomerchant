package paygent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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
	CAFilePath     string
	ProductionMode bool
}

func New(config *Config) *Paygent {
	return &Paygent{
		Config: config,
	}
}

func (paygent *Paygent) Client() (*http.Client, error) {
	var (
		certBytes, certKeyBytes []byte
		caCertPool              = x509.NewCertPool()
	)

	if pemData, err := ioutil.ReadFile(paygent.Config.ClientFilePath); err == nil {
		var (
			originalPemData []byte
			block           *pem.Block
		)

		caCertPool.AppendCertsFromPEM(pemData)

		for len(pemData) > 0 {
			originalPemData = pemData
			if block, pemData = pem.Decode(pemData); block == nil {
				break
			}

			if block.Type == "CERTIFICATE" && len(certBytes) == 0 {
				if len(certBytes) == 0 {
					certBytes = originalPemData[0 : len(originalPemData)-len(pemData)-1]
				}
			}

			if block.Type == "RSA PRIVATE KEY" {
				certKeyBytes = originalPemData[0 : len(originalPemData)-len(pemData)-1]
			}
		}
	}

	if pemData, err := ioutil.ReadFile(paygent.Config.CAFilePath); err == nil {
		caCertPool.AppendCertsFromPEM(pemData)
	}

	if cert, err := tls.X509KeyPair(certBytes, certKeyBytes); err == nil {
		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		return &http.Client{Transport: transport}, nil
	} else {
		return nil, err
	}
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

	u, err := url.Parse(domain)
	u.Path = urlPath
	return u.String(), err
}

func (paygent *Paygent) Request(telegramKind string, params gomerchant.Params) (gomerchant.Params, error) {
	if client, err := paygent.Client(); err == nil {
		if serviceURL, err := paygent.serviceURLOfTelegramKind(telegramKind); err == nil {
			var urlValues url.Values
			for key, value := range params {
				urlValues.Add(key, fmt.Sprint(value))
			}
			response, err := client.PostForm(serviceURL, urlValues)
			if err == nil && response.StatusCode == 200 {
				var bodyBytes []byte
				response.Body.Read(bodyBytes)
				fmt.Println(string(bodyBytes))
			}
			return gomerchant.Params{}, err
		} else {
			return gomerchant.Params{}, err
		}
	} else {
		return gomerchant.Params{}, err
	}
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