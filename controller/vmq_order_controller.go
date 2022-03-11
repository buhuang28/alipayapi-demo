package controller

import (
	"AliPayService/dto"
	"AliPayService/pay_data"
	"AliPayService/service"
	"AliPayService/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"time"
)

type VMQOrderController struct{}

func NewVMQOrderController() VMQOrderController {
	return VMQOrderController{}
}

//创建订单
func (a *VMQOrderController) CreateOrder(c *gin.Context) {
	var orderDto dto.VMQCreateOrderDto
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &orderDto)
	if err != nil {
		log.Info(err)
		log.Info(orderDto)
		log.Info("json反序列化失败1")
		c.JSON(200, nil)
		return
	}
	sign := util.GetSign(orderDto.PayId, orderDto.Param,
		strconv.FormatInt(orderDto.PayType, 10),
		fmt.Sprintf("%.2f", orderDto.Price),
		pay_data.SESSIONKEY) //这里的 SESSIONKEY 只是一个任意字符串，用于计算MD5
	if sign != orderDto.Sign {
		c.JSON(200, nil)
		log.Info("签名对不上")
		return
	}
	order := service.CreateOrder(orderDto.PayId, orderDto.Param, orderDto.Price, pay_data.ALIPAY_TYPE)
	var result dto.VMQCreateOrderResult
	result.Code = order.Status
	result.Data.PayURL = order.Data.(string)
	result.Data.OrderID = orderDto.PayId
	result.Data.PayID = orderDto.PayId
	c.JSON(200, result)
}

//支付回调
func (a *VMQOrderController) CallBack(c *gin.Context) {
	log.Info("收到支付回调")
	raw, err := c.GetRawData()
	if err != nil {
		log.Error(err)
	}
	var callBackData dto.AliCallBackData
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
	c.Bind(&callBackData)
	if err != nil {
		fmt.Println(err)
		c.String(200, "error")
		return
	}
	success := service.AliVerifySign(string(raw))
	log.Info("回调校验:", success)
	if success {
		orderChan := service.GetOrderChan(callBackData.OutTradeNo)
		close(orderChan)
		c.String(200, "success")
		go func() {
			for i := 0; i < 10; i++ {
				ok := Notice(callBackData)
				if !ok {
					time.Sleep(time.Minute * 20)
				} else {
					return
				}
			}
		}()
	} else {
		c.String(200, "error")
	}
	return
}

func Notice(callBackData dto.AliCallBackData) bool {
	noticeData := make(map[string]string)
	float, _ := strconv.ParseFloat(callBackData.ReceiptAmount, 10)
	noticeData["price"] = fmt.Sprintf("%.2f", float)
	noticeData["payId"] = callBackData.OutTradeNo
	noticeData["type"] = "1"
	noticeData["reallyPrice"] = fmt.Sprintf("%.2f", float)
	noticeData["param"] = callBackData.Body
	noticeData["sign"] = util.GetSign(callBackData.OutTradeNo, callBackData.Body,
		"1", fmt.Sprintf("%.2f", float),
		fmt.Sprintf("%.2f", float), pay_data.SESSIONKEY) //这里的 SESSIONKEY 只是一个任意字符串，用于计算MD5
	ok, resp := util.GetRequest("http://127.0.0.1:8882/PayCallBackOrder", nil, nil, noticeData)
	log.Info("回调结果:", ok, string(resp.Data))
	if string(resp.Data) != "success" {
		return false
	}
	return true
}

//关闭订单
