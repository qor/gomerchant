package paygent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

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
		cookieJar, _ := cookiejar.New(nil)

		return &http.Client{
			Transport: transport,
			Jar:       cookieJar,
		}, nil
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
			break
		}
	}

	u, err := url.Parse(domain)
	u.Path = urlPath
	return u, err
}

var ResponseParser = regexp.MustCompile(`(?m)(\w+?)=(.*?)\r\n`)

func (paygent *Paygent) Request(telegramKind string, params gomerchant.Params) (gomerchant.Params, error) {
	var (
		serviceURL  *url.URL
		urlValues   = url.Values{}
		results     = gomerchant.Params{}
		client, err = paygent.Client()
	)

	if err == nil {
		serviceURL, err = paygent.serviceURLOfTelegramKind(telegramKind)

		if err == nil {
			urlValues.Add("merchant_id", paygent.Config.MerchantID)
			urlValues.Add("connect_id", paygent.Config.ConnectID)
			urlValues.Add("connect_password", paygent.Config.ConnectPassword)
			urlValues.Add("telegram_version", paygent.Config.TelegramVersion)
			urlValues.Add("telegram_kind", telegramKind)

			for key, value := range params {
				urlValues.Add(key, fmt.Sprint(value))
			}

			response, err := client.Post(serviceURL.String(), "application/x-www-form-urlencoded", strings.NewReader(urlValues.Encode()))

			if err == nil {
				if response.StatusCode == 200 {
					defer response.Body.Close()
					var bodyBytes []byte
					bodyBytes, err = ioutil.ReadAll(response.Body)
					if err == nil {
						for _, value := range ResponseParser.FindAllStringSubmatch(string(bodyBytes), -1) {
							results.Set(value[1], value[2])
						}
						return results, nil
					}
				}
				err = fmt.Errorf("status code: %v", response.StatusCode)
			}
		}
	}

	return results, err
}

func (*Paygent) Purchase(amount uint64, params *gomerchant.PurchaseParams) (gomerchant.PurchaseResponse, error) {
	return gomerchant.PurchaseResponse{}, nil
}

func (paygent *Paygent) Authorize(amount uint64, params *gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	requestParams := gomerchant.Params{
		"payment_id":     params.OrderID,
		"payment_amount": amount,
		"payment_class":  10,
	}

	if paymentMethod := params.PaymentMethod; paymentMethod != nil {
		if creditCard := paymentMethod.CreditCard; creditCard != nil {
			requestParams["card_number"] = creditCard.Number
			requestParams["card_valid_term"] = getValidTerm(creditCard)
			requestParams["3dsecure_ryaku"] = 1
		}
	}

	results, err := paygent.Request("020", requestParams)
	fmt.Println(results)
	return gomerchant.AuthorizeResponse{}, err
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
