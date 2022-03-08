package service

import (
	"AliPayService/dto"
	"AliPayService/pay_data"
	"AliPayService/util"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
	"net/url"
	"sort"
	"strings"
	"time"
)

//初始化一个公共参数
func NewAliClient(aliMethod, createTime string, biz dto.BizContent) map[string]string {
	data := make(map[string]string)
	data["app_id"] = pay_data.ALIPAY_APP_ID
	data["method"] = aliMethod
	data["charset"] = pay_data.ALIPAY_CHARSET
	data["sign_type"] = pay_data.ALIPAY_SIGN_TYPE
	data["timestamp"] = createTime
	data["version"] = pay_data.ALIPAY_VERSION
	//回调通知地址，用来确认用户是否付款
	if pay_data.ALIPAY_NOTIFY_URL != "" {
		data["notify_url"] = pay_data.ALIPAY_NOTIFY_URL
	}
	if pay_data.ALIPAY_FORMAT != "" {
		data["format"] = pay_data.ALIPAY_FORMAT
	}
	marshal, _ := json.Marshal(biz)
	data["biz_content"] = string(marshal)
	var signParamList []string
	for k, v := range data {
		if k != "" && v != "" {
			signParamList = append(signParamList, k+"="+v)
		}
	}
	sort.Strings(signParamList)
	var src = strings.Join(signParamList, "&")
	rsa, err := util.Sha256WithRsa(pay_data.APP_PRIVATE_KEY, src)
	if err != nil {
		return nil
	}
	data["sign"] = rsa
	return data
}

//预创建订单
func AliPayTradePreCreate(orderId, price, body, createTime string) string {
	var biz dto.BizContent
	biz.OutTradeNo = orderId
	biz.TotalAmount = price
	biz.Subject = pay_data.ALIPAY_SUBJECT
	biz.Body = body
	data := NewAliClient(pay_data.ALIPAY_TRADE_PRECREATE_METHOD, createTime, biz)
	postJson, resp := util.PostForm(pay_data.ALIPAY_GATEWAY, nil, nil, data)
	if !postJson {
		return ""
	}
	bytes, err := fastjson.ParseBytes(resp.Data)
	if err != nil {
		return ""
	}
	aliResponse := bytes.Get("alipay_trade_precreate_response")
	code := string(aliResponse.GetStringBytes("code"))
	if code != pay_data.ALIPAY_SUCCESS_CODE {
		return ""
	}
	return string(aliResponse.GetStringBytes("qr_code"))
}

//关闭订单
func AliPayTradeClose(orderId string) {
	var biz dto.BizContent
	biz.OutTradeNo = orderId
	data := NewAliClient(pay_data.ALIPAY_TRADE_CLOSE_METHOD, time.Now().Format(pay_data.ALIPAY_TIME_STAMP_FORMAT), biz)
	success, resp := util.GetRequest(pay_data.ALIPAY_GATEWAY, nil, nil, data)
	if !success {
		return
	}
	bytes, err := fastjson.ParseBytes(resp.Data)
	if err != nil {
		log.Infof("关闭订单 %v 失败:%v,返回数据:%v", orderId, err, string(resp.Data))
		return
	}
	closeResponse := bytes.Get("alipay_trade_close_response")
	code := string(closeResponse.GetStringBytes("code"))
	if code == pay_data.ALIPAY_SUCCESS_CODE {
		log.Infof("关闭订单 %v 成功", orderId)
	} else {
		log.Infof("关闭订单 %v 失败", orderId)
	}
}

//查询订单
func AliPayTradeQuery(orderId string) dto.QuertOrderResult {
	var queryResult dto.QuertOrderResult
	var biz dto.BizContent
	biz.OutTradeNo = orderId
	data := NewAliClient(pay_data.ALIPAY_TRADE_QUERY_METHOD, time.Now().Format(pay_data.ALIPAY_TIME_STAMP_FORMAT), biz)
	success, resp := util.GetRequest(pay_data.ALIPAY_GATEWAY, nil, nil, data)
	if !success {
		success, resp = util.GetRequest(pay_data.ALIPAY_GATEWAY, nil, nil, data)
		if !success {
			marshal, _ := json.Marshal(data)
			log.Infof("查询订单 %v 失败,查询参数:%v", orderId, string(marshal))
			return queryResult
		}
	}
	bytes, err := fastjson.ParseBytes(resp.Data)
	if err != nil {
		log.Infof("查询订单 %v 失败:%v,返回数据:%v", orderId, err, string(resp.Data))
		return queryResult
	}
	closeResponse := bytes.Get("alipay_trade_query_response")
	code := string(closeResponse.GetStringBytes("code"))
	if code != pay_data.ALIPAY_SUCCESS_CODE {
		log.Infof("查询不到该订单%v", orderId)
		return queryResult
	}
	queryResult.BuyerLogonId = string(closeResponse.GetStringBytes("buyer_logon_id"))
	queryResult.BuyerUserId = string(closeResponse.GetStringBytes("buyer_user_id"))
	queryResult.TradeNo = string(closeResponse.GetStringBytes("trade_no"))
	queryResult.TradeStatus = string(closeResponse.GetStringBytes("trade_status"))
	return queryResult
}

//验签
/*
3. 需要严格按照如下描述校验通知数据的正确性：
• 商户需要验证该通知数据中的 out_trade_no 是否为商户系统中创建的订单号。
• 判断 total_amount 是否确实为该订单的实际金额（即商户订单创建时的金额）。
• 校验通知中的 seller_id（或者 seller_email) 是否为 out_trade_no 这笔单据的对应的操作方（有的时候，一个商户可能有多个 seller_id/seller_email）。
上述有任何一个验证不通过，则表明本次通知是异常通知，务必忽略。在上述验证通过后商户必须根据支付宝不同类型的业务通知，
正确的进行不同的业务处理，并且过滤重复的通知结果数据。在支付宝的业务通知中，只有交易通知状态为 TRADE_SUCCESS 或 TRADE_FINISHED 时，
支付宝才会认定为买家付款成功。
*/
func AliVerifySign(callbackData string) bool {
	callbackData, _ = url.QueryUnescape(callbackData)
	split := strings.Split(callbackData, "&")
	sign := ""
	var signContent []string
	for _, v := range split {
		if strings.Contains(v, "sign") {
			if len(v) > 20 {
				sign = v[5:]
			}
			continue
		}
		signContent = append(signContent, v)
	}
	sort.Strings(signContent)
	src := strings.Join(signContent, "&")
	return util.Rsa2PubSign(pay_data.ALI_PUBLIC_KEY, src, sign)
}
