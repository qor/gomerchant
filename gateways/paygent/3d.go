package paygent

import "fmt"

type ThreeDomainSecureParams struct {
	UserAgent  string
	TermURL    string
	HttpAccept string
}

type paramsInterface interface {
	Get(string) (interface{}, bool)
}

func get3DModeParams(params paramsInterface) (bool, *ThreeDomainSecureParams) {
	if value, ok := params.Get("3DMode"); ok {
		if fmt.Sprint(value) == "true" {
			if value, ok := params.Get("3DParams"); ok {
				if v, ok := value.(ThreeDomainSecureParams); ok {
					return true, &v
				}
				if v, ok := value.(*ThreeDomainSecureParams); ok {
					return true, v
				}
			}
		}
	}
	return false, nil
}
