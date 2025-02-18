package paygent

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/qor/gomerchant"
	"github.com/youmark/pkcs8"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Paygent struct {
	Config *Config
}

type Config struct {
	MerchantID           string `required:"true"`
	ConnectID            string `required:"true"`
	ConnectPassword      string `required:"true"`
	MerchantName         string
	TelegramVersion      string
	ThreeDSAcceptanceKey string // 3D Secure result acceptance hash key

	CertPassword      string
	ClientFilePath    string // this is required, if ClientFileContent is blank
	ClientFileContent string // this is required, if ClientFilePath is blank
	CAFilePath        string // this is required, if CAFileContent is blank
	CAFileContent     string // this is required, if CAFilePath is blank

	ProductionMode  bool
	SecurityCodeUse bool
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

	if paygent.Config.ClientFileContent == "" {
		if pemData, err := ioutil.ReadFile(paygent.Config.ClientFilePath); err == nil {
			paygent.Config.ClientFileContent = string(pemData)
		} else {
			return nil, err
		}
	}

	if paygent.Config.CAFileContent == "" {
		if pemData, err := ioutil.ReadFile(paygent.Config.CAFilePath); err == nil {
			paygent.Config.CAFileContent = string(pemData)
		} else {
			return nil, err
		}
	}

	{ // ClientFile
		var (
			originalPemData []byte
			block           *pem.Block
			pemData         = []byte(paygent.Config.ClientFileContent)
		)

		caCertPool.AppendCertsFromPEM(pemData)

		for len(pemData) > 0 {
			originalPemData = pemData
			if block, pemData = pem.Decode(pemData); block == nil {
				break
			}

			switch block.Type {
			case "ENCRYPTED PRIVATE KEY":
				{
					if privateKey, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, []byte(paygent.Config.CertPassword)); err == nil {
						privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
						if err != nil {
							return nil, err
						}

						certKeyBytes = []byte("-----BEGIN RSA PRIVATE KEY-----\n" + base64.StdEncoding.EncodeToString(privateKeyBytes) + "\n-----END RSA PRIVATE KEY-----")
					} else {
						return nil, err
					}
				}
			case "RSA PRIVATE KEY":
				{
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
			case "CERTIFICATE":
				{
					if len(certBytes) == 0 {
						certBytes = originalPemData[0 : len(originalPemData)-len(pemData)-1]
					}
				}
			}
		}
	}

	{ // CAFile
		caCertPool.AppendCertsFromPEM([]byte(paygent.Config.CAFileContent))
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

var ResponseParser = regexp.MustCompile(`(?s)(\w+?)=(<!DOCTYPE.*HTML>|.*?)(\r\n|$)`)

type Response struct {
	RawBody        string
	Result         string
	ResponseCode   string
	ResponseDetail string
	gomerchant.Params
}

func (paygent *Paygent) Request(telegramKind string, params gomerchant.Params) (Response, error) {
	var (
		response    *http.Response
		serviceURL  *url.URL
		urlValues   = url.Values{}
		results     = Response{Params: gomerchant.Params{}}
		client, err = paygent.Client()
	)
	if err == nil {
		serviceURL, err = paygent.serviceURLOfTelegramKind(telegramKind)

		if err == nil {
			urlValues.Add("merchant_id", paygent.Config.MerchantID)
			urlValues.Add("connect_id", paygent.Config.ConnectID)
			urlValues.Add("connect_password", paygent.Config.ConnectPassword)
			if paygent.Config.TelegramVersion != "" {
				urlValues.Add("telegram_version", paygent.Config.TelegramVersion)
			} else {
				urlValues.Add("telegram_version", "1.0")
			}
			urlValues.Add("telegram_kind", telegramKind)

			for key, value := range params {
				urlValues.Add(key, fmt.Sprint(value))
			}

			response, err = client.Post(serviceURL.String(), "application/x-www-form-urlencoded", strings.NewReader(urlValues.Encode()))
			if err == nil {
				if response.StatusCode == 200 {
					defer response.Body.Close()
					var bodyBytes []byte
					bodyBytes, err = io.ReadAll(response.Body)
					utf8Bytes := bodyBytes
					contentType := response.Header.Get("Content-Type")
					if !strings.Contains(contentType, "charset=UTF-8") {
						shiftJISToUTF8 := transform.NewReader(bytes.NewReader(bodyBytes), japanese.ShiftJIS.NewDecoder())
						utf8Bytes, _ = io.ReadAll(shiftJISToUTF8)
					}
					results.RawBody = string(utf8Bytes)
					results.Params.Set("RawBody", results.RawBody)

					if err == nil {
						for _, value := range ResponseParser.FindAllStringSubmatch(string(utf8Bytes), -1) {
							if value[1] == "result" {
								results.Result = value[2]
							}

							if value[1] == "response_code" {
								results.ResponseCode = value[2]
							}

							if value[1] == "response_detail" {
								results.ResponseDetail = value[2]
							}

							results.Set(value[1], value[2])
						}

						if results.Result == "1" {
							if results.ResponseDetail != "" {
								err = errors.New(results.ResponseDetail)
							} else {
								err = errors.New("failed to process this request")
							}
						}
						return results, err
					}
				}
				err = fmt.Errorf("status code: %v", response.StatusCode)
			}
		}
	}

	return results, err
}

func (paygent *Paygent) Authorize(amount uint64, params gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	var (
		response      gomerchant.AuthorizeResponse
		requestParams = gomerchant.Params{
			"trading_id":     params.OrderID,
			"payment_amount": amount,
			"payment_class":  10,
		}
	)

	if ok, threeDomainParams := get3DModeParams(params); ok {
		requestParams["http_user_agent"] = threeDomainParams.UserAgent
		requestParams["term_url"] = threeDomainParams.TermURL
		requestParams["http_accept"] = threeDomainParams.HttpAccept
	} else {
		requestParams["3dsecure_ryaku"] = 1
	}

	if paymentMethod := params.PaymentMethod; paymentMethod != nil {
		if savedCreditCard := paymentMethod.SavedCreditCard; savedCreditCard != nil {
			requestParams["stock_card_mode"] = 1
			requestParams["customer_id"] = savedCreditCard.CustomerID
			requestParams["customer_card_id"] = savedCreditCard.CreditCardID
			if savedCreditCard.ThreeDSAuthID != "" {
				requestParams["3ds_auth_id"] = savedCreditCard.ThreeDSAuthID
				requestParams["3dsecure_use_type"] = "2" // 3D Secure 2.0
			} else {
				if paygent.Config.SecurityCodeUse {
					requestParams["security_code_use"] = 1
				}
				requestParams["card_conf_number"] = savedCreditCard.CVC
			}

		} else if creditCard := paymentMethod.CreditCard; creditCard != nil {
			if paygent.Config.SecurityCodeUse {
				requestParams["security_code_use"] = 1
			}
			requestParams["card_number"] = creditCard.Number
			requestParams["card_valid_term"] = getValidTerm(creditCard)
			requestParams["card_conf_number"] = creditCard.CVC
			if creditCard.ThreeDSAuthID != "" {
				requestParams["3ds_auth_id"] = creditCard.ThreeDSAuthID
				requestParams["3dsecure_use_type"] = "2" // 3D Secure 2.0
			}

		} else {
			return response, gomerchant.ErrNotSupportedPaymentMethod
		}
	} else {
		return response, gomerchant.ErrNotSupportedPaymentMethod
	}

	results, err := paygent.Request("020", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			response.TransactionID = fmt.Sprint(paymentID)
		}

		// If 3D Mode
		if ok, _ := get3DModeParams(params); ok {
			if result, ok := results.Get("out_acs_html"); ok && fmt.Sprint(result) != "" {
				response.HandleRequest = true
				response.RequestHandler = func(writer http.ResponseWriter, request *http.Request, _ gomerchant.Params) error {
					_, e := io.WriteString(writer, fmt.Sprint(result))
					return e
				}
			}
		}
	}

	response.Params = results.Params

	return response, err
}

func (paygent *Paygent) Capture(transactionID string, params gomerchant.CaptureParams) (gomerchant.CaptureResponse, error) {
	var (
		response      gomerchant.CaptureResponse
		requestParams = gomerchant.Params{"payment_id": transactionID}
	)

	results, err := paygent.Request("022", requestParams)

	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			response.TransactionID = fmt.Sprint(paymentID)
		}
	}
	response.Params = results.Params

	return response, err
}

