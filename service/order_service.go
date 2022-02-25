package service

import (
	"AliPayService/dao"
	"AliPayService/db"
	"AliPayService/dto"
	"AliPayService/pay_data"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	OrderChanMap     = make(map[string]chan struct{})
	OrderChanMapLock sync.Mutex
)

func AddOrderChan(orderId string, c chan struct{}) {
	OrderChanMapLock.Lock()
	defer OrderChanMapLock.Unlock()
	OrderChanMap[orderId] = c
}

func DeleteOrderChan(orderId string) {
	OrderChanMapLock.Lock()
	defer OrderChanMapLock.Unlock()
	delete(OrderChanMap, orderId)
}

//先创建订单，然后再请求支付婊接口  以后可能对接上微信接口
func CreateOrder(orderId, body string, price float64, platform int32) dto.Result {
	var order dao.Order
	order.OrderId = orderId
	order.Platform = platform
	order.Body = body //授权码
	order.Price = price
	order.OrderStatus = pay_data.WAIT_PAY

	now := time.Now()
	ch := make(chan string)
	go func() {
		//创建订单
		qrCode := AliPayTradePreCreate(orderId, fmt.Sprintf("%.2f", price), body, now.Format(pay_data.ALIPAY_TIME_STAMP_FORMAT))
		ch <- qrCode
		close(ch)
		//订单创建失败就返回
		if qrCode == "" {
			return
		}
		a := make(chan struct{})
		AddOrderChan(orderId, a)
		select {
		case <-a:
			//回调
			break
		case <-time.After(time.Minute * 5):
			//查询支付宝订单支付状态
			query := AliPayTradeQuery(orderId)
			var o dao.Order
			o.BuyerUserId = query.BuyerUserId
			o.TradeNo = query.TradeNo
			o.TradeStatus = query.TradeStatus
			o.BuyerLogonId = query.BuyerLogonId
			switch query.TradeStatus {
			case pay_data.ALIPAY_ORDER_SUCCESS:
				//修改数据库
				o.OrderStatus = pay_data.PAY_SUCCESS
			case pay_data.ALIPAY_ORDER_WAIT, pay_data.ALIPAY_ORDER_CLOSE:
				//修改数据库
				//其实应该都是返回的待支付的
				o.OrderStatus = pay_data.PAY_OUTTIME
			default:
				//压根就没扫过码的
				o.OrderStatus = pay_data.PAY_OUTTIME
			}
			updateError := o.UpdateOrderStatus(db.DbLink, orderId)
			if updateError != nil {
				log.Errorf("更新订单%v错误:%v", orderId, updateError)
			}
			close(a)
			break
		}
		DeleteOrderChan(orderId)
	}()

	var result dto.Result

	qrCode := <-ch

	order.CreateTime = now.Unix()
	order.EndTime = now.Unix() + pay_data.WAIT_TIME
	err := order.Create(db.DbLink)

	if err != nil {
		if qrCode != "" {
			//关闭订单
			AliPayTradeClose(orderId)
			close(OrderChanMap[orderId])
		}
		return result.Error(0, "创建订单失败")
	}
	if qrCode == "" {
		return result.Error(0, "创建订单失败")
	}
	return result.Success("", qrCode)
}

//回调
func AliPayCallBack() {

}
