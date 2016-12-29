package paygent

import "fmt"

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

func getPaymentID(params paramsInterface) (string, bool) {
	paymentID, ok := params.Get("payment_id")
	return fmt.Sprint(paymentID), ok
}
