package alipay

import (
	"errors"
	"time"
)

// Common alipay common params
type Common struct {
	AppID        string            `json:"app_id"`
	Method       string            `json:"method"`
	Format       string            `json:"format"`
	Charset      string            `json:"charset"`
	SignType     string            `json:"sign_type"`
	Sign         string            `json:"-"`
	Timestamp    string            `json:"timestamp"`
	Version      string            `json:"version"`
	ReturnURL    string            `json:"return_url"`
	NotifyURL    string            `json:"notify_url"`
	AppAuthToken string            `json:"app_auth_token"`
	BizContent   map[string]string `json:"-"`
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
		commont.SignType = "RSA2"
	}

	if common.Timestamp == "" {
		common.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	if common.Version == "" {
		common.Version = "1.0"
	}
}
