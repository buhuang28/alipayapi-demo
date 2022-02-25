package service

import (
	"AliPayService/pay_data"
	"fmt"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	create := AliPayTradePreCreate("20000000", "0.01", "mytest", time.Now().Format(pay_data.ALIPAY_TIME_STAMP_FORMAT))
	fmt.Println(create)
}

func TestCloseOrder(t *testing.T) {
	AliPayTradeClose("20000000")
}

func TestQueryOrder(t *testing.T) {
	AliPayTradeQuery("20000000")
}

func TestVerifySign(t *testing.T) {
	AliVerifySign(`支付宝回调的原生字符串数据`)
}
