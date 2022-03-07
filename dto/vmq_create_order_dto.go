package dto

//创建订单的数据
type VMQCreateOrderDto struct {
	PayId   string  `json:"payId"` //自定义订单id 对应OutTradeNo
	PayType int64   `json:"type"` //1支付婊，2微信 （VMQ里这两个是相反的）
	Price   float64 `json:"price"`//订单价格
	Sign    string  `json:"sign"`//
	Param   string  `json:"param"`
}

//创建订单返回的数据
type VMQCreateOrderResult struct {
	Code int64     `json:"code"`
	Msg  string    `json:"msg"`
	Data VMQOrderData `json:"data"`
}

type VMQOrderData struct {
	PayID       string  `json:"payId"`
	OrderID     string  `json:"orderId"`
	PayType     int64   `json:"payType"`
	Price       float64 `json:"price"`
	ReallyPrice float64 `json:"reallyPrice"`
	PayURL      string  `json:"payUrl"`
	IsAuto      int64   `json:"isAuto"`
	State       int64   `json:"state"`
	TimeOut     int64   `json:"timeOut"`
	Date        int64   `json:"date"`
}

type VMQCallBackData struct {
	PayId string `form:"payId"`
	Param string `form:"param"`
	Type int64 `form:"type"`
	Price float64 `form:"price"`
	ReallyPrice float64 `form:"reallyPrice"`
	Sign string `form:"sign"`
}
