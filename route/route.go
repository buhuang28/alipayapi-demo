package route

import (
	"AliPayService/controller"
	"github.com/gin-gonic/gin"
)

func GinRun() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	payController := controller.NewVMQOrderController()
	router.POST("/createOrder", payController.CreateOrder)
	//router.GET("/deleteOrder",payController.CallBack)
	//router.GET("/callback",payController.CallBack)
	router.POST("/callback", payController.CallBack)
	router.Run(":8888")
}
