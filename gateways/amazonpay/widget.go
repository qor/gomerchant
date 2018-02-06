package amazonpay

// WidgetJSURL return widget js url
func (amazonPay AmazonPay) WidgetJSURL() string {
	if amazonPay.Config.ProductionMode {
		return "https://static-na.payments-amazon.com/OffAmazonPayments/us/js/Widgets.js"
	}
	return "https://static-na.payments-amazon.com/OffAmazonPayments/us/sandbox/js/Widgets.js"
}
