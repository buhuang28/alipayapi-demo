package main

import (
	"AliPayService/logs"
	"AliPayService/route"
	log "github.com/sirupsen/logrus"
)

func init() {
	logs.InitLog()
	log.Info("开始运行")
}

func main() {
	route.GinRun()
}
