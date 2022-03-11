package pay_data

import "fmt"

var (
	//是否使用沙箱数据
	DEV = false
	//应用私钥
	APP_PRIVATE_KEY = ""
	//应用公钥
	APP_PUBLIC_KEY = ""
	//支付宝公钥(用于回调鉴定)
	ALI_PUBLIC_KEY = ""
	//应用APPID
	ALIPAY_APP_ID = ""
	//网关
	ALIPAY_GATEWAY = ""
)

func init() {
	if DEV {
		APP_PRIVATE_KEY = DEV_ALIPAY_APP_PRIVATE_KEY
		ALI_PUBLIC_KEY = DEV_ALIPAY_ALI_PUBLIC_KEY
		ALIPAY_APP_ID = DEV_ALIPAY_APPID
		ALIPAY_GATEWAY = DEV_ALIPAY_GATEWAY
		fmt.Println("启用开发者设置")
	} else {
		//这里的替换成自己真实的支付婊应用数据
		APP_PRIVATE_KEY = REAL_APP_PRIVATE_KEY
		ALI_PUBLIC_KEY = REAL_ALI_PUBLIC_KEY
		APP_PUBLIC_KEY = REAL_APP_PUBLIC_KEY
		ALIPAY_APP_ID = REAL_ALIPAY_APPID
		ALIPAY_GATEWAY = REAL_ALIPAY_GATEWAY
		fmt.Println("启用真实环境设置")
	}
}

const (
	DEV_ALIPAY_GATEWAY            = "https://openapi.alipaydev.com/gateway.do"
	REAL_ALIPAY_GATEWAY           = "https://openapi.alipay.com/gateway.do"
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
