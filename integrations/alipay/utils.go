package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"sort"
	"strings"
	"time"
)

// Params alipay common params
type Params struct {
	AppID        string      `json:"app_id"`
	Method       string      `json:"method"`
	Format       string      `json:"format"`
	Charset      string      `json:"charset"`
	SignType     string      `json:"sign_type"`
	Sign         string      `json:"-"`
	Timestamp    string      `json:"timestamp"`
	Version      string      `json:"version"`
	ReturnURL    string      `json:"return_url"`
	NotifyURL    string      `json:"notify_url"`
	AppAuthToken string      `json:"app_auth_token"`
	BizContent   interface{} `json:"biz_content"`
}

// Sign common  params
func (alipay *Alipay) Sign(params *Params) (string, error) {
	if params.Method == "" {
		return "", errors.New("method is invalid")
	}

	if params.AppID == "" {
		params.AppID = alipay.Config.AppID
	}

	if params.Format == "" {
		params.Format = "JSON"
	}

	if params.Charset == "" {
		params.Version = "utf-8"
	}

	if params.SignType == "" {
		params.SignType = "RSA2"
	}

	if params.Timestamp == "" {
		params.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	if params.Version == "" {
		params.Version = "1.0"
	}

	reqParams := map[string]string{}
	result, err := json.Marshal(&params)
	if err == nil {
		err = json.Unmarshal(result, reqParams)
		params.Sign, err = alipay.sign(reqParams)
		reqParams["sign"] = params.Sign
	}

	return toSortedQuery(reqParams), err
}

func (alipay *Alipay) sign(params map[string]string) (s string, err error) {
	if alipay.Config.PrivateKey == "" {
		return "", errors.New("invalid private key")
	}

	block, _ := pem.Decode([]byte(alipay.Config.PrivateKey))
	if block == nil {
		return "", errors.New("invalid private key")
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		hash := crypto.SHA256
		if v, ok := params["sign_type"]; ok && v == "RSA" {
			hash = crypto.SHA1
		}
		h := hash.New()
		h.Write([]byte(toSortedQuery(params)))
		hashed := h.Sum(nil)

		s, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, hash, hashed)
		return base64.StdEncoding.EncodeToString(s), err
	}

	return "", err
}

func toSortedQuery(params map[string]string) string {
	apiParams := []string{}

	for key, value := range params {
		if len(value) > 0 {
			apiParams = append(apiParams, key+"="+value)
		}
	}

	sort.Strings(apiParams)
	return strings.Join(apiParams, "&")
}
