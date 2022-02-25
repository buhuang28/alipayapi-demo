package entity

type Order struct {
	ID uint `gorm:"primary_key"`
	//订单Id，用来创建支付婊订单（out_trade_no：商户订单号,保证唯一）
	OrderId string
	//支付宝订单号
	TradeNo string
	//购买人账号
	BuyerLogonId string
	//购买人Id
	BuyerUserId string
	//订单支付状态 ---和OrderSatus区别在于这个是查询接口查询到的
	TradeStatus string
	//订单价格
	Price float64
	//订单创建时间
	CreateTime int64
	//订单结束时间
	EndTime int64
	//订单平台 1:支付婊，2:微信
	Platform int32
	//订单状态  1付款,-1未付款,-2关闭
	OrderStatus int32
	//订单Body 会回传的数据
	Body string
	//备注，预留字段
	Note string
}
