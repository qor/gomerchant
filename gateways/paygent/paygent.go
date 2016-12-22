package paygent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
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
	MerchantID      string
	ConnectID       string
	ConnectPassword string
	TelegramVersion string

	ClientFilePath string
	CertPassword   string
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
				if x509.IsEncryptedPEMBlock(block) {
					if results, err := x509.DecryptPEMBlock(block, []byte(paygent.Config.CertPassword)); err == nil {
						certKeyBytes = []byte("-----BEGIN RSA PRIVATE KEY-----\n" + base64.StdEncoding.EncodeToString(results) + "\n-----END RSA PRIVATE KEY-----")
					} else {
						return nil, err
					}
				} else {
					certKeyBytes = originalPemData[0 : len(originalPemData)-len(pemData)-1]
				}
			}
		}
	}

	if pemData, err := ioutil.ReadFile(paygent.Config.CAFilePath); err == nil {
		caCertPool.AppendCertsFromPEM(pemData)
	}

	if cert, err := tls.X509KeyPair(certBytes, certKeyBytes); err == nil {
		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates:  []tls.Certificate{cert},
			RootCAs:       caCertPool,
			Renegotiation: tls.RenegotiateFreelyAsClient,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		return &http.Client{Transport: transport}, nil
	} else {
		return nil, err
	}
}

func (paygent *Paygent) serviceURLOfTelegramKind(telegramKind string) (*url.URL, error) {
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
	return u, err
}

func (paygent *Paygent) Request(telegramKind string, params gomerchant.Params) (gomerchant.Params, error) {
	if client, err := paygent.Client(); err == nil {
		if serviceURL, err := paygent.serviceURLOfTelegramKind(telegramKind); err == nil {
			urlValues := url.Values{}
			urlValues.Add("merchant_id", paygent.Config.MerchantID)
			urlValues.Add("connect_id", paygent.Config.ConnectID)
			urlValues.Add("connect_password", paygent.Config.ConnectPassword)
			urlValues.Add("telegram_version", paygent.Config.TelegramVersion)
			urlValues.Add("telegram_kind", telegramKind)

			for key, value := range params {
				urlValues.Add(key, fmt.Sprint(value))
			}

			serviceURL.RawQuery = urlValues.Encode()
			response, err := client.PostForm(serviceURL.String(), url.Values{})

			if err == nil {
				if response.StatusCode == 200 {
					var bodyBytes []byte
					response.Body.Read(bodyBytes)
					fmt.Println(string(bodyBytes))
				} else {
					err = fmt.Errorf("status code: %v", response.StatusCode)
				}
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
