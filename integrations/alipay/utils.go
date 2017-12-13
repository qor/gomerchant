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

// Common alipay common params
type Common struct {
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
func (alipay *Alipay) Sign(common *Common, availableAttrs ...string) error {
	if common.Method == "" {
		return errors.New("method is invalid")
	}

	if common.AppID == "" {
		common.AppID = alipay.Config.AppID
	}

	if common.Format == "" {
		common.Format = "JSON"
	}

	if common.Charset == "" {
		common.Version = "utf-8"
	}

	if common.SignType == "" {
		common.SignType = "RSA2"
	}

	if common.Timestamp == "" {
		common.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	if common.Version == "" {
		common.Version = "1.0"
	}

	params := map[string]string{}
	result, err := json.Marshal(&params)
	if err == nil {
		err = json.Unmarshal(result, params)
		common.Sign, err = alipay.sign(params)
	}

	return err
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
		apiParams := []string{}

		for key, value := range params {
			if len(value) > 0 {
				apiParams = append(apiParams, key+"="+value)
			}
		}

		sort.Strings(apiParams)

		hash := crypto.SHA256
		if v, ok := params["sign_type"]; ok && v == "RSA" {
			hash = crypto.SHA1
		}
		h := hash.New()
		h.Write([]byte(strings.Join(apiParams, "&")))
		hashed := h.Sum(nil)

		s, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, hash, hashed)
		return base64.StdEncoding.EncodeToString(s), err
	}

	return "", err
}
