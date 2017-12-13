package alipay

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/qor/gomerchant"
)

// Response aliyun response
type Response struct {
	Code    string `json:"code"`
	Mesg    string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
	Sign    string `json:"sign"`
	gomerchant.Params
}

var client = http.DefaultClient

// Request request alipay API
func (alipay *Alipay) Request(params *Params) (Response, error) {
	response := Response{}

	sign, err := alipay.Sign(params)
	if err != nil {
		return response, err
	}

	buf := strings.NewReader(sign)
	req, err := http.NewRequest(method, alipay.Config.APIDomain, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return response, err
	}
}