func (paygent *Paygent) Refund(transactionID string, amount uint, params gomerchant.RefundParams) (response gomerchant.RefundResponse, err error) {
	var (
		results       Response
		requestParams = gomerchant.Params{
			"payment_id":     transactionID,
			"payment_amount": amount,
			"reduction_flag": 1,
		}
	)

	if params.Captured {
		results, err = paygent.Request("029", requestParams)
	} else {
		results, err = paygent.Request("028", requestParams)
	}

	response.Params = results.Params
	if paymentID, ok := getPaymentID(results); ok {
		response.TransactionID = paymentID
	}

	return response, err
}

func (paygent *Paygent) Void(transactionID string, params gomerchant.VoidParams) (response gomerchant.VoidResponse, err error) {
	var results Response

	if params.Captured {
		results, err = paygent.Request("023", gomerchant.Params{"payment_id": transactionID})
	} else {
		results, err = paygent.Request("021", gomerchant.Params{"payment_id": transactionID})
	}

	response.Params = results.Params
	if paymentID, ok := getPaymentID(results); ok {
		response.TransactionID = paymentID
	}

	return response, err
}

func (paygent *Paygent) Query(transactionID string) (gomerchant.Transaction, error) {
	results, err := paygent.Request("094", gomerchant.Params{"payment_id": transactionID})
	transaction := extractTransactionFromPaygentResponse(results)
	transaction.Params = results.Params
	return transaction, err
}

