package paygent

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/now"
	"github.com/qor/gomerchant"
)

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

func extractTransactionFromPaygentResponse(params Response) (transaction gomerchant.Transaction) {
	transaction.ID, _ = getPaymentID(params)

	if v, ok := params.Get("currency_code"); ok {
		transaction.Currency = fmt.Sprint(v)
	} else {
		transaction.Currency = "JCB"
	}

	if v, ok := params.Get("payment_init_date"); ok {
		if t, err := now.Parse(fmt.Sprint(v)); err == nil {
			transaction.CreatedAt = &t
		}
	}

	if v, ok := params.Get("payment_amount"); ok {
		if i, err := strconv.Atoi(fmt.Sprint(v)); err == nil {
			transaction.Amount = uint(i)
		}
	}

	if v, ok := params.Get("payment_status"); ok {
		transaction.Status = fmt.Sprint(v)
		switch transaction.Status {
		case "20", "30", "35":
			transaction.Paid = true
		case "40", "41":
			transaction.Paid = true
			transaction.Captured = true
		case "32", "33", "42", "55", "60":
			transaction.Cancelled = true
		}
	}

	return
}

// Paygent Status Code meaning
// 10 Applied
// 11 Authorization failed
// 13 3D Secure suspended
// 14 3D Secure authentication

// 20 Authorization succeeded
// 30 Sales being requested
// 35 Sales requested

// 40 Cleared
// 41 Cleared (sales cancellation overdue)

// 32 Authorization cancelled
// 33 Authorization expired
// 42 Sales being cancelled
// 55 Sales cancellation requested
// 60 Sales cancelled
