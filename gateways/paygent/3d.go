package paygent

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/qor/gomerchant"
)

type SecureCodeParams struct {
	UserAgent  string
	TermURL    string
	HttpAccept string
}

type paramsInterface interface {
	Get(string) (interface{}, bool)
}

func get3DModeParams(params paramsInterface) (bool, *SecureCodeParams) {
	if value, ok := params.Get("3DMode"); ok {
		if fmt.Sprint(value) == "true" {
			if value, ok := params.Get("3DParams"); ok {
				if v, ok := value.(SecureCodeParams); ok {
					return true, &v
				}
				if v, ok := value.(*SecureCodeParams); ok {
					return true, v
				}
			}
		}
	}
	return false, nil
}

func (paygent *Paygent) SecureCodeAuthorize(amount uint64, secureCodeParams SecureCodeParams, params gomerchant.AuthorizeParams) (gomerchant.AuthorizeResponse, error) {
	if params.Params == nil {
		params.Params = gomerchant.Params{}
	}
	params.Set("3DMode", true)
	params.Set("3DParams", secureCodeParams)

	return paygent.Authorize(amount, params)
}

func (paygent *Paygent) CompleteAuthorize(paymentID string, params gomerchant.CompleteAuthorizeParams) (gomerchant.CompleteAuthorizeResponse, error) {
	if req, ok := params.Get("request"); ok {
		if request, ok := req.(*http.Request); ok {
			request.ParseForm()
			response, err := paygent.Request("024", gomerchant.Params{"MD": request.Form.Get("MD"), "PaRes": request.Form.Get("PaRes")})
			return gomerchant.CompleteAuthorizeResponse{Params: response.Params}, err
		}
	}
	return gomerchant.CompleteAuthorizeResponse{}, errors.New("no valid request params found")
}
