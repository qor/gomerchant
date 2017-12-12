package alipay

import (
	"errors"
	"net/url"
)

// Common alipay common params
type Common struct {
	AppID        string
	Method       string
	Format       string
	Charset      string
	SignType     string
	Sign         string
	Timestamp    string
	Version      string
	ReturnURL    string
	NotifyURL    string
	AppAuthToken string
	BizContent   map[string]string
}

// Sign common  params
func (alipay *Alipay) Sign(common *Common, availableAttrs ...string) error {
	values := url.Values{}

	if common.AppID == "" {
		common.AppID = alipay.Config.AppID
	}

	if common.Method == "" {
		return errors.New("method is invalid")
	}

	if common.Format == "" {
		common.Format = "json"
	}
}
