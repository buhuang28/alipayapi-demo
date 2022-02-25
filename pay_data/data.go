package pay_data

const (
	ALIPAY_GATEWAY_DEV            = "https://openapi.alipaydev.com/gateway.do"
	ALIPAY_GATEWAY                = "https://openapi.alipay.com/gateway.do"
	ALIPAY_TRADE_PRECREATE_METHOD = "alipay.trade.precreate"
	ALIPAY_TRADE_CLOSE_METHOD     = "alipay.trade.close"
	ALIPAY_TRADE_QUERY_METHOD     = "alipay.trade.query"

	ALIPAY_CHARSET           = "UTF-8"
	ALIPAY_SIGN_TYPE         = "RSA2"
	ALIPAY_VERSION           = "1.0"
	ALIPAY_SUBJECT           = "buhuang"
	ALIPAY_FORMAT            = "JSON"
	ALIPAY_SUCCESS_CODE      = "10000"
	ALIPAY_TIME_STAMP_FORMAT = "2006-01-02 15:04:05"
)

const (
	ALIPAY_TYPE int32 = 1
	WXPAY_TYPE  int32 = 2

	WAIT_PAY    int32 = -1
	PAY_OUTTIME int32 = -2 //也就是Close
	PAY_SUCCESS int32 = 1

	WAIT_TIME int64 = 300

	ALIPAY_ORDER_CLOSE   = "TRADE_CLOSED"
	ALIPAY_ORDER_SUCCESS = "TRADE_SUCCESS"
	ALIPAY_ORDER_WAIT    = "WAIT_BUYER_PAY"
)
