package router

import (
	"github.com/gin-gonic/gin"
	"serverList/controller"
)

//这边用来抛测试用
func initRmqRouter(e *gin.Engine) {
	route := controller.RmqReturnData{}
	e.GET("/rmq/publishMessage", route.RmqPublishMessage)
}