func (paygent *Paygent) InquiryNotification(noticeID string) (response gomerchant.InquiryResponse, err error) {
	reqParam := gomerchant.Params{}
	if len(noticeID) != 0 {
		reqParam["payment_notice_id"] = noticeID
	}
	results, err := paygent.Request("091", reqParam)
	response.Params = results.Params
	if paymentID, ok := getPaymentID(results); ok {
		response.TransactionID = paymentID
		if tradingID, ok := results.Get("trading_id"); ok {
			response.TradingID = fmt.Sprint(tradingID)
		}

		if paymentNoticeID, ok := results.Get("payment_notice_id"); ok {
			response.PaymentNoticeID = fmt.Sprint(paymentNoticeID)
		}

		if paymentInitDate, ok := results.Get("payment_init_date"); ok {
			response.PaymentInitDate = fmt.Sprint(paymentInitDate)
		}

		if paymentChangeDate, ok := results.Get("change_date"); ok {
			response.PaymentChangeDate = fmt.Sprint(paymentChangeDate)
		}

		if paymentAmount, ok := results.Get("payment_amount"); ok {
			response.PaymentAmount = fmt.Sprint(paymentAmount)
		}

		if basePaymentID, ok := results.Get("base_payment_id"); ok {
			response.BasePaymentID = fmt.Sprint(basePaymentID)
		}

		if paymentStatus, ok := results.Get("payment_status"); ok {
			response.PaymentStatus = fmt.Sprint(paymentStatus)
		}
	}

	if successCode, ok := results.Get("success_code"); ok {
		response.SuccessCode = fmt.Sprint(successCode)
	}

	if successDetail, ok := results.Get("success_detail"); ok {
		response.SuccessDetail = fmt.Sprint(successDetail)
	}
	return response, err
}

// This is rakuten pay authorize function
// Before user confirmed on rakuten page status is 10:already applid
// After user confirmed status change to 20: Authorization OK
func (paygent *Paygent) RakutePayApplicationMessage(amount uint64, params gomerchant.ApplicationParams) (gomerchant.ApplicationResponse, error) {
	var (
		requestParams = gomerchant.Params{
			"payment_amount":   amount,
			"merchandise_type": params.MerchandiseType,
			"pc_mobile_type":   params.PCMobileType,
			"button_type":      params.ButtonType,
			"return_url":       params.ReturnUrl,
			"cancel_url":       params.CancelUrl,
		}
	)

	for i, g := range params.Goods {
		requestParams[fmt.Sprintf("goods[%d]", i)] = g.Name
		requestParams[fmt.Sprintf("goods_id[%d]", i)] = g.ID
		requestParams[fmt.Sprintf("goods_price[%d]", i)] = g.Price
		requestParams[fmt.Sprintf("goods_amount[%d]", i)] = g.Amount
	}
	var res gomerchant.ApplicationResponse
	results, err := paygent.Request("270", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			res.TransactionID = fmt.Sprint(paymentID)
		}

		if tradeGenerationDate, ok := results.Get("trade_generation_date"); ok {
			res.TradeGenerationDate = fmt.Sprint(tradeGenerationDate)
		}

		//Rakuten pay reponse redirect_html can not be find, so here need to do more logic
		redirectHTML := strings.Split(results.RawBody, "redirect_html=")
		if len(redirectHTML) == 2 {
			res.RedirectHTML = redirectHTML[1]

		}
		return res, nil
	}
	return res, err
}

// This is rakuten pay capture function
func (paygent *Paygent) RakutenPaySalesMessage(transactionID string) (gomerchant.CaptureResponse, error) {
	var (
		response      gomerchant.CaptureResponse
		requestParams = gomerchant.Params{
			"payment_id": transactionID,
		}
	)

	results, err := paygent.Request("271", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			response.TransactionID = fmt.Sprint(paymentID)
		}
	}
	response.Params = results.Params
	return response, err
}

// This is rakuten pay void function
func (paygent *Paygent) RakutenPayCancellationMessage(transactionID string) (gomerchant.VoidResponse, error) {
	var (
		response      gomerchant.VoidResponse
		requestParams = gomerchant.Params{
			"payment_id": transactionID,
		}
	)

	results, err := paygent.Request("272", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			response.TransactionID = fmt.Sprint(paymentID)
		}
	}
	response.Params = results.Params
	return response, err
}

