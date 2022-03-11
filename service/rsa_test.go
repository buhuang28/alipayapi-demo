package service

import (
	"AliPayService/pay_data"
	"fmt"
	"testing"
	"time"
)

//创建订单，返回支付婊收款二维码链接
func TestCreateOrder(t *testing.T) {
	create := AliPayTradePreCreate("20000001", "0.01", "mytest", time.Now().Format(pay_data.ALIPAY_TIME_STAMP_FORMAT))
	fmt.Println(create)
}

//关闭超时订单
func TestCloseOrder(t *testing.T) {
	AliPayTradeClose("20000001")
}

//查询订单状态
func TestQueryOrder(t *testing.T) {
	orderId := "20000001"
	query := AliPayTradeQuery(orderId)
	fmt.Println(query)
}

//校验支付婊收款回调字符串是否正确，记得再按照支付婊官方的要求校验金额和订单号之类的信息。
func TestVerifySign(t *testing.T) {
	AliVerifySign(`支付宝回调的原生字符串数据`)
}