func (paygent *Paygent) RakutenPayCorrectionMessage(transactionID string, amount uint) (gomerchant.RefundResponse, error) {
	var (
		response      gomerchant.RefundResponse
		requestParams = gomerchant.Params{
			"payment_id":     transactionID,
			"payment_amount": amount,
			//Because it's hard to specify every item price in our system
			//Like order with discount, it's so hard to calculate every item price and need equals total amounts.
			//So we set whole order as a goods to rakuten pay
			//If we could fix this problem later. Should be care with `del_flg`. Please read the documentation carefully [https://theplanttokyo.atlassian.net/browse/LAX-3319]
			"goods_id[0]":     gomerchant.RAKUTEN_PAY_PRODUCT_ID,
			"goods_price[0]":  amount,
			"goods_amount[0]": 1,
		}
	)

	results, err := paygent.Request("273", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			response.TransactionID = fmt.Sprint(paymentID)
		}
	}
	response.Params = results.Params
	return response, err
}

// Paypay authrioze function
func (paygent *Paygent) PayPayApplicationMessage(amount uint64, params gomerchant.ApplicationParams) (gomerchant.ApplicationResponse, error) {
	var (
		requestParams = gomerchant.Params{
			"payment_amount": amount,
			"return_url":     params.ReturnUrl,
			"cancel_url":     params.CancelUrl,
		}
	)
	var res gomerchant.ApplicationResponse
	results, err := paygent.Request("420", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			res.TransactionID = fmt.Sprint(paymentID)
		}

		if tradeGenerationDate, ok := results.Get("trade_generation_date"); ok {
			res.TradeGenerationDate = fmt.Sprint(tradeGenerationDate)
		}

		//Rakuten pay reponse redirect_html can not be find, so here need to do more logic
		redirectHTML := strings.Split(results.RawBody, "redirect_html=")
		if len(redirectHTML) == 2 {
			res.RedirectHTML = redirectHTML[1]

		}
		return res, nil
	}
	return res, err
}

func (paygent *Paygent) PayPaySalesMessage(transactionID string) (gomerchant.CaptureResponse, error) {
	var (
		response      gomerchant.CaptureResponse
		requestParams = gomerchant.Params{
			"payment_id": transactionID,
		}
	)

	results, err := paygent.Request("422", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			response.TransactionID = fmt.Sprint(paymentID)
		}
	}
	response.Params = results.Params
	return response, err
}

func (paygent *Paygent) PayPayCancelAndRefundMessage(transactionID string, amount uint) (gomerchant.RefundResponse, error) {
	var (
		response      gomerchant.RefundResponse
		requestParams = gomerchant.Params{
			"payment_id": transactionID,
		}
	)

	if amount > 0 {
		requestParams["repayment_amount"] = amount
	}

	results, err := paygent.Request("421", requestParams)
	if err == nil {
		if paymentID, ok := results.Get("payment_id"); ok {
			response.TransactionID = fmt.Sprint(paymentID)
		}
	}
	response.Params = results.Params
	return response, err
}

func (paygent *Paygent) Start3DS2Authentication(ctx context.Context, params gomerchant.Start3DS2AuthenticationParams) (response gomerchant.Start3DS2AuthenticationResponse, err error) {
	var (
		requestParams = gomerchant.Params{
			"trading_id":          params.OrderID,
			"payment_amount":      params.Amount,
			"term_url":            params.TermURL,
			"authentication_type": "01",
			"merchant_name":       paygent.Config.MerchantName,
		}
	)
	if params.PaymentMethod == nil {
		return response, gomerchant.ErrNotSupportedPaymentMethod
	}
	if savedCreditCard := params.PaymentMethod.SavedCreditCard; savedCreditCard != nil {
		requestParams["card_set_method"] = "customer"
		requestParams["customer_id"] = savedCreditCard.CustomerID
		requestParams["customer_card_id"] = savedCreditCard.CreditCardID
		// requestParams["card_conf_number"] = savedCreditCard.CVC

	} else if creditCard := params.PaymentMethod.CreditCard; creditCard != nil {
		requestParams["card_set_method"] = "direct"
		requestParams["card_number"] = creditCard.Number
		requestParams["card_valid_term"] = getValidTerm(creditCard)
		requestParams["card_conf_number"] = creditCard.CVC
	} else {
		return response, gomerchant.ErrNotSupportedPaymentMethod
	}
	for k, v := range params.Params {
		requestParams[k] = v
	}
	results, err := paygent.Request("450", requestParams)
	if err == nil {
		splitHTML := strings.Split(results.RawBody, "out_acs_html=")
		if len(splitHTML) == 2 {
			response.OutAcsHTML = strings.TrimSpace(splitHTML[1])
		}
		if res, ok := results.Get("result"); ok {
			response.Result = fmt.Sprint(res)
		}
	}
	response.Params = results.Params
	return response, err
}
